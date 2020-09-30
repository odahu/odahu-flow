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
API commands for odahuflow cli
"""
import logging
import time

import click

from odahuflow.cli.utils import click_utils
from odahuflow.cli.utils.client import pass_obj
from odahuflow.cli.utils.error_handler import check_id_or_file_params_present, TIMEOUT_ERROR_MESSAGE, \
    IGNORE_NOT_FOUND_ERROR_MESSAGE
from odahuflow.cli.utils.output import DEFAULT_OUTPUT_FORMAT, format_output, validate_output_format
from odahuflow.cli.utils.verifiers import positive_number
from odahuflow.sdk import config
from odahuflow.sdk.clients.api import EntityAlreadyExists, WrongHttpStatusCode
from odahuflow.sdk.clients.api_aggregated import parse_resources_file_with_one_item
from odahuflow.sdk.clients.deployment import ModelDeployment, ModelDeploymentClient, READY_STATE, \
    FAILED_STATE

DEFAULT_WAIT_TIMEOUT = 5
# 20 minutes
DEFAULT_DEPLOYMENT_TIMEOUT = 20 * 60


LOGGER = logging.getLogger(__name__)


@click.group(cls=click_utils.BetterHelpGroup)
@click.option('--url', help='API server host', default=config.API_URL)
@click.option('--token', help='API server jwt token', default=config.API_TOKEN)
@click.pass_context
def deployment(ctx: click.core.Context, url: str, token: str):
    """
    Allow you to perform actions on deployments.\n
    Alias for the command is dep.
    """
    ctx.obj = ModelDeploymentClient(url, token)


@deployment.command()
@click.option('--md-id', '--id', help='Model deployment ID')
@click.option('--output-format', '-o', 'output_format', help='Output format  [json|table|yaml|jsonpath]',
              default=DEFAULT_OUTPUT_FORMAT, callback=validate_output_format)
@pass_obj
def get(client: ModelDeploymentClient, md_id: str, output_format: str):
    """
    \b
    Get deployments.
    The command without id argument retrieve all deployments.
    \b
    Get all deployments in json format:
        odahuflowctl dep get --output-format json
    \b
    Get deployment with "git-repo" id:
        odahuflowctl dep get --id model-wine
    \b
    Using jsonpath:
        odahuflowctl dep get -o 'jsonpath=[*].spec.reference'
    \f
    :param client: Model deployment HTTP client
    :param md_id: Model deployment ID
    :param output_format: Output format
    :return:
    """
    mds = [client.get(md_id)] if md_id else client.get_all()

    format_output(mds, output_format)


@deployment.command()
@click.option('--md-id', '--id', help='Model deployment ID')
@click.option('--file', '-f', type=click.Path(), required=True, help='Path to the file with deployment')
@click.option('--wait/--no-wait', default=True,
              help='no wait until scale will be finished')
@click.option('--timeout', default=DEFAULT_DEPLOYMENT_TIMEOUT, type=int, callback=positive_number,
              help='timeout in seconds. for wait (if no-wait is off)')
@click.option('--image', type=str, help='Override Docker image from file')
@click.option('--ignore-if-exists', is_flag=True,
              help='Ignore if entity is already exists on API server. Return success status code')
@pass_obj
def create(client: ModelDeploymentClient, md_id: str, file: str, wait: bool, timeout: int, image: str,
           ignore_if_exists: bool):
    """
    \b
    Create a deployment.
    You should specify a path to file with a deployment. The file must contain only one deployment.
    For now, CLI supports YAML and JSON file formats.
    If you want to create multiple deployments, you should use "odahuflowctl bulk apply" instead.
    If you provide the deployment id parameter, it will override before sending to API server.
    \b
    Usage example:
        * odahuflowctl dep create -f dep.yaml --id examples-git
    \f
    :param timeout: timeout in seconds. for wait (if no-wait is off)
    :param wait: no wait until deployment will be finished
    :param client: Model deployment HTTP client
    :param md_id: Model deployment ID
    :param file: Path to the file with only one deployment
    :param image: Override Docker image from file
    :param ignore_if_exists: Return success status code if entity is already exists
    """
    md = parse_resources_file_with_one_item(file).resource
    if not isinstance(md, ModelDeployment):
        raise ValueError(f'Model deployment expected, but {type(md)} provided')

    if md_id:
        md.id = md_id

    if image:
        md.spec.image = image

    try:
        res = client.create(md)
    except EntityAlreadyExists as e:
        if ignore_if_exists:
            LOGGER.debug(f'--ignore-if-exists was passed: {e} will be suppressed')
            click.echo('Deployment already exists')
            return
        raise

    click.echo(res)

    wait_deployment_finish(timeout, wait, md.id, client)


@deployment.command()
@click.option('--md-id', '--id', help='Model deployment ID')
@click.option('--file', '-f', type=click.Path(), required=True, help='Path to the file with deployment')
@click.option('--wait/--no-wait', default=True,
              help='no wait until scale will be finished')
@click.option('--timeout', default=DEFAULT_DEPLOYMENT_TIMEOUT, type=int, callback=positive_number,
              help='timeout in seconds. for wait (if no-wait is off)')
@click.option('--image', type=str, help='Override Docker image from file')
@pass_obj
def edit(client: ModelDeploymentClient, md_id: str, file: str, wait: bool, timeout: int, image: str):
    """
    \b
    Update a deployment.
    You should specify a path to file with a deployment. The file must contain only one deployment.
    For now, CLI supports YAML and JSON file formats.
    If you want to update multiple deployments, you should use "odahuflowctl bulk apply" instead.
    If you provide the deployment id parameter, it will override before sending to API server.
    \b
    Usage example:
        * odahuflowctl dep update -f dep.yaml --id examples-git
    \f
    :param client: Model deployment HTTP client
    :param md_id: Model deployment ID
    :param file: Path to the file with only one deployment
    :param timeout: timeout in seconds. for wait (if no-wait is off)
    :param wait: no wait until edit will be finished
    :param image: Override Docker image from file
    """
    md = parse_resources_file_with_one_item(file).resource
    if not isinstance(md, ModelDeployment):
        raise ValueError(f'Model deployment expected, but {type(md)} provided')

    if md_id:
        md.id = md_id

    if image:
        md.spec.image = image

    click.echo(client.edit(md))

    wait_deployment_finish(timeout, wait, md.id, client)


@deployment.command()
@click.option('--md-id', '--id', help='Model deployment ID')
@click.option('--file', '-f', type=click.Path(), help='Path to the file with deployment')
@click.option('--wait/--no-wait', default=True,
              help='no wait until scale will be finished')
@click.option('--timeout', default=DEFAULT_DEPLOYMENT_TIMEOUT, type=int, callback=positive_number,
              help='timeout in seconds. for wait (if no-wait is off)')
@click.option('--ignore-not-found/--not-ignore-not-found', default=False,
              help='ignore if Model Deployment is not found')
@pass_obj
def delete(client: ModelDeploymentClient, md_id: str, file: str, ignore_not_found: bool,
           wait: bool, timeout: int):
    """
    \b
    Delete a deployment.
    For this command, you must provide a deployment ID or path to file with one deployment.
    The file must contain only one deployment.
    If you want to delete multiple deployments, you should use "odahuflowctl bulk delete" instead.
    For now, CLI supports YAML and JSON file formats.
    The command will fail if you provide both arguments.
    \b
    Usage example:
        * odahuflowctl dep delete --id examples-git
        * odahuflowctl dep delete -f dep.yaml
    \f
    :param timeout: timeout in seconds. for wait (if no-wait is off)
    :param wait: no wait until deletion will be finished
    :param client: Model deployment HTTP client
    :param md_id: Model deployment ID
    :param file: Path to the file with only one deployment
    :param ignore_not_found: ignore if Model Deployment is not found
    """
    check_id_or_file_params_present(md_id, file)

    if file:
        md = parse_resources_file_with_one_item(file).resource
        if not isinstance(md, ModelDeployment):
            raise ValueError(f'Model deployment expected, but {type(md)} provided')

        md_id = md.id

    try:
        message = client.delete(md_id)

        wait_delete_operation_finish(timeout, wait, md_id, client)
        click.echo(message)
    except WrongHttpStatusCode as e:
        if e.status_code != 404 or not ignore_not_found:
            raise e

        click.echo(IGNORE_NOT_FOUND_ERROR_MESSAGE.format(kind=ModelDeployment.__name__, id=md_id))


def wait_delete_operation_finish(timeout: int, wait: bool, md_id: str, md_client: ModelDeploymentClient):
    """
    Wait delete operation

    :param timeout: timeout in seconds. for wait (if no-wait is off)
    :param wait: no wait until deletion will be finished
    :param md_id: Model Deployment name
    :param md_client: Model Deployment Client

    :return: None
    """
    if not wait:
        return

    start = time.time()
    if timeout <= 0:
        raise Exception('Invalid --timeout argument: should be positive integer')

    while True:
        elapsed = time.time() - start
        if elapsed > timeout:
            raise Exception('Time out: operation has not been confirmed')

        try:
            md_client.get(md_id)
        except WrongHttpStatusCode as e:
            if e.status_code == 404:
                return
            LOGGER.info('Callback have not confirmed completion of the operation')

        print(f'Model deployment {md_id} is still being deleted...')
        time.sleep(DEFAULT_WAIT_TIMEOUT)


def wait_deployment_finish(timeout: int, wait: bool, md_id: str, md_client: ModelDeploymentClient):
    """
    Wait for deployment to finish according to command line arguments

    :param timeout: timeout in seconds. for wait (if no-wait is off)
    :param wait: no wait until deletion will be finished
    :param md_id: Model Deployment name
    :param md_client: Model Deployment Client

    :return: None
    """
    if not wait:
        return

    start = time.time()
    if timeout <= 0:
        raise Exception('Invalid --timeout argument: should be positive integer')

    while True:
        elapsed = time.time() - start
        if elapsed > timeout:
            raise Exception(TIMEOUT_ERROR_MESSAGE)

        try:
            md: ModelDeployment = md_client.get(md_id)
            if md.status.state == READY_STATE:
                if md.spec.min_replicas <= md.status.available_replicas:
                    print(f'Model {md_id} was deployed. '
                          f'Deployment process took {round(time.time() - start)} seconds')
                    return
                else:
                    print(f'Model {md_id} was deployed. '
                          f'Number of available pods is {md.status.available_replicas}/{md.spec.min_replicas}')
            elif md.status.state == FAILED_STATE:
                raise Exception(f'Model deployment {md_id} was failed')
            elif md.status.state == "":
                print(f"Can't determine the state of {md.id}. Sleeping...")
            else:
                print(f'Current deployment state is {md.status.state}. Sleeping...')
        except WrongHttpStatusCode:
            LOGGER.info('Callback have not confirmed completion of the operation')

        LOGGER.debug('Sleep before next request')
        time.sleep(DEFAULT_WAIT_TIMEOUT)
