#
#    Copyright 2020 EPAM Systems
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
Local packaging commands for odahuflow cli
"""
import logging
from typing import List, Dict, Optional

import click
from odahuflow.cli.utils.client import pass_obj
from odahuflow.sdk import config
from odahuflow.sdk.clients.api_aggregated import \
    parse_resources_file, \
    parse_resources_dir, OdahuflowCloudResourceUpdatePair
from odahuflow.sdk.clients.packaging import ModelPackagingClient
from odahuflow.sdk.clients.packaging_integration import PackagingIntegrationClient
from odahuflow.sdk.local.packaging import start_package, cleanup_packaging_docker_containers
from odahuflow.sdk.models import K8sPackager, ModelPackaging, PackagingIntegration

LOGGER = logging.getLogger(__name__)


@click.group()
@click.option('--url', help='API server host', default=config.API_URL)
@click.option('--token', help='API server jwt token', default=config.API_TOKEN)
@click.pass_context
def packaging(ctx: click.core.Context, url: str, token: str):
    """
    Local packaging process.\n
    Alias for the command is pack.
    """
    ctx.obj = ModelPackagingClient(url, token)


@packaging.command('cleanup-containers')
@pass_obj
def cleanup_containers():
    """
    Delete all packaging docker containers.
    Usage example:\n
        * odahuflowctl local pack cleanup\n
    \f
    """
    cleanup_packaging_docker_containers()


@packaging.command()
@click.option('--pack-id', '--id', help='Model packaging ID', required=True)
@click.option('--manifest-file', '-f', type=click.Path(), multiple=True,
              help='Path to a ODAHU-flow manifest file')
@click.option('--manifest-dir', '-d', type=click.Path(), multiple=True,
              help='Path to a directory with ODAHU-flow manifest files')
@click.option('--artifact-path', type=click.Path(),
              help='Path to a training artifact')
@click.option('--artifact-name', '-a', type=str, help='Override artifact name from file')
@click.option('--disable-package-targets/--no-disable-package-targets', 'is_target_disabled',
              default=True, help='Disable all targets in packaging')
@pass_obj
def run(client: ModelPackagingClient, pack_id: str, manifest_file: List[str], manifest_dir: List[str],
        artifact_path: str, artifact_name: str, is_target_disabled: bool):
    """
    Start a packaging process locally.\n
    Usage example:\n
        * odahuflowctl local pack run --id wine\n
    \f
    """
    entities: List[OdahuflowCloudResourceUpdatePair] = []
    for file_path in manifest_file:
        entities.extend(parse_resources_file(file_path).changes)

    for dir_path in manifest_dir:
        entities.extend(parse_resources_dir(dir_path))

    mp: Optional[ModelPackaging] = None

    packagers: Dict[str, PackagingIntegration] = {}
    for entity in map(lambda x: x.resource, entities):
        if isinstance(entity, PackagingIntegration):
            packagers[entity.id] = entity
        elif isinstance(entity, ModelPackaging) and entity.id == pack_id:
            mp = entity

    if not mp:
        LOGGER.debug(
            f'The {pack_id} packaging not found in the manifest files.'
            f' Trying to retrieve it from API server'
        )
        mp = client.get(pack_id)

    integration_name = mp.spec.integration_name
    packager = packagers.get(integration_name)
    if not packager:
        LOGGER.debug(
            f'The {integration_name} packager not found in the manifest files.'
            f' Trying to retrieve it from API server'
        )
        packager = PackagingIntegrationClient.construct_from_other(client).get(integration_name)

    if artifact_name:
        mp.spec.artifact_name = artifact_name

        LOGGER.debug('Override the artifact namesdsdsdsdsdsdsd')

    if is_target_disabled:
        mp.spec.targets = []

    k8s_packager = K8sPackager(
        model_packaging=mp,
        packaging_integration=packager,
        targets=[],
    )

    result = start_package(k8s_packager, artifact_path)

    click.echo('Packager results:')
    for key, value in result.items():
        click.echo(f'* {key} - {value}')
