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
Helpers for working with `Trained Model Binary`
Docs: https://odahu.github.io/gen_glossary.html#term-trained-model-binary
"""
import json
import logging
import os
import uuid

import pydantic
import yaml
from odahuflow.sdk.gppi import entrypoint_invoke
from odahuflow.sdk.gppi.models import OdahuflowProjectManifest, OdahuflowProjectManifestBinaries
from odahuflow.sdk.io_proc_utils import run

_logger = logging.getLogger(__name__)

PROJECT_FILE = 'odahuflow.project.yaml'

VALIDATION_FAILED_EXCEPTION_MESSAGE = 'Exception raised while invoke model entrypoint. ' \
                                      'GPPI Self check is failed. ' \
                                      'Fix your model and try again.'


def _get_conda_bin_executable(executable_name):
    """
    Return path to the specified executable, assumed to be discoverable within the 'bin'
    subdirectory of a conda installation.
    """

    if "CONDA_EXE" in os.environ:
        conda_bin_dir = os.path.dirname(os.environ["CONDA_EXE"])
        executable = os.path.join(conda_bin_dir, executable_name)
    else:
        executable = executable_name
    return executable


def _get_activate_conda_command(conda_env_name):
    #  Checking for newer conda versions
    if 'CONDA_EXE' in os.environ:
        conda_path = _get_conda_bin_executable("conda")
        activate_conda_env = ['source ' + os.path.dirname(conda_path) +
                              '/../etc/profile.d/conda.sh']
        activate_conda_env += ["conda activate {0} 1>&2".format(conda_env_name)]
    else:
        activate_path = _get_conda_bin_executable("activate")
        # in case os name is not 'nt', we are not running on windows. It introduces
        # bash command otherwise.
        if os.name != "nt":
            activate_conda_env = ["source %s %s 1>&2" % (activate_path, conda_env_name)]
        else:
            activate_conda_env = ["conda %s %s 1>&2" % (activate_path, conda_env_name)]
    return ' && '.join(activate_conda_env)


def _check_conda():
    conda_path = _get_conda_bin_executable("conda")
    try:
        run(conda_path, "--help", stream_output=False)
    except Exception as exc_info:
        raise Exception("Could not find Conda executable at {0}. "
                        "Ensure Conda is installed as per the instructions "
                        "at https://conda.io/docs/user-guide/install/index.html") from exc_info


def load_odahuflow_project_manifest(manifest_path) -> OdahuflowProjectManifest:
    """
    Extract model manifest from file to object

    :param model: path to unpacked model (folder)
    :type model: str
    :return: None
    """
    manifest_file = os.path.join(manifest_path, PROJECT_FILE)
    if not manifest_file:
        raise Exception(f'Can not find manifest file {manifest_file}')

    with open(manifest_file, 'r') as manifest_stream:
        data = manifest_stream.read()

        try:
            data = json.loads(data)
        except json.JSONDecodeError:
            try:
                data = yaml.safe_load(data)
            except yaml.YAMLError as decode_error:
                raise ValueError(f'Cannot decode ModelPacking resource file: {decode_error}')

        try:
            return OdahuflowProjectManifest(**data)
        except pydantic.ValidationError as valid_error:
            raise Exception(f'Legion manifest file is in incorrect format: {valid_error}')


class ExecutionEnvironment:

    def __init__(self, name: str, binaries: OdahuflowProjectManifestBinaries, manifest_path: str, skip_deps=False):

        self.check_environment_executables()

        if not name:
            name = str(uuid.uuid4())

        self._name = name
        self._binaries = binaries
        self._manifest_path = manifest_path

        self.create_env()

        if skip_deps:
            _logger.warning(f'Flag "skip_deps"=True. Installing deps is skipped')
        else:
            _logger.info(f'Start to install dependencies for {self.__class__}')
            self.install_dependencies()

    def __str__(self):
        return f'{self.__class__}: name: {self._name}'

    @property
    def binaries(self):
        return self._binaries

    def execute(self, command: str, cwd: str = None, stream_output: bool = True):
        raise NotImplementedError

    def install_dependencies(self):
        raise NotImplementedError

    def create_env(self):
        raise NotImplementedError

    def check_environment_executables(self):
        raise NotImplementedError


class CondaExecutionEnvironment(ExecutionEnvironment):

    def check_environment_executables(self):
        _check_conda()
        _logger.info('Conda environment executables are detected')

    def env_created(self):
        conda_exec = _get_conda_bin_executable('conda')
        _1, stdout, _2 = run(conda_exec, 'env', 'list', '--json', stream_output=False)
        env_names = [os.path.basename(env) for env in json.loads(stdout)['envs']]
        return self._name in env_names

    def create_env(self):
        env_id = self._name
        conda_exec = _get_conda_bin_executable('conda')

        if self.env_created():
            _logger.info(f'Conda env with name {env_id!r} already exists. Creating step is skipped')
        else:
            _logger.info(f'Creating conda env with name {env_id!r}')
            run(conda_exec, 'create', '--yes', '--name', env_id, stream_output=False)

    def install_dependencies(self):

        env_id = self._name
        conda_exec = _get_conda_bin_executable('conda')

        # Install requirements from dep. list
        conda_dep_list = os.path.join(self._manifest_path, self.binaries.conda_path)
        _logger.info(f'Installing mandatory requirements from {conda_dep_list} to {env_id!r}')
        run(conda_exec, 'env', 'update', f'--name={env_id}', f'--file={conda_dep_list}', stream_output=False)

    def execute(self, command: str, cwd: str = None, stream_output: bool = True):
        activate_cmd = _get_activate_conda_command(self._name)
        _logger.info(f'cwd: {cwd}')
        return run('bash', '-c', f'{activate_cmd} && {command}', cwd=cwd, stream_output=False)


class GPPITrainedModelBinary:
    """
    Class that allows interact with TrainedModelBinary stored in GPPI interface
        * Validate `TrainedModelBinary`
        * Install `TrainedModelBinary` dependencies
        * Run self-check `TrainedModelBinary`

    Additional Info:
        * https://odahu.github.io/gen_glossary.html#term-general-python-prediction-interface
    """

    def __init__(self, manifest_path: str, use_current_env: bool = False, env_name: str = '', skip_deps=False):
        """

        :param manifest_path: path to dir where odahuflow.project.yaml
        :param use_current_env: if False, new python env will be created and all commands
        """
        self.manifest_path: str = manifest_path
        try:
            self.manifest: OdahuflowProjectManifest = load_odahuflow_project_manifest(manifest_path)
        except Exception as e:
            raise Exception(VALIDATION_FAILED_EXCEPTION_MESSAGE) from e
        self.use_current_env = use_current_env
        self.skip_deps = skip_deps

        self.exec_env = None
        if not use_current_env:
            _logger.info(f'Start initializing environment')
            self.exec_env: ExecutionEnvironment = self._init_exec_env(env_name)

    @property
    def model_dir(self):
        return os.path.join(self.manifest_path, self.manifest.model.workDir)

    def execute(self, command: str, cwd: str = None, stream_output: bool = True):
        if self.use_current_env:
            _logger.info(f'use_current_env flag = True. Start to execute command without changing environment')
            exit_code, output, err_ = run('bash', '-c', command, cwd=cwd, stream_output=stream_output)
        else:
            _logger.info(f'use_current_env flag = False. Start to execute command in {self.exec_env}')
            exit_code, output, err_ = self.exec_env.execute(command, cwd=cwd, stream_output=stream_output)

        _logger.info('Subprocess result: %s\n\nSTDOUT:\n%s\n\nSTDERR:%s\n ======' %
                     (exit_code, output, err_))

        return output

    def self_check(self):
        """
        Try to invoke public API of GPPI
        1) Invoke entrypoint.init()
        2) Invoke entrypoint.info()
        3) If input samples is found in package run entrypoint.predict_on_matrix()
        :return:
        """

        _logger.info('Start procedure of self checking for GPPI')
        try:
            self.execute(
                f'python  {entrypoint_invoke.__file__} -v --model {self.model_dir} '
                f'--entrypoint {self.manifest.model.entrypoint} self_check',
                stream_output=False
            )
        except Exception as e:
            raise Exception(VALIDATION_FAILED_EXCEPTION_MESSAGE) from e

        _logger.info('Self check is successful. GPPI model is packed as expected.')

    def predict(self, input_file: str, output_dir: str, output_file_name):
        """
        Do predictions
        :param input_file: JSON file with data for predictions
        :param output_dir: output directory where results will be saved
        :return:
        """
        try:
            self.execute(
                f'python  {entrypoint_invoke.__file__} -v --model {self.model_dir} '
                f'--entrypoint {self.manifest.model.entrypoint} predict {input_file} {output_dir} '
                f'{"--output_file_name "+output_file_name if output_file_name else ""}',
                stream_output=False
            )
        except Exception as e:
            raise Exception(VALIDATION_FAILED_EXCEPTION_MESSAGE) from e

        _logger.info('Prediction successfully completed')

    def info(self) -> str:
        """
        Show input/output model schema
        :return:
        """
        try:
            out = self.execute(
                f'python  {entrypoint_invoke.__file__} -v --model {self.model_dir} '
                f'--entrypoint {self.manifest.model.entrypoint} info',
                stream_output=False
            )
        except Exception as e:
            raise Exception(VALIDATION_FAILED_EXCEPTION_MESSAGE) from e

        return out

    def _init_exec_env(self, env_name: str) -> ExecutionEnvironment:
        if self.manifest.binaries.dependencies == 'conda':
            env = CondaExecutionEnvironment(env_name, self.manifest.binaries, self.manifest_path, self.skip_deps)
        else:
            raise RuntimeError(f'{self.manifest.binaries.dependencies} dependency manager is not recognized')

        return env
