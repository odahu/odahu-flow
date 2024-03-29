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
import sys
import typing

import click

from odahuflow.cli.utils import click_utils
from odahuflow.cli.utils.click_utils import auth_options
from odahuflow.cli.utils.client import pass_obj
from odahuflow.sdk.clients.api import RemoteAPIClient
from odahuflow.sdk.clients.api_aggregated import apply as api_aggregated_apply
from odahuflow.sdk.clients.api_aggregated import parse_resources_file, OdahuflowCloudResourceUpdatePair

LOGGER = logging.getLogger(__name__)


@click.group(cls=click_utils.BetterHelpGroup)
@auth_options
@click.pass_context
def bulk(ctx: click.core.Context, api_client: RemoteAPIClient):
    """
    Bulk operations on Odahuflow resources
    """
    ctx.obj = api_client


@bulk.command()
@click.argument('file', required=True, type=click.Path())
@pass_obj
def apply(client: RemoteAPIClient, file: str):
    """
    Create/Update Odahuflow resources on an API.\n
    You should specify a path to file with resources.
    For now, CLI supports yaml and JSON file formats.
    Usage example:\n
        * odahuflowctl bulk apply resources.odahuflow.yaml
    \f
    :param client: Generic API HTTP client
    :param file: Path to the file with odahuflow resources
    """
    process_bulk_operation(client, file, False)


@bulk.command()
@click.argument('file', required=True, type=click.Path())
@pass_obj
def delete(client: RemoteAPIClient, file: str):
    """
    Remove Odahuflow resources from an API.\n
    You should specify a path to file with resources.
    For now, CLI supports yaml and JSON file formats.
    Usage example:\n
        * odahuflowctl bulk delete resources.odahuflow.yaml
    \f
    :param client: Generic API HTTP client
    :param file: Path to the file with odahuflow resources
    """
    process_bulk_operation(client, file, True)


def _print_resources_info_counter(objects: typing.Tuple[OdahuflowCloudResourceUpdatePair]) -> str:
    """
    Output count of resources with their's types

    :param objects: resources to print
    :return: str
    """
    names = ', '.join(f'{type(obj.resource).__name__} {obj.resource_id}' for obj in objects)
    if objects:
        return f'{len(objects)} ({names})'
    else:
        return str(len(objects))


def process_bulk_operation(api_client: RemoteAPIClient, filename: str, is_removal: bool):
    """
    Apply bulk operation helper

    :param api_client: base API client to extract connection options from
    :param filename: path to file with resources
    :param is_removal: is it removal operation
    """
    odahuflow_resources = parse_resources_file(filename)
    result = api_aggregated_apply(odahuflow_resources, api_client, is_removal)
    output = ['Operation completed']
    if result.created:
        output.append(f'created resources: {_print_resources_info_counter(result.created)}')
    if result.changed:
        output.append(f'changed resources: {_print_resources_info_counter(result.changed)}')
    if result.removed:
        output.append(f'removed resources: {_print_resources_info_counter(result.removed)}')
    click.echo(', '.join(output))

    if result.errors:
        click.echo('Some errors detected:')
        for error in result.errors:
            click.echo(f'\t{error}')
        sys.exit(1)
