#
#    Copyright 2017 EPAM Systems
#
#    Licensed under the Apache License, Version 2.0 (the "License");
#    you may not use this file except in compliance with the License.
#    You may obtain a copy of the License at
#
#        http://www.apache.org/licenses/LICENSE-2.0
#
#    Unless required by applicable law or agreed to in writing, software
#    distributed under the License is distributed on an "AS IS" BASIS,
#    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
#    See the License for the specific language governing permissions and
#    limitations under the License.
#
"""
API client
"""
import json
import logging
import random
import string
import sys
import threading
from collections.abc import AsyncIterable
from typing import Any, Callable, Dict, Iterator, Mapping, Optional, Tuple, Union
from urllib.parse import urlencode, urlparse
from http.client import responses

from requests.adapters import HTTPAdapter
from urllib3 import Retry


import aiohttp
import requests
import requests.exceptions

import odahuflow.sdk.config
from odahuflow.sdk.clients.oauth_handler import OAuthLoginResult, do_client_cred_authentication, do_refresh_token, \
    start_oauth2_callback_handler
from odahuflow.sdk.clients.oidc import fetch_openid_configuration
from odahuflow.sdk.config import update_config_file
from odahuflow.sdk.definitions import API_VERSION

LOGGER = logging.getLogger(__name__)


class WrongHttpStatusCode(Exception):
    """
    Exception for wrong HTTP status code
    """

    def __init__(self, status_code: int, http_result: Dict[str, str] = None):
        """
        Initialize Wrong Http Status Code exception

        :param status_code: HTTP status code
        :param http_result: HTTP data
        """
        if http_result is None:
            http_result = {}

        default_message = responses[status_code]

        message = http_result.get("message", default_message)

        super().__init__(f'Got error from server: {message} (status: {status_code})')

        self.status_code = status_code


class EntityAlreadyExists(WrongHttpStatusCode):
    """
    Exception for 409 conflict: entity already exists
    """
    pass


class APIConnectionException(Exception):
    """
    Exception that says that client can not reach API server
    """

    pass


class IncorrectAuthorizationToken(APIConnectionException):
    """
    Exception that says that provided API authorization token is incorrect
    """

    pass


class IncorrectClientCredentials(APIConnectionException):
    """
    Exception that says that provided API authorization token is incorrect
    """

    pass


class LoginRequired(APIConnectionException):
    """
    Exception that says that login is required to do calls
    """

    pass


def get_authorization_redirect(web_redirect: str, after_login: Callable) -> str:
    """
    Try to detect, parse and build OAuth2 redirect

    :param web_redirect: returned redirect
    :param after_login: function that have to be called after successful login
    :return: str -- new redirect
    """
    loc = urlparse(web_redirect)

    state = ''.join(random.choice(string.ascii_letters) for _ in range(10))
    local_check_address = start_oauth2_callback_handler(after_login, state, web_redirect)

    get_parameters = {
        'client_id': odahuflow.sdk.config.ODAHUFLOWCTL_OAUTH_CLIENT_ID,
        'response_type': 'code',
        'state': state,
        'redirect_uri': local_check_address,
        'scope': odahuflow.sdk.config.ODAHUFLOWCTL_OAUTH_SCOPE
    }
    web_redirect = f'{loc.scheme}://{loc.netloc}{loc.path}?{urlencode(get_parameters)}'
    return web_redirect


class URLBuilder:

    def __init__(self, base_url: str = odahuflow.sdk.config.API_URL):
        self._base_url = base_url
        self._version = API_VERSION

    @property
    def base_url(self):
        return self._base_url

    def build_url(self, url_template):
        sub_url = url_template.format(version=self._version)
        target_url = self._base_url.strip('/') + sub_url
        return target_url

    def build_request_kwargs(self,
                             url_template: str,
                             payload: Mapping[Any, Any] = None,
                             action: str = 'GET',
                             stream: bool = False,
                             token: Optional[str] = None
                             ) -> Dict[str, Any]:
        target_url = self.build_url(url_template)
        headers = {}
        if token:
            headers['Authorization'] = f'Bearer {token}'
        if stream:
            headers['Content-type'] = 'text/event-stream'

        request_kwargs = {
            'method': action,
            'url': target_url,
            'params' if action.lower() == 'get' else 'json': payload,
            'headers': headers
        }
        return request_kwargs


class Authenticator:

    def __init__(self,
                 client_id: str, client_secret: str, non_interactive: bool, base_url: str,
                 token: Optional[str] = odahuflow.sdk.config.API_TOKEN,
                 issuer_url: Optional[str] = odahuflow.sdk.config.ISSUER_URL):

        self._client_id = client_id
        self._client_secret = client_secret
        self._non_interactive = non_interactive
        self._base_url = base_url
        self._interactive_login_finished = threading.Event()
        self._token = token
        self._issuer_url = issuer_url

        # Force if set
        if odahuflow.sdk.config.ODAHUFLOWCTL_NONINTERACTIVE:
            self._non_interactive = True

    @staticmethod
    def login_required(response):
        """
        Check whether the login is required or a client is already authorised
        :param response:
        :return:
        """

        # We assume if there were redirects then credentials are out of date and we can refresh or build auth url

        def not_authorized_resp_code():
            if hasattr(response, 'status_code'):
                return response.status_code == 401
            elif hasattr(response, 'status'):
                return response.status == 401
            return False

        return bool(response.history) or not_authorized_resp_code()

    def login(self, url: str, limit_stack=False):
        """
        Authorise client. Next methods are used by priority:
        1. Refreshing token (if token exists)
        2. Interactive mode (if enabled)
        :return: None if success, IncorrectAuthorizationToken exception otherwise
        """

        # If it is an error after refreshed token - fail
        if limit_stack:
            raise IncorrectAuthorizationToken(
                f'{self._credentials_error_status} even after refreshing. \n'
                'Please try to log in again'
            )

        # use default value if self._issuer_url is empty
        self._issuer_url = self._issuer_url or odahuflow.sdk.config.API_ISSUING_URL

        LOGGER.debug('Redirect has been detected. Trying to refresh a token')
        if self._refresh_token_exists:
            LOGGER.debug('Refresh token for %s has been found, trying to use it', odahuflow.sdk.config.API_ISSUING_URL)
            self._login_with_refresh_token()
        elif self._client_id and self._client_secret and fetch_openid_configuration(self._issuer_url):
            self._login_with_client_credentials()
        elif self._interactive_mode_enabled:
            # Start interactive flow
            self._login_interactive_mode(url)
        else:
            raise IncorrectAuthorizationToken(
                f'{self._credentials_error_status}.\n'
                'Please provide correct temporary token or disable non interactive mode'
            )

    @property
    def token(self):
        return self._token

    def _login_with_refresh_token(self):
        login_result = do_refresh_token(odahuflow.sdk.config.API_REFRESH_TOKEN, odahuflow.sdk.config.API_ISSUING_URL)
        if not login_result:
            raise IncorrectAuthorizationToken(
                'Refresh token is not correct.\n'
                'Please login again'
            )
        else:
            self._update_config_with_new_oauth_config(login_result)

    def _login_with_client_credentials(self):

        # use default value if self._issuer_url is empty
        self._issuer_url = self._issuer_url or odahuflow.sdk.config.API_ISSUING_URL

        login_result = do_client_cred_authentication(
            issue_token_url=fetch_openid_configuration(self._issuer_url), client_id=self._client_id,
            client_secret=self._client_secret
        )
        if not login_result:
            raise IncorrectClientCredentials(
                'Client credentials are not correct.\n'
                'Please login again'
            )
        else:
            self._update_config_with_new_oauth_config(login_result)

    def _login_interactive_mode(self, url):
        self._interactive_login_finished.clear()
        target_url = get_authorization_redirect(url, self._after_login)
        print('%s. \nPlease open %s' % (self._credentials_error_status, target_url))
        self._interactive_login_finished.wait()

    @property
    def _refresh_token_exists(self):
        return odahuflow.sdk.config.API_REFRESH_TOKEN and odahuflow.sdk.config.API_ISSUING_URL

    @property
    def _interactive_mode_enabled(self):
        return not self._non_interactive

    @property
    def _credentials_error_status(self):
        if self._token:
            credentials_error_status = 'Credentials are not correct'
        else:
            credentials_error_status = 'Credentials are missed'
        return credentials_error_status

    def _update_config_with_new_oauth_config(self, login_result: OAuthLoginResult) -> None:
        """
        Update config with new oauth credentials

        :param login_result: result of login
        :return: None
        """
        self._token = login_result.id_token
        update_config_file(API_URL=self._base_url,
                           API_TOKEN=login_result.id_token,
                           API_REFRESH_TOKEN=login_result.refresh_token,
                           API_ACCESS_TOKEN=login_result.access_token,
                           ISSUER_URL=self._issuer_url,
                           API_ISSUING_URL=login_result.issuing_url)

    def _after_login(self, login_result: OAuthLoginResult) -> None:
        """
        Handle action after login

        :param login_result: result of login
        :return: None
        """
        self._interactive_login_finished.set()
        self._update_config_with_new_oauth_config(login_result)
        print('You have been authorized on endpoint %s as %s / %s' %
              (self._base_url, login_result.user_name, login_result.user_email))
        sys.exit(0)


def _handle_query_response(text: str, payload: Mapping[Any, Any], status_code: int) -> Dict:
    try:
        answer = json.loads(text)
        LOGGER.debug('Got answer: {!r} with code {} for URL {!r}'
                     .format(answer, status_code, payload))
    except ValueError:
        answer = {}

    if status_code == 409:
        raise EntityAlreadyExists(status_code, answer)
    if 400 <= status_code < 600:
        raise WrongHttpStatusCode(status_code, answer)

    LOGGER.debug('Query has been completed, parsed and validated')
    return answer


class RemoteAPIClient:
    """
    Base API client
    """

    def __init__(self,
                 base_url: str = odahuflow.sdk.config.API_URL,
                 token: Optional[str] = odahuflow.sdk.config.API_TOKEN,
                 client_id: Optional[str] = '',
                 client_secret: Optional[str] = '',
                 retries: Optional[int] = odahuflow.sdk.config.RETRY_ATTEMPTS,
                 timeout: Optional[Union[int, Tuple[int, int]]] = 10,
                 non_interactive: Optional[bool] = True,
                 issuer_url: Optional[str] = odahuflow.sdk.config.ISSUER_URL):
        """
        Build client

        :param base_url: base url, for example: http://api.example.com
        :param token: token for token based auth
        :param client_id: client_id for Client Credentials OAuth2 flow
        :param client_secret: client_secret for Client Credentials OAuth2 flow
        :param retries: command retries or less then 2 if disabled
        :param timeout: timeout for connection in seconds. 0 for disabling
        :param non_interactive: disable any interaction
        """
        self._base_url = base_url
        self.url_builder = URLBuilder(base_url)
        self.authenticator = Authenticator(client_id, client_secret, non_interactive, base_url, token, issuer_url)
        self.timeout = timeout
        self.retries = retries

        retry_strategy = Retry(
            total=retries,
            status_forcelist=[429, 500, 502, 503, 504],
            backoff_factor=odahuflow.sdk.config.BACKOFF_FACTOR
        )
        adapter = HTTPAdapter(max_retries=retry_strategy)
        self.default_client = requests.Session()
        self.default_client.mount("http://", adapter)
        self.default_client.mount("https://", adapter)

    @property
    def timeout(self):
        return self._timeout

    @timeout.setter
    def timeout(self, value):
        self._timeout = value

    @classmethod
    def construct_from_other(cls, other):
        """
        Construct API-based client from another API-based client

        :param other: API-based client to get connection options from
        :return: self -- new client
        """
        return cls(
            other.url_builder.base_url, other.authenticator.token,
            retries=other.retries, timeout=other.timeout
        )

    def _request(self,
                 url_template: str,
                 payload: Mapping[Any, Any] = None,
                 action: str = 'GET',
                 stream: bool = False,
                 timeout: Optional[int] = None,
                 limit_stack: bool = False,
                 client: requests.Session = None):
        """
        Make HTTP request

        :param url_template: target URL
        :param payload: data to be placed in body of request
        :param action: request action, e.g. get / post / delete
        :param stream: use stream mode or not
        :param timeout: custom timeout in seconds (overrides default). 0 for disabling
        :param limit_stack: do not start refreshing token if it is possible
        :return: :py:class:`requests.Response` -- response
        """

        request_kwargs = self.url_builder.build_request_kwargs(url_template, payload, action, stream,
                                                               token=self.authenticator.token)
        connection_timeout = self.timeout

        try:
            if client:
                response = client.request(timeout=connection_timeout, stream=stream, **request_kwargs)
            else:
                response = self.default_client.request(timeout=connection_timeout, stream=stream, **request_kwargs)
        except Exception as raised_exception:
            raise APIConnectionException('Can not reach {}'.format(self._base_url)) from raised_exception

        if self.authenticator.login_required(response):
            try:
                self.authenticator.login(str(response.url), limit_stack=limit_stack)
            except IncorrectAuthorizationToken as login_exc:
                raise login_exc
            return self._request(
                url_template,
                payload=payload,
                action=action,
                stream=stream,
                timeout=timeout,
                limit_stack=True
            )
        else:
            return response

    def query(self, url_template: str, payload: Mapping[Any, Any] = None, action: str = 'GET'):
        """
        Perform query to API server

        :param url_template: url template from odahuflow.const.api
        :param payload: payload (will be converted to JSON) or None
        :param action: HTTP method (GET, POST, PUT, DELETE)
        :return: dict[str, any] -- response content
        """
        response = self._request(url_template, payload, action)
        return _handle_query_response(response.text, payload, response.status_code)

    def stream(self, url_template: str, action: str = 'GET', params: Mapping[str, Any] = None) -> Iterator[str]:
        """
        Perform query to API server

        :param url_template: url template from odahuflow.const.api
        :param params: payload (will be converted to JSON) or None
        :param action: HTTP method (GET, POST, PUT, DELETE)
        :return: response content
        """
        response = self._request(url_template, payload=params, action=action, stream=True)

        with response:
            if not response.ok:
                raise WrongHttpStatusCode(response.status_code)

            for line in response.iter_lines():
                yield line.decode("utf-8")

    def info(self):
        """
        Perform info query on API server

        :return:
        """
        try:
            return self.query("/health")
        except ValueError as e:
            raise APIConnectionException(*e.args) from e


class AsyncRemoteAPIClient:

    def __init__(self,
                 base_url: str = odahuflow.sdk.config.API_URL,
                 token: Optional[str] = odahuflow.sdk.config.API_TOKEN,
                 client_id: Optional[str] = '',
                 client_secret: Optional[str] = '',
                 retries: Optional[int] = 3,
                 timeout: Optional[Union[int, Tuple[int, int]]] = 10,
                 non_interactive: Optional[bool] = True,
                 issuer_url: Optional[str] = odahuflow.sdk.config.ISSUER_URL):
        """
        Build client

        :param base_url: base url, for example: http://api.example.com
        :param token: token for token based auth
        :param client_id: client_id for Client Credentials OAuth2 flow
        :param client_secret: client_secret for Client Credentials OAuth2 flow
        :param retries: command retries or less then 2 if disabled
        :param timeout: timeout for connection in seconds. 0 for disabling
        :param non_interactive: disable any interaction
        """
        self._base_url = base_url
        self.url_builder = URLBuilder(base_url)
        self.authenticator = Authenticator(client_id, client_secret, non_interactive, base_url, token, issuer_url)
        self.timeout = timeout
        self.retries = retries

    async def _request(
            self, url_template: str,
            payload: Mapping[Any, Any] = None,
            action: str = 'GET',
            stream=False,
            session: aiohttp.ClientSession = None
    ) -> AsyncIterable:
        """
        Perform async request to API server

        :param url_template: url template from odahuflow.const.api
        :param payload: payload (will be converted to JSON) or None
        :param action: HTTP method (GET, POST, PUT, DELETE)
        :param stream: use stream mode or not
        :param session: aiohttp.Session for making response
        :return: AsyncIterable of data. If stream = False, only one line will be returned
        """
        assert session is not None

        request_kwargs = self.url_builder.build_request_kwargs(url_template, payload, action, stream=False,
                                                               token=self.authenticator.token)
        left_retries = self.retries
        raised_exception = None
        while left_retries > 0:
            try:
                async with session.request(**request_kwargs) as resp:
                    LOGGER.debug(resp)
                    if self.authenticator.login_required(resp):
                        raise LoginRequired()

                    if stream:
                        LOGGER.debug('Status code: "{}", Response: "<stream>"'.format(resp.status))
                        async for line in self._handle_stream_response(resp):
                            yield line
                    else:
                        resp_text = await resp.text()
                        LOGGER.debug('Status code: "{}", Response: "{}"'.format(resp.status, resp_text))
                        data = _handle_query_response(resp_text, payload, resp.status)
                        yield data
                    break
            except aiohttp.ClientConnectionError as exception:
                LOGGER.error('Failed to connect to {}: {}. Retrying'.format(self._base_url, exception))
                raised_exception = exception
                left_retries -= 1
        else:
            raise APIConnectionException('Can not reach {}'.format(self._base_url)) from raised_exception

    @staticmethod
    async def _handle_stream_response(response: aiohttp.ClientResponse, chunk_size=500) -> AsyncIterable:
        """
        Helper method that do next things:
        1) Async iterating over response content by chunks
        2) Decode bytes to text
        3) Build string lines from text stream
        :param response:
        :param chunk_size:
        :return:
        """

        if 400 <= response.status < 600:
            raise WrongHttpStatusCode(response.status)

        encoding = response.get_encoding()
        pending = None

        async for bytes_chunk in response.content.iter_chunked(chunk_size):
            chunk = bytes_chunk.decode(encoding)

            if pending is not None:
                chunk = pending + chunk

            lines = chunk.splitlines()

            if lines and lines[-1] and chunk and lines[-1][-1] == chunk[-1]:
                pending = lines.pop()
            else:
                pending = None

            for line in lines:
                yield line

        if pending is not None:
            yield pending

    async def _request_with_login(
            self, url_template: str, payload: Mapping[Any, Any] = None, action: str = 'GET', stream=False
    ):
        """
        Perform async request to API server.
        Create session under the hood.
        Make two attempts to get data. Second attempt is used after login if required

        :param url_template: url template from odahuflow.const.api
        :param payload: payload (will be converted to JSON) or None
        :param action: HTTP method (GET, POST, PUT, DELETE)
        :param stream: use stream mode or not
        :return: AsyncIterable of data. If stream = False, only one line will be returned
        """
        async with aiohttp.ClientSession(conn_timeout=self.timeout, trust_env=True) as session:
            try:
                async for data in self._request(url_template, payload, action, stream, session):
                    yield data
            except LoginRequired:
                self.authenticator.login(
                    self.url_builder.build_url(url_template)
                )
                async for data in self._request(url_template, payload, action, stream, session):
                    yield data

    async def query(self, url_template: str, payload: Mapping[Any, Any] = None, action: str = 'GET') -> Any:
        """
        Perform query to API server.

        :param url_template: url template from odahuflow.const.api
        :param payload: payload (will be converted to JSON) or None
        :param action: HTTP method (GET, POST, PUT, DELETE)
        :return: dict[str, any] -- response content
        """
        resp = None
        async for res in self._request_with_login(url_template, payload, action, stream=False):
            resp = res
        return resp

    async def stream(self, url_template: str, action: str = 'GET', params: Mapping[str, Any] = None) -> AsyncIterable:
        """
        Perform query to API server in stream mode

        :param url_template: url template from odahuflow.const.api
        :param action: HTTP method (GET, POST, PUT, DELETE)
        :param params: payload (will be converted to JSON) or None
        :return: async iterable that return response content line by line
        """
        async for line in self._request_with_login(url_template, action=action, stream=True, payload=params):
            yield line

    async def info(self):
        """
        Perform info query on API server
        :return:
        """
        try:
            return await self.query("/health")
        except ValueError:
            pass

    @property
    def timeout(self):
        return self._timeout

    @timeout.setter
    def timeout(self, value):
        self._timeout = value

    @classmethod
    def construct_from_other(cls, other):
        """
        Construct API-based client from another API-based client

        :param other: API-based client to get connection options from
        :return: self -- new client
        """
        return cls(
            other.url_builder.base_url, other.authenticator.token,
            retries=other.retries, timeout=other.timeout
        )
