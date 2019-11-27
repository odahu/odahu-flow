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
import click

from odahuflow.sdk.gppi.executor import GPPITrainedModelBinary


@click.group()
def gppi():
    """
    Allow you to perform actions on odahuflow gppi models
    """
    pass


@gppi.command()
@click.argument('path', type=click.Path(exists=True, file_okay=False))
@click.option('--use-current-env/--not-use-current-env', 'use_current_env',
              help='Use current environment', default=False)
@click.option('--env-name', '-e', 'env_name', help='Environment name to run GPPI model', type=click.STRING)
@click.option('--skip-deps/--not-skip-deps', 'skip_deps', help='Not install dependencies', default=False)
def test(path: str, use_current_env: bool = False, env_name: str = '', skip_deps: bool = False):
    """
    Initialize GPPI model and try to execute api

    Next API are tested:

    Executes .info() method

    Executes .predict_on_matrix() method
    with data deserialized from head_input.pkl as a parameter
    (only if head_input.pkl exists as a model artifact)
    """
    model = GPPITrainedModelBinary(path, use_current_env, env_name, skip_deps)
    model.self_check()
    click.echo(f'OK\nGPPI is correct and could be packaged or deployed')
