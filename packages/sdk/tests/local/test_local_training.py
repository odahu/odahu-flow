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
from typing import Optional
from unittest.mock import ANY

import pytest
from pytest_mock import MockFixture

from odahuflow.sdk.local import training
from odahuflow.sdk.local.training import compile_artifact_name_template, TemplateContext, start_train, \
    DEFAULT_MODEL_DIR_TEMPLATE
from odahuflow.sdk.models import K8sTrainer, ModelTraining, ModelTrainingSpec, ModelIdentity


# Format: [(template, context, expected_result), ...]
test_data = [
    (
        '{{ .Name }}_abc_{{ .RandomUUID }}',
        TemplateContext(Name='wine', Version='2.3', RandomUUID='123-123'),
        'wine_abc_123-123'
    ),
    (
        '{{ .Name }}-{{ .Version }}-{{ .RandomUUID }}',
        TemplateContext(Name='model_name', Version='99', RandomUUID='uuid_here'),
        'model_name-99-uuid_here'
    ),
    (
        'Model {{   .Name      }}',
        TemplateContext(Name=None, Version=None, RandomUUID=None),
        'Model None'
    )
]


@pytest.mark.parametrize(['template', 'ctx', 'expected'], test_data)
def test_compile_artifact_name_template(template, ctx, expected):
    assert compile_artifact_name_template(template, ctx) == expected


DEFAULT_OUTPUT_DIR = '/odahu/default_output'


@pytest.mark.parametrize(['output_dir', 'expected_output_dir'],
                         [(None, DEFAULT_OUTPUT_DIR),
                          ('/custom/output', '/custom/output')]
                         )
def test_start_training__output_dir(output_dir: Optional[str], expected_output_dir: str, mocker: MockFixture):
    """
    Tests output_dir parameter. If it is provided, the result model directory is created under it.
    Otherwise, it is created in a default output path from config
    """
    trainer = K8sTrainer(model_training=ModelTraining(spec=ModelTrainingSpec(model=ModelIdentity())))

    mocker.patch.object(training, 'create_mt_config_file')
    config_mock = mocker.patch.object(training, 'config')
    os_makedirs_mock = mocker.patch.object(training.os, 'makedirs')
    compile_artifact_name_template_mock = mocker.patch.object(training, 'compile_artifact_name_template')
    launch_training_container_mock = mocker.patch.object(training, 'launch_training_container')
    launch_gppi_validation_container_mock = mocker.patch.object(training, 'launch_gppi_validation_container')

    config_mock.LOCAL_MODEL_OUTPUT_DIR = DEFAULT_OUTPUT_DIR
    compile_artifact_name_template_mock.return_value = 'model_dir_name'

    start_train(trainer, output_dir)

    expected_model_dir_path = os.path.join(expected_output_dir, 'model_dir_name')

    os_makedirs_mock.assert_called_with(expected_model_dir_path, exist_ok=ANY)
    launch_training_container_mock.assert_called_with(trainer, expected_model_dir_path)
    launch_gppi_validation_container_mock.assert_called_with(trainer, expected_model_dir_path)


@pytest.mark.parametrize(['template', 'expected_template'],
                         [(None, DEFAULT_MODEL_DIR_TEMPLATE),
                          ('{{ .Prop1 }}__{{ .prop2 }}', '{{ .Prop1 }}__{{ .prop2 }}'),
                          ('{{ .Prop1 }}_abc.zip', '{{ .Prop1 }}_abc'),
                          (' {{ .Prop1   }}_fizz.zip.zip  ', '{{ .Prop1   }}_fizz.zip')],
                         )
def test_start_training__artifact_name_template(template: Optional[str],
                                                expected_template: str, mocker: MockFixture):
    """
    Tests artifactNameTemplate property from training configuration. If it is provided,
    the result model directory is named according to it. Otherwise, the default template is used.
    If artifact name template ends with .zip it is trimmed.
    """
    model = ModelIdentity(artifact_name_template=template)
    trainer = K8sTrainer(model_training=ModelTraining(spec=ModelTrainingSpec(model=model)))

    mocker.patch.object(training, 'create_mt_config_file')
    os_makedirs_mock = mocker.patch.object(training.os, 'makedirs')
    compile_artifact_name_template_mock = mocker.patch.object(training, 'compile_artifact_name_template')
    launch_training_container_mock = mocker.patch.object(training, 'launch_training_container')
    launch_gppi_validation_container_mock = mocker.patch.object(training, 'launch_gppi_validation_container')

    compile_artifact_name_template_mock.return_value = 'compiled_model_dir_name'

    output_dir = '/test/output'
    start_train(trainer, output_dir)

    compile_artifact_name_template_mock.assert_called_with(go_template=expected_template, context=ANY)

    expected_model_dir_path = os.path.join(output_dir, 'compiled_model_dir_name')

    os_makedirs_mock.assert_called_with(expected_model_dir_path, exist_ok=ANY)
    launch_training_container_mock.assert_called_with(trainer, expected_model_dir_path)
    launch_gppi_validation_container_mock.assert_called_with(trainer, expected_model_dir_path)
