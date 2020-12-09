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
Test user authentication
"""

import click
from odahuflow.robot.profiler_loader import (
    get_variables,
    AUTH_TOKEN_PARAM_NAME,
    API_URL_PARAM_NAME,
)
from odahuflow.sdk.config import update_config_file


@click.command()
@click.argument("profile", type=click.Path(exists=True, dir_okay=True, readable=True))
def main(profile: str) -> None:
    """
    Authenticate the test user from odahuflow profile file and save its credentials

    \f
    :param profile: file with odahuflow secrets
    """
    test_variables = get_variables(profile)
    api_token, api_url = (
        test_variables[AUTH_TOKEN_PARAM_NAME],
        test_variables[API_URL_PARAM_NAME],
    )

    update_config_file(API_URL=api_url, API_TOKEN=api_token)
    click.echo("Config was updated!")
