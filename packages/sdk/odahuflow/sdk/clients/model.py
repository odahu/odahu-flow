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
Model HTTP API client and utils
"""
import json
import logging

import requests
from urllib3.exceptions import HTTPError

import odahuflow.sdk.config
from odahuflow.sdk.clients.deployment import ModelDeploymentClient
from odahuflow.sdk.clients.route import ModelRouteClient
from odahuflow.sdk.clients.api import Authenticator, IncorrectAuthorizationToken
from odahuflow.sdk.utils import ensure_function_succeed


LOGGER = logging.getLogger(__name__)


def calculate_url(base_url: str, model_route: str = None, model_deployment: str = None,
                  url_prefix: str = None, mr_client: ModelRouteClient = None):
    """
    Calculate url for model

    :param base_url: base model server url
    :param model_route: model route name
    :param model_deployment: model deployment name. Default route URL will be returned
    :param url_prefix: model prefix
    :param mr_client: ModelRoute client to use
    :return: model url
    """
    if not base_url:
        raise ValueError("Base url is required")

    if url_prefix:
        return f'{base_url}{url_prefix}'

    if model_route:
        if mr_client is None:
            mr_client = ModelRouteClient(base_url=base_url)

        model_route = mr_client.get(model_route)

        LOGGER.debug('Found model route: %s', model_route)
        return model_route.status.edge_url

    if model_deployment:
        md_client = ModelDeploymentClient(base_url=base_url)
        model_route = md_client.get_default_route(model_deployment)
        LOGGER.debug('Found default model route: %s', model_route)
        return model_route.status.edge_url

    raise NotImplementedError("Cannot create a model url")


class ModelClient:
    """
    Model HTTP client
    """

    def __init__(self,
                 base_url,
                 model_route=None,
                 model_deployment=None,
                 url_prefix=None,
                 token=None,
                 http_client=requests,
                 http_exception=requests.exceptions.RequestException,
                 timeout=None,
                 client_id='',
                 client_secret='',
                 issuer_url=''):
        """
        Build client

        :param base_url: model base url
        :type base_url: str
        :param model_route: model route name
        :type model_route: str
        :param model_deployment: model deployment name
        :type model_deployment: str
        :param url_prefix: model url prefix
        :type url_prefix: str
        :param token: API token value to use (default: None)
        :type token: str
        :param http_client: HTTP client (default: requests)
        :type http_client: python class that implements requests-like post & get methods
        :param http_exception: http_client exception class, which can be thrown by http_client in case of some errors
        :type http_exception: python class that implements Exception class interface
        :param timeout: timeout for connections
        :type timeout: int
        :param client_id: client_id for Client Credentials OAuth2 flow
        :type client_id: str
        :param client_secret: client_secret for Client Credentials OAuth2 flow
        :type client_secret: str
        :param issuer_url: url for credential login
        :type issuer_url: str
        """
        self._url = calculate_url(base_url, model_route, model_deployment, url_prefix)
        self._base_url = base_url
        self._http_client = http_client
        self._http_exception = http_exception
        self._timeout = timeout

        token = token if token else odahuflow.sdk.config.API_TOKEN
        issuer_url = issuer_url if issuer_url else odahuflow.sdk.config.ISSUER_URL

        self._authenticator = Authenticator(
            client_id=client_id,
            client_secret=client_secret,
            non_interactive=True,
            base_url=self._base_url,
            token=token,
            issuer_url=issuer_url
        )

        LOGGER.debug('Model client params: %s, %s, %s, %s, %s', self._url, token, http_client, http_exception, timeout)

    @property
    def api_url(self):
        """
        Build API root URL

        :return: str -- api root url
        """
        return '{host}/api/model'.format(host=self._url)

    @property
    def info_url(self):
        """
        Build API info URL

        :return: str -- info url
        """
        return self.api_url + '/info'

    @staticmethod
    def _parse_response(response):
        """
        Parse model response (requests or FlaskClient)

        :param response: model HTTP response
        :type response: object with .text or .data and .status_code attributes
        :return: dict -- parsed response
        """
        data = response.text if hasattr(response, 'text') else response.data

        if isinstance(data, bytes):
            data = data.decode('utf-8')

        try:
            data = json.loads(data)
        except ValueError:
            pass

        if not 200 <= response.status_code < 400:
            url = response.url if hasattr(response, 'url') else None
            raise Exception('Wrong status code returned: {}. Data: "{}". URL: "{}"'
                            .format(response.status_code, data, url))

        return data

    @property
    def _additional_kwargs(self):
        """
        Get additional HTTP client key-value arguments like timestamp

        :return: dict -- additional kwargs
        """
        kwargs = {}
        if self._authenticator.token:
            kwargs['headers'] = {'Authorization': f'Bearer {self._authenticator.token}'}
        if self._timeout is not None:
            kwargs['timeout'] = self._timeout
        return kwargs

    def _request(self, http_method, url, data=None, files=None, retries=10, sleep=3, **kwargs):
        """
        Send request with provided method and other parameters
        :param http_method: HTTP method
        :type http_method: str
        :param url: url to send request to
        :type url: str
        :param data: request data
        :type data: any
        :param files: files to send with request
        :type files: dict
        :param retries: How many times to retry executing a request
        :type retries: int
        :param sleep: How much time to sleep between retries in case of errors
        :type sleep: int
        :return: dict -- parsed model response
        """
        http_method = http_method.lower()

        client_method = getattr(self._http_client, http_method)

        if data:
            kwargs['data'] = data
        if files:
            kwargs['files'] = files

        def check_function():
            try:
                return client_method(url, **kwargs)
            except (self._http_exception, HTTPError) as e:
                LOGGER.error('Failed to connect to {}: {}.'.format(url, e))

        response = ensure_function_succeed(check_function, retries, sleep)
        if response is None:
            raise self._http_exception('HTTP request failed')

        if self._authenticator.login_required(response):
            try:
                self._authenticator.login(str(response.url), limit_stack=False)
            except IncorrectAuthorizationToken as login_exc:
                raise login_exc
            return self._request(
                http_method,
                url,
                data=data,
                files=files,
                retries=retries,
                sleep=sleep,
                **self._additional_kwargs
            )

        return self._parse_response(response)

    def invoke(self, **parameters):
        """
        Invoke model with parameters

        :param parameters: parameters for model
        :type parameters: dict[str, object] -- dictionary with parameters
        :return: dict -- parsed model response
        """
        return self._request('post', f'{self.api_url}/invoke', **self._additional_kwargs, json=parameters)

    def info(self):
        """
        Get model info

        :return: dict -- parsed model info
        """
        return self._request('get', self.info_url, **self._additional_kwargs)
