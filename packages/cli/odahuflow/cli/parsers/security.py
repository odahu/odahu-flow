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
import sys

import click
from odahuflow.sdk.clients import api
from odahuflow.sdk.config import update_config_file

LOG = logging.getLogger(__name__)


@click.command()
@click.option('--url', 'api_host', help='API server host')
@click.option('--token', help='API server jwt token')
def login(api_host: str, token: str):
    """
    Authorize on API endpoint.
    Check that credentials is correct and save to the config
    """
    try:
        api_clint = api.RemoteAPIClient(api_host, token, non_interactive=bool(token))

        api_clint.info()
        update_config_file(API_URL=api_host, API_TOKEN=token)

        print('Success! Credentials have been saved.')
    except api.IncorrectAuthorizationToken as wrong_token:
        LOG.error('Wrong authorization token\n%s', wrong_token)
        sys.exit(1)


@click.command()
def logout():
    """
    Remove all authorization data from the configuration file
    """
    update_config_file(API_URL=None,
                       API_TOKEN=None,
                       API_REFRESH_TOKEN=None,
                       API_ACCESS_TOKEN=None,
                       API_ISSUING_URL=None)

    print('All authorization credentials have been removed')
