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
import base64
import logging
from typing import List, Dict, Optional

import click
from click import ClickException

from odahuflow.cli.utils import click_utils
from odahuflow.cli.utils.client import pass_obj
from odahuflow.sdk import config
from odahuflow.sdk.clients.api import WrongHttpStatusCode
from odahuflow.sdk.clients.api_aggregated import \
    parse_resources_file, \
    parse_resources_dir, OdahuflowCloudResourceUpdatePair
from odahuflow.sdk.clients.connection import ConnectionClient
from odahuflow.sdk.clients.packaging import ModelPackagingClient
from odahuflow.sdk.clients.packaging_integration import PackagingIntegrationClient
from odahuflow.sdk.local.packaging import start_package, cleanup_packaging_docker_containers
from odahuflow.sdk.models import Connection, K8sPackager, ModelPackaging, PackagerTarget, \
    PackagingIntegration, Target, TargetSchema

LOGGER = logging.getLogger(__name__)


@click.group(name='packaging', cls=click_utils.BetterHelpGroup)
@click.option('--url', help='API server host', default=config.API_URL)
@click.option('--token', help='API server jwt token', default=config.API_TOKEN)
@click.pass_context
def packaging_group(ctx: click.core.Context, url: str, token: str):
    """
    Local packaging process.\n
    Alias for the command is pack.
    """
    ctx.obj = ModelPackagingClient(url, token)


@packaging_group.command('cleanup-containers')
def cleanup_containers():
    """
    \b
    Delete all packaging docker containers.
    \b
    Usage example:
        * odahuflowctl local pack cleanup-containers
    \f
    """
    cleanup_packaging_docker_containers()


def fetch_local_entities(manifest_file: List[str], manifest_dir: List[str]) -> List[OdahuflowCloudResourceUpdatePair]:
    """
    Collect entities from manifest files in local FS and return a result
    Manifests can be collected from file or files inside a directory
    Manifests from different sources are combined together
    """

    entities: List[OdahuflowCloudResourceUpdatePair] = []
    for file_path in manifest_file:
        entities.extend(parse_resources_file(file_path).changes)

    for dir_path in manifest_dir:
        entities.extend(parse_resources_dir(dir_path))

    return entities


def parse_entities(
        entities: List[OdahuflowCloudResourceUpdatePair], pack_id: str
) -> (ModelPackaging, Dict[str, PackagingIntegration], Dict[str, Connection]):
    mp: Optional[ModelPackaging] = None

    packagers: Dict[str, PackagingIntegration] = {}
    connections: Dict[str, Connection] = {}

    for entity in map(lambda x: x.resource, entities):
        if isinstance(entity, PackagingIntegration):
            packagers[entity.id] = entity
        elif isinstance(entity, Connection):
            connections[entity.id] = entity
        elif isinstance(entity, ModelPackaging) and entity.id == pack_id:
            mp = entity

    return mp, packagers, connections


def _decode_connection(connection: Connection) -> None:
    encoding = 'utf-8'
    decode_fields = ['password', 'key_secret', 'key_id', 'public_key']

    for f in decode_fields:
        v: str = getattr(connection.spec, f)
        if not v:
            continue

        try:
            decoded = base64.b64decode(v, validate=True)
        except Exception as e:
            LOGGER.error(f'Unable to decode base64 .spec.{f} of connection {connection.id}')
            raise e

        try:
            decoded_string = decoded.decode(encoding)
        except Exception as e:
            LOGGER.error(f'Unable to decode utf-8 .spec.{f} of connection {connection.id}')
            raise e

        setattr(connection.spec, f, decoded_string)


def get_packager(
        name: str, local: Dict[str, PackagingIntegration], remote_api: PackagingIntegrationClient
) -> PackagingIntegration:
    """
    Fetch Packager entity by name looking in local manifests at first
    and trying to fetch from web api after
    :param name: name of packager
    :param local: entities parsed from local manifests
    :param remote_api: client to fetch packagers from API server
    :return:
    """
    packager = local.get(name)
    if not packager:
        click.echo(
            f'The "{name}" packager is not found in the manifest files. '
            f'Trying to retrieve it from API server'
        )
        packager = remote_api.get(name)
    return packager


def get_packager_targets(
        targets: List[Target], connections: Dict[str, Connection], remote_api: ConnectionClient
) -> List[PackagerTarget]:
    """
    Build targets for calling packager. Fetch and base64 decode connections by names using local manifest and
    ODAHU connections API
    :param targets: Targets from packaging manifest
    :param connections: Connections found in local manifest files
    :param remote_api: ConnectionClient to fetch missing Connections
    """

    packager_targets: List[PackagerTarget] = []

    for t in targets:
        conn = connections.get(t.connection_name)
        if not conn:
            click.echo(
                f'The "{t.connection_name}" connection of "{t.name}" target is not found in the manifest files. '
                f'Trying to retrieve it from API server'
            )
            conn = remote_api.get_decrypted(t.connection_name)

        _decode_connection(conn)

        packager_targets.append(
            PackagerTarget(conn, t.name)
        )

    return packager_targets


def get_packager_target(
        target: Target, connections: Dict[str, Connection], remote_api: ConnectionClient
) -> PackagerTarget:
    """
    Build targets for calling packager. Fetch and base64 decode connections by names using local manifest and
    ODAHU connections API
    :param target: Targets from packaging manifest
    :param connections: Connections found in local manifest files
    :param remote_api: ConnectionClient to fetch missing Connections
    """

    conn = connections.get(target.connection_name)
    if not conn:
        click.echo(
            f'The "{target.connection_name}" connection of "{target.name}" target is not found in the manifest files. '
            f'Trying to retrieve it from API server'
        )
        conn = remote_api.get_decrypted(target.connection_name)

    _decode_connection(conn)

    return PackagerTarget(conn, target.name)


def _deprecation_warning(is_target_disabled: bool):
    if is_target_disabled is None:
        click.echo('[FeatureWarning] Current behavior is to disable all packager targets for a local run. '
                   'In future releases this behavior will be removed '
                   'and all targets from ModelPackaging manifest will be enabled by default. '
                   'Consider using --disable-target option to disable specific targets by name')
    else:
        click.echo('[FeatureWarning] --disable-package-targets/--no-disable-package-targets options '
                   'are deprecated and will be removed in future releases. '
                   'By default all targets will be enabled. '
                   'Consider using --disable-target option to disable specific targets by name')


def validate_targets(
        targets: List[Target], pi: PackagingIntegration, disabled_targets: List[str],
        connections: Dict[str, Connection], remote_api: ConnectionClient
) -> List[PackagerTarget]:

    result: List[PackagerTarget] = []
    errors = []

    targets_in_spec = {t.name for t in targets}
    schema = {t.name: t for t in pi.spec.schema.targets}

    # Set defaults
    for ts in pi.spec.schema.targets:
        if ts.name in targets_in_spec:  # already in spec
            continue
        if ts.name in disabled_targets:  # disabled
            continue
        targets_in_spec.add(ts.name)
        targets.append(Target(ts.default, ts.name, ))
        LOGGER.info(f'{ts.name} target default value is set to ({ts.default})')

    # validate targets
    for t in targets:
        t_schema: TargetSchema = schema.get(t.name)
        if t_schema is None:
            errors.append(f'cannot find "{t.name}" target in "{pi.id}" packaging integration')
            continue

        try:
            pt = get_packager_target(t, connections, remote_api)
        except WrongHttpStatusCode as e:
            if e.status_code == 404:
                errors.append(f'"{t.connection_name}" connection of "{t.name}" target is not found')
                continue
            raise e

        if pt.connection.spec.type not in t_schema.connection_types:
            errors.append(f'"{pt.name}" target has invalid connection type, "{pt.connection.spec.type}", '
                          f'for "{pi.id}" packaging integration the expected are {t_schema.connection_types}')
            continue

        result.append(pt)

    if errors:
        raise ClickException('Next errors were found during .spec.targets validation:\n'
                             + '\n'.join(errors))

    return result


@packaging_group.command()
@click.option('--pack-id', '--id', help='Model packaging ID', required=True)
@click.option('--manifest-file', '-f', type=click.Path(), multiple=True,
              help='Path to an ODAHU-flow manifest file')
@click.option('--manifest-dir', '-d', type=click.Path(), multiple=True,
              help='Path to a directory with ODAHU-flow manifest files')
@click.option('--artifact-path', type=click.Path(),
              help='Path to a training artifact')
@click.option('--artifact-name', '-a', type=str, help='Override artifact name from file')
# TODO Breaking Changes: remove --disable-package-targets/--no-disable-package-targets' options
@click.option('--disable-package-targets/--no-disable-package-targets', 'is_target_disabled',
              default=None, help='Disable all targets in packaging')
@click.option('--disable-target', multiple=True,
              help='Disable target in packaging')
@pass_obj
def run(client: ModelPackagingClient, pack_id: str, manifest_file: List[str], manifest_dir: List[str],
        artifact_path: str, artifact_name: str, is_target_disabled: bool, disable_target: List[str]):
    """
    \b
    Start a packaging process locally.
    \b
    Usage example:
        * odahuflowctl local pack run --id wine
    \f
    """

    _deprecation_warning(is_target_disabled)

    if is_target_disabled is None:  # Backward compatibility
        is_target_disabled = True

    entities: List[OdahuflowCloudResourceUpdatePair] = fetch_local_entities(manifest_file, manifest_dir)

    mp, packagers, connections = parse_entities(entities, pack_id)

    if not mp:
        click.echo(
            f'The "{pack_id}" packaging is not found in the manifest files. '
            f'Trying to retrieve it from API server'
        )
        mp = client.get(pack_id)

    packager = get_packager(
        mp.spec.integration_name, packagers, PackagingIntegrationClient.construct_from_other(client)
    )

    if artifact_name:
        mp.spec.artifact_name = artifact_name
        LOGGER.debug('Override the artifact name')

    if disable_target:
        LOGGER.debug(f'Next targets are disabled: {", ".join(disable_target)}')

    mp_spec_targets = mp.spec.targets if mp.spec.targets is not None else []

    # check for deprecated option
    if is_target_disabled:
        targets = []
    else:
        # Disable targets
        mp_spec_targets = [t for t in mp_spec_targets if t.name not in disable_target]

        # Validate targets
        targets = validate_targets(mp_spec_targets, packager, disable_target,
                                   connections, ConnectionClient.construct_from_other(client))

    k8s_packager = K8sPackager(
        model_packaging=mp,
        packaging_integration=packager,
        targets=targets,
    )

    result = start_package(k8s_packager, artifact_path)

    click.echo('Packager results:')
    for key, value in result.items():
        click.echo(f'* {key} - {value}')
