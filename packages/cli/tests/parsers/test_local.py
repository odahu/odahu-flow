import json
import os
import pathlib
import tempfile
from typing import List
from unittest.mock import patch, Mock

from click.testing import CliRunner, Result
from docker.types import Mount
from pytest_mock import MockFixture

from odahuflow.cli.parsers.local import training
from odahuflow.sdk.gppi.executor import PROJECT_FILE
from odahuflow.sdk.local import training as training_sdk
from odahuflow.sdk.local.training import MODEL_OUTPUT_CONTAINER_PATH
from odahuflow.sdk.models import ModelTraining, ToolchainIntegration, ModelTrainingSpec


def test_list_local_trainings(tmpdir):
    folders = ['wine-1@-12',
               '&wine-1-12',
               '[wine-1@-12',
               '2Wine-1@-12',
               'Wine-1@-123',
               '@wine-1@-12',
               'zine-1@-12',
               'Awine-1@']
    for folder in folders:
        tmpdir.mkdir(folder)
        pathlib.Path(tmpdir, folder, PROJECT_FILE).touch()
    with patch('odahuflow.sdk.local.training.config.LOCAL_MODEL_OUTPUT_DIR', tmpdir):
        assert training_sdk.list_local_trainings() == ['&wine-1-12',
                                          '2Wine-1@-12',
                                          '@wine-1@-12',
                                          'Awine-1@',
                                          'Wine-1@-123',
                                          '[wine-1@-12',
                                          'wine-1@-12',
                                          'zine-1@-12']


def test_local_training_relative_output_dir(mocker: MockFixture, cli_runner: CliRunner):
    """
    Tests issue #208 - Converting relative path to output-dir to absolute,
    because Docker requires mount paths to be absolute
    :param mocker: mocker fixture
    :param cli_runner: Click runner fixture
    """

    trainer_mock = mocker.patch.object(training, 'K8sTrainer', autospec=True).return_value
    trainer_mock.model_training.spec.model.artifact_name_template = 'model_dir_template'
    api_client = Mock()

    mocker.patch.object(training_sdk, 'create_mt_config_file')
    mocker.patch.object(training_sdk, 'stream_container_logs')
    mocker.patch.object(training_sdk, 'raise_error_if_container_failed')

    docker_mock: Mock = mocker.patch.object(training_sdk.docker, 'from_env').return_value
    docker_mock.api.inspect_container = Mock(return_value={})

    with tempfile.TemporaryDirectory(dir=os.curdir) as temp_dir:
        temp_dir_path = pathlib.Path(temp_dir)
        training_yml = temp_dir_path / 'training.yml'
        training_yml.write_text(json.dumps(
            {**ModelTraining(id='training1', spec=ModelTrainingSpec(toolchain='toolchain1')).to_dict(),
             **{'kind': 'ModelTraining'}}))

        toolchain_yml = temp_dir_path / 'toolchain.yml'
        toolchain_yml.write_text(json.dumps({**ToolchainIntegration(id='toolchain1').to_dict(),
                                             **{'kind': 'ToolchainIntegration'}}))

        temp_dir_relative_path: str = os.path.relpath(temp_dir_path)

        result: Result = cli_runner.invoke(
            training.training_group,
            ['run', '--output-dir', temp_dir_relative_path, '--manifest-file', str(training_yml),
             '--manifest-file', str(toolchain_yml), '--id', 'training1'],
            obj=api_client)

        assert result.exit_code == 0, f'command invocation ended with exit code {result.exit_code}'
        assert result.exception is None, f'command invocation ended with exception {result.exception}'

        abs_path = os.path.abspath(temp_dir)

        for call in docker_mock.containers.run.call_args_list:
            call_kwargs = call[1]
            mounts: List[Mount] = call_kwargs.get('mounts', [])

            for mount in filter(lambda m: m.get('Target') == MODEL_OUTPUT_CONTAINER_PATH, mounts):
                assert os.path.dirname(mount.get('Source')) == abs_path
