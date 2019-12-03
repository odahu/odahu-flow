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
import os
import sys

import click
from odahuflow.sdk.gppi.executor import GPPITrainedModelBinary

ODAHUFLOW_GPPI_MODEL_PATH_ENV_NAME = 'ODAHUFLOW_GPPI_MODEL_PATH'
ODAHUFLOW_CONDA_ENV_NAME = 'ODAHUFLOW_CONDA'


class GppiCommandContextObj:

    def __init__(self, model_binary: GPPITrainedModelBinary):
        self.model_binary = model_binary


@click.group()
@click.option('--gppi-model-path', '-m', type=click.Path(exists=True, file_okay=False),
              envvar=ODAHUFLOW_GPPI_MODEL_PATH_ENV_NAME)
@click.option('--use-current-env/--not-use-current-env', 'use_current_env',
              help='Use current environment', default=False)
@click.option('--env-name', '-e', 'env_name', help='Environment name to run GPPI model', type=click.STRING,
              envvar=ODAHUFLOW_CONDA_ENV_NAME)
@click.option('--skip-deps/--not-skip-deps', 'skip_deps', help='Not install dependencies', default=False)
@click.pass_context
def gppi(ctx, gppi_model_path: str, use_current_env: bool = False, env_name: str = '', skip_deps: bool = False):
    """
    Allow you to perform actions on odahuflow gppi models
    """

    if not gppi_model_path:
        click.echo(f'--gppi-model-path OR ${ODAHUFLOW_GPPI_MODEL_PATH_ENV_NAME} env var must be provided')
        sys.exit(1)

    mb = GPPITrainedModelBinary(gppi_model_path, use_current_env, env_name, skip_deps)
    ctx.obj = GppiCommandContextObj(mb)


@gppi.command()
@click.pass_context
def test(ctx):
    """
    Initialize GPPI model and try to execute api

    Next API are tested:

    Executes .info() method

    Executes .predict_on_matrix() method
    with data deserialized from head_input.pkl as a parameter
    (only if head_input.pkl exists as a model artifact)
    """
    ctx_obj: GppiCommandContextObj = ctx.obj
    ctx_obj.model_binary.self_check()
    click.echo(f'OK\nGPPI is correct and could be packaged or deployed')


@gppi.command()
@click.argument('input-file', type=click.Path(exists=True, dir_okay=False))
@click.argument('output-dir', type=click.Path(exists=True, file_okay=False))
@click.option('--output-file-name', '-o', default='results.json')
@click.pass_context
def predict(ctx, input_file: str, output_dir: str, output_file_name: str):
    """

    :param ctx:
    :param input_file: Input JSON file for predictions
    :param output_dir: Output directory where results should be saved
    :param output_file_name: Output filename with predictions
    :return:
    """
    ctx_obj: GppiCommandContextObj = ctx.obj
    ctx_obj.model_binary.predict(input_file, output_dir, output_file_name)
    full_path = os.path.join(output_dir, output_file_name)
    click.echo(f'Prediction file: {full_path}')
