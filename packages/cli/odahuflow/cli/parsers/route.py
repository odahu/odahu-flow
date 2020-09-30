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
from click import pass_obj

from odahuflow.cli.utils import click_utils
from odahuflow.cli.utils.error_handler import check_id_or_file_params_present, TIMEOUT_ERROR_MESSAGE, \
    IGNORE_NOT_FOUND_ERROR_MESSAGE
from odahuflow.cli.utils.output import DEFAULT_OUTPUT_FORMAT, format_output, validate_output_format
from odahuflow.sdk import config
from odahuflow.sdk.clients.api import WrongHttpStatusCode
from odahuflow.sdk.clients.api_aggregated import parse_resources_file_with_one_item
from odahuflow.sdk.clients.route import ModelRoute, ModelRouteClient, READY_STATE

DEFAULT_WAIT_TIMEOUT = 5
LOGGER = logging.getLogger(__name__)


@click.group(cls=click_utils.BetterHelpGroup)
@click.option('--url', help='API server host', default=config.API_URL)
@click.option('--token', help='API server jwt token', default=config.API_TOKEN)
@click.pass_context
def route(ctx: click.core.Context, url: str, token: str):
    """
    Allow you to perform actions on routes
    """
    ctx.obj = ModelRouteClient(url, token)


@route.command()
@click.option('--mr-id', '--id', 'mr_id', help='ModelRoute ID')
@click.option('--output-format', '-o', 'output_format', help='Output format  [json|table|yaml|jsonpath]',
              default=DEFAULT_OUTPUT_FORMAT, callback=validate_output_format)
@pass_obj
def get(client: ModelRouteClient, mr_id: str, output_format: str):
    """
    Get routes.\n
    The command without id argument retrieve all routes.\n
    Get all routes in json format:\n
        odahuflowctl route get --output-format json\n
    Get model route with "git-repo" id:\n
        odahuflowctl route get --id git-repo\n
    Using jsonpath:\n
        odahuflowctl route get -o 'jsonpath=[*].spec.reference'
    \f
    :param client: ModelRoute HTTP client
    :param mr_id: ModelRoute ID
    :param output_format: Output format
    :return:
    """
    routes = [client.get(mr_id)] if mr_id else client.get_all()

    format_output(routes, output_format)


@route.command()
@click.option('--mr-id', '--id', 'mr_id', help='ModelRoute ID')
@click.option('--file', '-f', type=click.Path(), required=True, help='Path to the file with model route')
@click.option('--wait/--no-wait', default=True,
              help='no wait until scale will be finished')
@click.option('--timeout', default=600, type=int,
              help='timeout in seconds. for wait (if no-wait is off)')
@pass_obj
def create(client: ModelRouteClient, mr_id: str, file: str, wait: bool, timeout: int):
    """
    Create a model route.\n
    You should specify a path to file with a model route. The file must contain only one model route.
    For now, CLI supports yaml and JSON file formats.
    If you want to create multiples routes then you should use "odahuflowctl bulk apply" instead.
    If you provide the model route id parameter than it will be overridden before sending to API server.\n
    Usage example:\n
        * odahuflowctl route create -f route.yaml --id examples-git
    \f
    :param timeout: timeout in seconds. for wait (if no-wait is off)
    :param wait: no wait until operation will be finished
    :param client: ModelRoute HTTP client
    :param mr_id: ModelRoute ID
    :param file: Path to the file with only one model route
    """
    route_resource = parse_resources_file_with_one_item(file).resource
    if not isinstance(route_resource, ModelRoute):
        raise ValueError(f'ModelRoute expected, but {type(route_resource)} provided')

    if mr_id:
        route_resource.id = mr_id

    click.echo(client.create(route_resource))

    wait_operation_finish(timeout, wait, mr_id, client)


@route.command()
@click.option('--route-id', '--id', 'mr_id', help='ModelRoute ID')
@click.option('--file', '-f', type=click.Path(), required=True, help='Path to the file with model route')
@click.option('--wait/--no-wait', default=True,
              help='no wait until scale will be finished')
@click.option('--timeout', default=600, type=int,
              help='timeout in seconds. for wait (if no-wait is off)')
@pass_obj
def edit(client: ModelRouteClient, mr_id: str, file: str, wait: bool, timeout: int):
    """
    Update a model route.\n
    You should specify a path to file with a model route. The file must contain only one model route.
    For now, CLI supports yaml and JSON file formats.
    If you want to update multiples routes then you should use "odahuflowctl bulk apply" instead.
    If you provide the model route id parameter than it will be overridden before sending to API server.\n
    Usage example:\n
        * odahuflowctl route update -f route.yaml --id examples-git
    \f
    :param timeout: timeout in seconds. for wait (if no-wait is off)
    :param wait: no wait until operation will be finished
    :param client: Model route HTTP client
    :param mr_id: Model route ID
    :param file: Path to the file with only one model route
    """
    route_resource = parse_resources_file_with_one_item(file).resource
    if not isinstance(route_resource, ModelRoute):
        raise ValueError(f'ModelRoute expected, but {type(route_resource)} provided')

    if mr_id:
        route_resource.id = mr_id

    click.echo(client.edit(route_resource))

    wait_operation_finish(timeout, wait, mr_id, client)


@route.command()
@click.option('--route-id', '--id', 'mr_id', help='ModelRoute ID')
@click.option('--file', '-f', type=click.Path(), help='Path to the file with model route')
@click.option('--ignore-not-found/--not-ignore-not-found', default=False,
              help='ignore if Model Deployment is not found')
@pass_obj
def delete(client: ModelRouteClient, mr_id: str, file: str, ignore_not_found: bool):
    """
    Delete a model route.\n
    For this command, you must provide a model route ID or path to file with one model route.
    The file must contain only one model route.
    If you want to delete multiples routes then you should use "odahuflowctl bulk delete" instead.
    For now, CLI supports yaml and JSON file formats.
    The command will be failed if you provide both arguments.\n
    Usage example:\n
        * odahuflowctl route delete --id examples-git\n
        * odahuflowctl route delete -f route.yaml
    \f
    :param client: ModelRoute HTTP client
    :param mr_id: ModelRoute ID
    :param file: Path to the file with only one model route
    :param ignore_not_found: ignore if Model Deployment is not found
    """
    check_id_or_file_params_present(mr_id, file)

    if file:
        route_resource = parse_resources_file_with_one_item(file).resource
        if not isinstance(route_resource, ModelRoute):
            raise ValueError(f'ModelRoute expected, but {type(route_resource)} provided')

        mr_id = route_resource.id

    try:
        click.echo(client.delete(mr_id))
    except WrongHttpStatusCode as e:
        if e.status_code != 404 or not ignore_not_found:
            raise e

        click.echo(IGNORE_NOT_FOUND_ERROR_MESSAGE.format(kind=ModelRoute.__name__, id=mr_id))


def wait_operation_finish(timeout: int, wait: bool, mr_id: str, mr_client: ModelRouteClient):
    """
    Wait route to finish according command line arguments

    :param timeout: timeout in seconds. for wait (if no-wait is off)
    :param wait: no wait until operation will be finished
    :param mr_id: Model Route id
    :param mr_client: Model Route Client

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
            mr = mr_client.get(mr_id)
            if mr.status.state == READY_STATE:
                print(f'Model Route {mr_id} is ready')
                return
            elif mr.status.state == "":
                print(f"Can't determine the state of {mr.id}. Sleeping...")
            else:
                print(f'Current route state is {mr.status.state}. Sleeping...')
        except WrongHttpStatusCode:
            LOGGER.info('Callback have not confirmed completion of the operation')

        LOGGER.debug('Sleep before next request')
        time.sleep(DEFAULT_WAIT_TIMEOUT)
