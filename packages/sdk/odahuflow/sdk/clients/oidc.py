#
#    Copyright 2020 EPAM Systems
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
OIDC configuration tool
"""

import requests

WELL_KNOWN_CONFIGURATION_URL = '.well-known/openid-configuration'

# Some OpenID Provider configuration keys
# (defined at https://openid.net/specs/openid-connect-discovery-1_0.html#ProviderMetadata)
TOKEN_ENDPOINT = 'token_endpoint'


def fetch_openid_configuration(issuer: str) -> str:
    """
    returns token_endpoint (candidate for API_ISSUING_URL variable)
    from ISSUER_URL variable

    :param issuer: ISSUER_URL config variable
    :return: token_endpoint
    """
    conf = OpenIdProviderConfiguration(issuer)
    conf.fetch_configuration()

    return conf.token_endpoint


class OpenIdProviderConfiguration:
    """
    OpenID Provider Configuration according to https://openid.net/specs/openid-connect-discovery-1_0.html
    """

    def __init__(self, issuer: str):
        """

        :param issuer: URL to OIDC issuer
        """
        self._issuer = issuer
        self._config_json = None

    def extract_value(self, key):
        if self._config_json is None:
            raise RuntimeError('You must fetch OpenID Provider configuration via .fetch_configuration() '
                               'or .async_fetch_configuration() method before refer some properties')

        return self._config_json.get(key)

    def fetch_configuration(self) -> None:
        """
        Fetch OpenID Provider configuration making HTTP request
        :return:
        """
        response = requests.get(self.configuration_url)
        if response.status_code == 200:
            self._config_json = response.json()
        elif 400 <= response.status_code < 600:
            raise RuntimeError(f'Some error during attempt to fetch OpenID Provider configuration. '
                               f'Reason: {response.reason}')
        else:
            raise RuntimeError(f'Not expected status code: {response.status_code}')

    @property
    def configuration_url(self) -> str:
        return f'{self._issuer}/{WELL_KNOWN_CONFIGURATION_URL}'

    @property
    def token_endpoint(self) -> str:
        return self.extract_value(TOKEN_ENDPOINT)
