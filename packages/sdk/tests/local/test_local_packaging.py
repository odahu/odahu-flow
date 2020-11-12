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
import os
import json
import docker

import pytest
from pytest_mock import MockFixture

from odahuflow.sdk.local import packaging
from odahuflow.sdk.local.packaging import start_package
from odahuflow.sdk.models import K8sPackager, ModelPackaging, ModelPackagingSpec, PackagingIntegration, \
    PackagingIntegrationSpec

# Format: ['artifact_name', 'artifact_path',
#          'expected_artifact_name', expected_artifact_path]
test_data = [
    (
        'wine-1.0', '/odahu/training',
        'wine-1.0', '/odahu/training'
    ),
    (
        'wine-1.0.zip', '/odahu/training',
        'wine-1.0', '/odahu/training'
    ),
    (
        'wine-1.0.zip.zip', None,
        'wine-1.0.zip', '/odahu/default_output'
    )
]

DEFAULT_OUTPUT_DIR = '/odahu/default_output'


@pytest.mark.parametrize(['artifact_name', 'artifact_path',
                          'expected_artifact_name', 'expected_artifact_path'],
                         test_data)
def test_start_package__artifact_name_artifact_path(artifact_name, artifact_path,
                                                    expected_artifact_name, expected_artifact_path,
                                                    mocker: MockFixture):
    packager = K8sPackager(
        model_packaging=ModelPackaging(spec=ModelPackagingSpec(artifact_name=artifact_name)),
        # mocking packaging_integration default_image
        packaging_integration=PackagingIntegration(spec=PackagingIntegrationSpec(default_image='default_image')))

    create_mp_config_file_mock = mocker.patch.object(packaging, 'create_mp_config_file')
    config_mock = mocker.patch.object(packaging, 'config')

    mocker.patch.object(docker, 'from_env')
    mocker.patch.object(json, 'dumps')

    mocker.patch.object(packaging, 'stream_container_logs')
    mocker.patch.object(packaging, 'raise_error_if_container_failed')
    read_mp_result_file_mock = mocker.patch.object(packaging, 'read_mp_result_file')

    config_mock.LOCAL_MODEL_OUTPUT_DIR = DEFAULT_OUTPUT_DIR

    start_package(packager, artifact_path)

    expected_packager = K8sPackager(
        model_packaging=ModelPackaging(spec=ModelPackagingSpec(artifact_name=expected_artifact_name)),
        # mocking packaging_integration default_image
        packaging_integration=PackagingIntegration(spec=PackagingIntegrationSpec(default_image='default_image')))

    expected_full_artifact_path = os.path.join(expected_artifact_path, expected_artifact_name)

    create_mp_config_file_mock.assert_called_with(expected_full_artifact_path, expected_packager)
    read_mp_result_file_mock.assert_called_with(expected_full_artifact_path)
