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

import json
import pathlib
import typing

import pytest
from click.testing import CliRunner
from pytest_mock import MockFixture

from .data import EntityTestData, generate_entities_for_test, TRAINING, PACKAGING, DEPLOYMENT, ROUTER

ENTITY_TEST_DATA: typing.Dict[str, EntityTestData] = {
    "training": TRAINING,
    "packaging": PACKAGING,
    "deployments": DEPLOYMENT,
    "routes": ROUTER,
}


def pytest_generate_tests(metafunc):
    generate_entities_for_test(metafunc, list(ENTITY_TEST_DATA.keys()))


@pytest.fixture
def entity_test_data(request) -> EntityTestData:
    return ENTITY_TEST_DATA[request.param]


def test_edit(tmp_path: pathlib.Path, mocker: MockFixture, cli_runner: CliRunner,
              entity_test_data: EntityTestData):
    entity_file = tmp_path / "entity.yaml"
    entity_file.write_text(
        json.dumps({**entity_test_data.entity.to_dict(), **{'kind': entity_test_data.kind}}))
    edit_client_mock = mocker.patch.object(entity_test_data.entity_client.__class__,
                                           'edit',
                                           return_value=entity_test_data.entity)
    get_client_mock = mocker.patch.object(entity_test_data.entity_client.__class__,
                                          'get',
                                          return_value=entity_test_data.entity)

    result = cli_runner.invoke(entity_test_data.click_group,
                               ['edit', '-f', entity_file],
                               obj=entity_test_data.entity_client)

    assert result.exit_code == 0
    edit_client_mock.assert_called_once_with(entity_test_data.entity)
    get_client_mock.assert_called_once()


def test_edit_no_wait(mocker: MockFixture, tmp_path, cli_runner: CliRunner,
                      entity_test_data: EntityTestData):
    entity_file = tmp_path / "entity.yaml"
    entity_file.write_text(
        json.dumps(
            {**entity_test_data.entity.to_dict(), **{'kind': entity_test_data.kind}}))
    edit_client_mock = mocker.patch.object(entity_test_data.entity_client.__class__,
                                           'edit',
                                           return_value=entity_test_data.entity)
    get_client_mock = mocker.patch.object(entity_test_data.entity_client.__class__,
                                          'get',
                                          return_value=entity_test_data.entity)

    result = cli_runner.invoke(entity_test_data.click_group,
                               ['edit', '-f', entity_file, '--no-wait'],
                               obj=entity_test_data.entity_client)

    assert result.exit_code == 0
    edit_client_mock.assert_called_once_with(entity_test_data.entity)
    get_client_mock.assert_not_called()


def test_create(tmp_path: pathlib.Path, mocker: MockFixture, cli_runner: CliRunner,
                entity_test_data: EntityTestData):
    entity_file = tmp_path / "entity.yaml"
    entity_file.write_text(
        json.dumps({**entity_test_data.entity.to_dict(), **{'kind': entity_test_data.kind}}))
    edit_client_mock = mocker.patch.object(entity_test_data.entity_client.__class__,
                                           'create',
                                           return_value=entity_test_data.entity)
    get_client_mock = mocker.patch.object(entity_test_data.entity_client.__class__,
                                          'get',
                                          return_value=entity_test_data.entity)

    result = cli_runner.invoke(entity_test_data.click_group,
                               ['create', '-f', entity_file],
                               obj=entity_test_data.entity_client)

    assert result.exit_code == 0
    edit_client_mock.assert_called_once_with(entity_test_data.entity)
    get_client_mock.assert_called_once()


def test_create_no_wait(mocker: MockFixture, tmp_path, cli_runner: CliRunner,
                        entity_test_data: EntityTestData):
    entity_file = tmp_path / "entity.yaml"
    entity_file.write_text(
        json.dumps(
            {**entity_test_data.entity.to_dict(), **{'kind': entity_test_data.kind}}))
    edit_client_mock = mocker.patch.object(entity_test_data.entity_client.__class__,
                                           'create',
                                           return_value=entity_test_data.entity)
    get_client_mock = mocker.patch.object(entity_test_data.entity_client.__class__,
                                          'get',
                                          return_value=entity_test_data.entity)

    result = cli_runner.invoke(entity_test_data.click_group,
                               ['create', '-f', entity_file, '--no-wait'],
                               obj=entity_test_data.entity_client)

    assert result.exit_code == 0
    edit_client_mock.assert_called_once_with(entity_test_data.entity)
    get_client_mock.assert_not_called()
