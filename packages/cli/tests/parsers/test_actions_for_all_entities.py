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
Test in this module must be applied for all entity commands
"""
import http
import json
import pathlib
import typing

import pytest
from click.testing import CliRunner
from pytest_mock import MockFixture

from odahuflow.cli.utils.error_handler import IGNORE_NOT_FOUND_ERROR_MESSAGE, \
    ID_AND_FILE_GIVEN_ERROR_MESSAGE, ID_AND_FILE_MISSED_ERROR_MESSAGE
from odahuflow.cli.utils.output import JSON_OUTPUT_FORMAT, \
    JSONPATH_OUTPUT_FORMAT
from odahuflow.sdk.clients.api import WrongHttpStatusCode
from .data import ENTITY_ID, EntityTestData, ROUTER, DEPLOYMENT, PACKAGING_INTEGRATION, PACKAGING, TOOLCHAIN, TRAINING, \
    CONNECTION, generate_entities_for_test

ENTITY_TEST_DATA: typing.Dict[str, EntityTestData] = {
    "connection": CONNECTION,
    "training": TRAINING,
    "toolchain": TOOLCHAIN,
    "packaging": PACKAGING,
    "packaging_integration": PACKAGING_INTEGRATION,
    "deployments": DEPLOYMENT,
    "routes": ROUTER,
}

WRONG_OUTPUT_FORMAT = 'wrong-format'


def pytest_generate_tests(metafunc):
    generate_entities_for_test(metafunc, list(ENTITY_TEST_DATA.keys()))


@pytest.fixture
def entity_test_data(request) -> EntityTestData:
    return ENTITY_TEST_DATA[request.param]


def test_get(mocker: MockFixture, cli_runner: CliRunner, entity_test_data: EntityTestData):
    client_mock = mocker.patch.object(
        entity_test_data.entity_client.__class__, 'get',
        return_value=entity_test_data.entity,
    )

    result = cli_runner.invoke(entity_test_data.click_group, [
        'get', '--id', ENTITY_ID, '-o', JSON_OUTPUT_FORMAT,
    ], obj=entity_test_data.entity_client)

    client_mock.assert_called_once_with(ENTITY_ID)
    assert result.exit_code == 0
    assert json.loads(result.output) == [entity_test_data.entity.to_dict()]


def test_get_all(mocker: MockFixture, cli_runner: CliRunner,
                 entity_test_data: EntityTestData):
    client_mock = mocker.patch.object(
        entity_test_data.entity_client.__class__,
        'get_all',
        return_value=[entity_test_data.entity],
    )

    result = cli_runner.invoke(entity_test_data.click_group,
                               ['get', '-o', JSON_OUTPUT_FORMAT],
                               obj=entity_test_data.entity_client)

    client_mock.assert_called_once_with()
    assert result.exit_code == 0
    assert json.loads(result.output) == [entity_test_data.entity.to_dict()]


def test_get_jsonpath(mocker: MockFixture, cli_runner: CliRunner,
                      entity_test_data: EntityTestData):
    client_mock = mocker.patch.object(
        entity_test_data.entity_client.__class__,
        'get',
        return_value=entity_test_data.entity,
    )

    result = cli_runner.invoke(entity_test_data.click_group,
                               ['get', '--id', ENTITY_ID, '-o',
                                f'{JSONPATH_OUTPUT_FORMAT}=[*].id'],
                               obj=entity_test_data.entity_client)
    client_mock.assert_called_once_with(ENTITY_ID)
    assert result.exit_code == 0
    assert result.output.strip() == ENTITY_ID


def test_get_default_output_format(mocker: MockFixture, cli_runner: CliRunner,
                                   entity_test_data: EntityTestData):
    client_mock = mocker.patch.object(entity_test_data.entity_client.__class__,
                                      'get',
                                      return_value=entity_test_data.entity)

    result = cli_runner.invoke(entity_test_data.click_group,
                               ['get', '--id', ENTITY_ID],
                               obj=entity_test_data.entity_client)

    client_mock.assert_called_once_with(ENTITY_ID)
    assert result.exit_code == 0
    assert ENTITY_ID in result.stdout


def test_get_wrong_output_format(cli_runner: CliRunner,
                                 entity_test_data: EntityTestData):
    result = cli_runner.invoke(entity_test_data.click_group,
                               ['get', '--id', ENTITY_ID, '-o', WRONG_OUTPUT_FORMAT],
                               obj=entity_test_data.entity_client)

    assert result.exit_code != 0
    assert f'invalid choice: {WRONG_OUTPUT_FORMAT}' in result.output


def test_edit_wrong_kind(tmp_path: pathlib.Path, cli_runner: CliRunner,
                         entity_test_data: EntityTestData):
    entity_file = tmp_path / "entity.yaml"
    entity_file.write_text(
        json.dumps({**entity_test_data.entity.to_dict(), **{'kind': 'Wrong'}}))

    result = cli_runner.invoke(entity_test_data.click_group,
                               ['edit', '-f', entity_file],
                               obj=entity_test_data.entity_client)

    assert result.exit_code != 0
    assert "Unknown kind of object: 'Wrong'" in str(result.exception)


def test_create_wrong_kind(tmp_path: pathlib.Path, cli_runner: CliRunner,
                           entity_test_data: EntityTestData):
    entity_file = tmp_path / "entity.yaml"
    entity_file.write_text(
        json.dumps({**entity_test_data.entity.to_dict(), **{'kind': 'Wrong'}}))

    result = cli_runner.invoke(entity_test_data.click_group,
                               ['create', '-f', entity_file],
                               obj=entity_test_data.entity_client)

    assert result.exit_code != 0
    assert "Unknown kind of object: 'Wrong'" in str(result.exception)


def test_delete_id_and_file_present(tmp_path, cli_runner: CliRunner,
                                    entity_test_data: EntityTestData):
    entity_file = tmp_path / "entity.yaml"
    entity_file.write_text(
        json.dumps({**entity_test_data.entity.to_dict(), **{'kind': entity_test_data.kind}}))
    result = cli_runner.invoke(entity_test_data.click_group,
                               ['delete', '--id', 'some-id', '-f', entity_file],
                               obj=entity_test_data.entity_client)

    assert result.exit_code != 0
    assert ID_AND_FILE_GIVEN_ERROR_MESSAGE in str(result.exception)


def test_delete_no_id_or_file_present(tmp_path, cli_runner: CliRunner,
                                    entity_test_data: EntityTestData):
    entity_file = tmp_path / "entity.yaml"
    entity_file.write_text(
        json.dumps({**entity_test_data.entity.to_dict(), **{'kind': entity_test_data.kind}}))
    result = cli_runner.invoke(entity_test_data.click_group,
                               ['delete'],
                               obj=entity_test_data.entity_client)

    assert result.exit_code != 0
    assert ID_AND_FILE_MISSED_ERROR_MESSAGE in str(result.exception)


def test_delete_wrong_kind(tmp_path: pathlib.Path, cli_runner: CliRunner,
                           entity_test_data: EntityTestData):
    entity_file = tmp_path / "entity.yaml"
    entity_file.write_text(json.dumps({**entity_test_data.entity.to_dict(), **{'kind': 'Wrong'}}))

    result = cli_runner.invoke(entity_test_data.click_group,
                               ['delete', '-f', entity_file],
                               obj=entity_test_data.entity_client)

    assert result.exit_code != 0
    assert "Unknown kind of object: 'Wrong'" in str(result.exception)


def test_delete_ignore_not_found_enabled(mocker: MockFixture, cli_runner: CliRunner,
                                         entity_test_data: EntityTestData):
    client_mock = mocker.patch.object(entity_test_data.entity_client.__class__,
                                      'delete',
                                      side_effect=WrongHttpStatusCode(
                                          http.HTTPStatus.NOT_FOUND))

    result = cli_runner.invoke(entity_test_data.click_group,
                               ['delete', '--id', ENTITY_ID,
                                '--ignore-not-found'],
                               obj=entity_test_data.entity_client)

    client_mock.assert_called_once_with(ENTITY_ID)
    assert result.exit_code == 0
    assert IGNORE_NOT_FOUND_ERROR_MESSAGE.format(kind=entity_test_data.kind, id=ENTITY_ID) in result.stdout


def test_delete_ignore_not_found_disabled(mocker: MockFixture, cli_runner: CliRunner,
                                          entity_test_data: EntityTestData):
    client_mock = mocker.patch.object(entity_test_data.entity_client.__class__,
                                      'delete',
                                      side_effect=WrongHttpStatusCode(
                                          http.HTTPStatus.NOT_FOUND))

    result = cli_runner.invoke(entity_test_data.click_group,
                               ['delete', '--id', ENTITY_ID],
                               obj=entity_test_data.entity_client)

    client_mock.assert_called_once_with(ENTITY_ID)
    assert result.exit_code != 0
    assert "Got error from server" in str(result.exception)


def test_delete_ignore_not_found_enabled_http_code(
        mocker: MockFixture, cli_runner: CliRunner, entity_test_data: EntityTestData,
):
    client_mock = mocker.patch.object(entity_test_data.entity_client.__class__,
                                      'delete',
                                      side_effect=WrongHttpStatusCode(
                                          http.HTTPStatus.BAD_REQUEST))

    result = cli_runner.invoke(entity_test_data.click_group,
                               ['delete', '--id', ENTITY_ID,
                                '--ignore-not-found'],
                               obj=entity_test_data.entity_client)

    client_mock.assert_called_once_with(ENTITY_ID)
    assert result.exit_code != 0
    assert "Got error from server" in str(result.exception)
