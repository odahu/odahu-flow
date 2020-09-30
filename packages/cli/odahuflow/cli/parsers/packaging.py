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
Packing commands for odahuflow cli
"""
import logging
import time
from http.client import HTTPException

import click
from requests import RequestException

from odahuflow.cli.utils import click_utils
from odahuflow.cli.utils.client import pass_obj
from odahuflow.cli.utils.error_handler import check_id_or_file_params_present, TIMEOUT_ERROR_MESSAGE, \
    IGNORE_NOT_FOUND_ERROR_MESSAGE
from odahuflow.cli.utils.logs import print_logs
from odahuflow.cli.utils.output import format_output, DEFAULT_OUTPUT_FORMAT, validate_output_format
from odahuflow.sdk import config
from odahuflow.sdk.clients.api import EntityAlreadyExists, WrongHttpStatusCode, APIConnectionException
from odahuflow.sdk.clients.api_aggregated import parse_resources_file_with_one_item
from odahuflow.sdk.clients.packaging import ModelPackaging, ModelPackagingClient, SUCCEEDED_STATE, FAILED_STATE

DEFAULT_WAIT_TIMEOUT = 3
# 1 hour
DEFAULT_PACKAGING_TIMEOUT = 60 * 60
LOG_READ_TIMEOUT_SECONDS = 60

LOGGER = logging.getLogger(__name__)


@click.group(cls=click_utils.BetterHelpGroup)
@click.option('--url', help='API server host', default=config.API_URL)
@click.option('--token', help='API server jwt token', default=config.API_TOKEN)
@click.pass_context
def packaging(ctx: click.core.Context, url: str, token: str):
    """
    Allow you to perform actions on packagings.\n
    Alias for the command is pack.
    """
    ctx.obj = ModelPackagingClient(url, token)


@packaging.command()
@click.option('--pack-id', '--id', help='Model packaging ID')
@click.option('--output-format', '-o', 'output_format', help='Output format  [json|table|yaml|jsonpath]',
              default=DEFAULT_OUTPUT_FORMAT, callback=validate_output_format)
@pass_obj
def get(client: ModelPackagingClient, pack_id: str, output_format: str):
    """
    \b
    Get packagings.
    The command without id argument retrieve all packagings.
    \b
    Get all packagings in json format:
        odahuflowctl pack get --output-format json
    \b
    Get packaging with "git-repo" id:
        odahuflowctl pack get --id git-repo
    \b
    Using jsonpath:
        odahuflowctl pack get -o 'jsonpath=[*].spec.reference'
    \f
    :param client: Model packaging HTTP client
    :param pack_id: Model packaging ID
    :param output_format: Output format
    :return:
    """
    packs = [client.get(pack_id)] if pack_id else client.get_all()

    format_output(packs, output_format)


@packaging.command()
@click.option('--pack-id', '--id', help='Model packaging ID')
@click.option('--file', '-f', type=click.Path(), required=True, help='Path to the file with packaging')
@click.option('--wait/--no-wait', default=True,
              help='no wait until scale will be finished')
@click.option('--artifact-name', type=str, help='Override artifact name from file')
@click.option('--timeout', default=DEFAULT_PACKAGING_TIMEOUT, type=int,
              help='timeout in seconds. for wait (if no-wait is off)')
@click.option('--ignore-if-exists', is_flag=True,
              help='Ignore if entity is already exists on API server. Return success status code')
@pass_obj
def create(client: ModelPackagingClient, pack_id: str, file: str, wait: bool, timeout: int,
           artifact_name: str, ignore_if_exists: bool):
    """
    \b
    Create a packaging.
    You should specify a path to file with a packaging. The file must contain only one packaging.
    For now, CLI supports YAML and JSON file formats.
    If you want to create multiple packagings, you should use "odahuflowctl bulk apply" instead.
    If you provide the packaging id parameter, it will override before sending to API server.
    \b
    Usage example:
        * odahuflowctl pack create -f pack.yaml --id examples-git
    \f
    :param timeout: timeout in seconds. for wait (if no-wait is off)
    :param wait: no wait until scale will be finished
    :param client: Model packaging HTTP client
    :param pack_id: Model packaging ID
    :param file: Path to the file with only one packaging
    :param artifact_name: Override artifact name from file
    :param ignore_if_exists: Return success status code if entity is already exists
    """
    pack = parse_resources_file_with_one_item(file).resource
    if not isinstance(pack, ModelPackaging):
        raise ValueError(f'Model packaging expected, but {type(pack)} provided')

    if pack_id:
        pack.id = pack_id

    if artifact_name:
        pack.spec.artifact_name = artifact_name

    try:
        mp = client.create(pack)
    except EntityAlreadyExists as e:
        if ignore_if_exists:
            LOGGER.debug(f'--ignore-if-exists was passed: {e} will be suppressed')
            click.echo('Packaging already exists')
            return
        raise

    click.echo(f"Start packing: {mp}")

    wait_packaging_finish(timeout, wait, pack.id, client)


@packaging.command()
@click.option('--pack-id', '--id', help='Model packaging ID')
@click.option('--file', '-f', type=click.Path(), required=True, help='Path to the file with packaging')
@click.option('--wait/--no-wait', default=True,
              help='no wait until scale will be finished')
@click.option('--artifact-name', type=str, help='Override artifact name from file')
@click.option('--timeout', default=DEFAULT_PACKAGING_TIMEOUT, type=int,
              help='timeout in seconds. for wait (if no-wait is off)')
@pass_obj
def edit(client: ModelPackagingClient, pack_id: str, file: str, wait: bool, timeout: int,
         artifact_name: str):
    """
    \b
    Update a packaging.
    You should specify a path to file with a packaging. The file must contain only one packaging.
    For now, CLI supports YAML and JSON file formats.
    If you want to update multiple packagings, you should use "odahuflowctl bulk apply" instead.
    If you provide the packaging id parameter, it will override before sending to API server.
    \b
    Usage example:
        * odahuflowctl pack update -f pack.yaml --id examples-git
    \f
    :param client: Model packaging HTTP client
    :param pack_id: Model packaging ID
    :param file: Path to the file with only one packaging
    :param timeout: timeout in seconds. for wait (if no-wait is off)
    :param wait: no wait until scale will be finished
    :param artifact_name: Override artifact name from file
    """
    pack = parse_resources_file_with_one_item(file).resource
    if not isinstance(pack, ModelPackaging):
        raise ValueError(f'Model packaging expected, but {type(pack)} provided')

    if pack_id:
        pack.id = pack_id

    if artifact_name:
        pack.spec.artifact_name = artifact_name

    mp = client.edit(pack)
    click.echo(f"Rerun packing: {mp}")

    wait_packaging_finish(timeout, wait, pack.id, client)


@packaging.command()
@click.option('--pack-id', '--id', help='Model packaging ID')
@click.option('--file', '-f', type=click.Path(), help='Path to the file with packaging')
@click.option('--ignore-not-found/--not-ignore-not-found', default=False,
              help='ignore if Model Packaging is not found')
@pass_obj
def delete(client: ModelPackagingClient, pack_id: str, file: str, ignore_not_found: bool):
    """
    \b
    Delete a packaging.
    For this command, you must provide a packaging ID or path to file with one packaging.
    The file must contain only one packaging.
    If you want to delete multiple packagings, you should use "odahuflowctl bulk delete" instead.
    For now, CLI supports YAML and JSON file formats.
    The command will fail if you provide both arguments.
    \b
    Usage example:
        * odahuflowctl pack delete --id examples-git
        * odahuflowctl pack delete -f pack.yaml
    \f
    :param client: Model packaging HTTP client
    :param pack_id: Model packaging ID
    :param file: Path to the file with only one packaging
    :param ignore_not_found: ignore if Model Packaging is not found
    """
    check_id_or_file_params_present(pack_id, file)

    if file:
        pack = parse_resources_file_with_one_item(file).resource
        if not isinstance(pack, ModelPackaging):
            raise ValueError(f'Model packaging expected, but {type(pack)} provided')

        pack_id = pack.id

    try:
        message = client.delete(pack_id)
        click.echo(message)
    except WrongHttpStatusCode as e:
        if e.status_code != 404 or not ignore_not_found:
            raise e

        click.echo(IGNORE_NOT_FOUND_ERROR_MESSAGE.format(kind=ModelPackaging.__name__, id=pack_id))


@packaging.command()
@click.option('--pack-id', '--id', help='Model packaging ID')
@click.option('--file', '-f', type=click.Path(), help='Path to the file with packaging')
@click.option('--follow/--not-follow', default=True,
              help='Follow logs stream')
@pass_obj
def logs(client: ModelPackagingClient, pack_id: str, file: str, follow: bool):
    """
    \b
    Stream packaging logs.
    For this command, you must provide a packaging ID or path to file with one packaging.
    The file must contain only one packaging.
    The command will fail if you provide both arguments.
    \b
    Usage example:
        * odahuflowctl pack delete --id examples-git
        * odahuflowctl pack delete -f pack.yaml
    \f
    :param follow: Follow logs stream
    :param client: Model packaging HTTP client
    :param pack_id: Model packaging ID
    :param file: Path to the file with only one packaging
    """
    check_id_or_file_params_present(pack_id, file)

    if file:
        pack = parse_resources_file_with_one_item(file).resource
        if not isinstance(pack, ModelPackaging):
            raise ValueError(f'Model packaging expected, but {type(pack)} provided')

        pack_id = pack.id

    for msg in client.log(pack_id, follow):
        print_logs(msg)


def wait_packaging_finish(timeout: int, wait: bool, mp_id: str, mp_client: ModelPackagingClient):
    """
    Wait for packaging to finish according to command line arguments

    :param wait:
    :param timeout:
    :param mp_id: Model Packaging name
    :param mp_client: Model Packaging Client
    """
    if not wait:
        return

    start = time.time()
    if timeout <= 0:
        raise Exception('Invalid --timeout argument: should be positive integer')

    # We create a separate client for logs because it has the different timeout settings
    log_mp_client = ModelPackagingClient.construct_from_other(mp_client)
    log_mp_client.timeout = mp_client.timeout, LOG_READ_TIMEOUT_SECONDS

    click.echo("Logs streaming...")

    while True:
        elapsed = time.time() - start
        if elapsed > timeout:
            raise Exception(TIMEOUT_ERROR_MESSAGE)

        try:
            mp = mp_client.get(mp_id)
            if mp.status.state == SUCCEEDED_STATE:
                click.echo(f'Model {mp_id} was packed. Packaging took {round(time.time() - start)} seconds')
                return
            elif mp.status.state == FAILED_STATE:
                raise Exception(f'Model packaging {mp_id} was failed.')
            elif mp.status.state == "":
                click.echo(f"Can't determine the state of {mp.id}. Sleeping...")
            else:
                for msg in log_mp_client.log(mp.id, follow=True):
                    print_logs(msg)

        except (WrongHttpStatusCode, HTTPException, RequestException, APIConnectionException) as e:
            LOGGER.info('Callback have not confirmed completion of the operation. Exception: %s', str(e))

        LOGGER.debug('Sleep before next request')
        time.sleep(DEFAULT_WAIT_TIMEOUT)
