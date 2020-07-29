#  Copyright 2020 EPAM Systems
#
#  Licensed under the Apache License, Version 2.0 (the "License");
#  you may not use this file except in compliance with the License.
#  You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
#  Unless required by applicable law or agreed to in writing, software
#  distributed under the License is distributed on an "AS IS" BASIS,
#  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
#  See the License for the specific language governing permissions and
#  limitations under the License.

from unittest.mock import Mock
from click.testing import CliRunner
from pytest_mock import MockFixture
from odahuflow.cli.parsers import security
from odahuflow.sdk.clients.api import APIConnectionException


# TODO: Write a complete set of unit tests for security module

def test_login_invalid_url(cli_runner: CliRunner, mocker: MockFixture):
    """
    Tests fix for issue #257 (https://github.com/odahu/odahu-flow/issues/257)
    If
    :param cli_runner: Click CLI runner fixture
    :param mocker: pytest mocker fixture
    """

    mocker.patch.object(security, 'update_config_file')

    api_client_class_mock: Mock = mocker.patch.object(security.api, 'RemoteAPIClient')
    api_client: Mock = api_client_class_mock.return_value
    api_client.info.side_effect = APIConnectionException

    result = cli_runner.invoke(security.login, ['--url', '--help'])
    assert result.exit_code != 0
    assert isinstance(result.exception, APIConnectionException)
