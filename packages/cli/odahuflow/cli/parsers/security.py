#
#    Copyright 2019 EPAM Systems
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
Security commands for odahuflow cli
"""
import logging

import click
from click import UsageError

from odahuflow.sdk.clients import api
from odahuflow.sdk.config import update_config_file

LOGGER = logging.getLogger(__name__)


def _reset_credentials():
    """
    Clean credentials from config file
    :return:
    """
    update_config_file(API_URL=None,
                       API_TOKEN=None,
                       API_REFRESH_TOKEN=None,
                       API_ACCESS_TOKEN=None,
                       API_ISSUING_URL=None,
                       ODAHUFLOWCTL_OAUTH_CLIENT_ID=None,
                       ODAHUFLOWCTL_OAUTH_CLIENT_SECRET=None,
                       ISSUER_URL=None,
                       MODEL_HOST=None)


@click.command()
@click.option('--url', 'api_host', help='API server host', required=True)
@click.option('--token', help='API server jwt token')
@click.option('--client_id', help='client_id for OAuth2 Client Credentials flow')
@click.option('--client_secret', help='client_secret for OAuth2 Client Credentials flow')
@click.option('--issuer', help='OIDC Issuer URL')
def login(api_host: str, token: str, client_id: str, client_secret: str, issuer: str):
    """
    Authorize on API endpoint.
    Check that credentials is correct and save to the config
    """

    # clean config from previous credentials
    _reset_credentials()

    # update config
    update_config_file(
        API_URL=api_host,
        API_TOKEN=token,
        ODAHUFLOWCTL_OAUTH_CLIENT_ID=client_id,
        ODAHUFLOWCTL_OAUTH_CLIENT_SECRET=client_secret,
        ISSUER_URL=issuer
    )

    # set predicates
    is_token_login = bool(token)
    is_client_cred_login = bool(client_id) or bool(client_secret)
    is_interactive_login = not (is_token_login or is_client_cred_login)

    # validate
    if is_token_login and is_client_cred_login:
        raise UsageError('You should use either --token or --client_id/--client_secret to login. '
                         'Otherwise skipp all options to launch interactive login mode')
    if is_client_cred_login and (not client_id or not client_secret):
        raise UsageError('You must pass both client_id and client_secret to client_credentials login')
    if is_client_cred_login and not issuer:
        raise UsageError('You must provide --issuer parameter to do client_credentials login. '
                         'Or set ISSUER_URL option in config file')

    # login
    api_client = api.RemoteAPIClient(api_host, token, client_id, client_secret,
                                     non_interactive=not is_interactive_login,
                                     issuer_url=issuer)
    try:
        api_client.info()
        print('Success! Credentials have been saved.')
    except api.IncorrectAuthorizationToken as wrong_token:
        LOGGER.error('Wrong authorization token\n%s', wrong_token)
        raise
    except api.APIConnectionException as connection_exc:
        LOGGER.error(f'Failed to connect to API host! '
                     f'Error: {connection_exc}')
        raise


@click.command()
def logout():
    """
    Remove all authorization data from the configuration file
    """
    _reset_credentials()
    click.echo('All authorization credentials have been removed')
