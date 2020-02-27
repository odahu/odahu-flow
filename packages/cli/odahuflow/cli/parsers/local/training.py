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
Local training commands for odahuflow cli
"""
import json
import logging
from typing import List, Dict, Optional

import click
from odahuflow.cli.utils.client import pass_obj
from odahuflow.cli.utils.output import PLAIN_TEXT_OUTPUT_FORMAT, JSON_OUTPUT_FORMAT
from odahuflow.sdk import config
from odahuflow.sdk.clients.api_aggregated import \
    parse_resources_file, \
    parse_resources_dir, OdahuflowCloudResourceUpdatePair
from odahuflow.sdk.clients.toolchain_integration import ToolchainIntegrationClient
from odahuflow.sdk.clients.training import ModelTraining, ModelTrainingClient
from odahuflow.sdk.local.training import start_train, list_local_trainings, cleanup_local_artifacts, \
    cleanup_training_docker_containers
from odahuflow.sdk.models import ToolchainIntegration, K8sTrainer

LOGGER = logging.getLogger(__name__)


@click.group()
@click.option('--url', help='API server host', default=config.API_URL)
@click.option('--token', help='API server jwt token', default=config.API_TOKEN)
@click.pass_context
def training(ctx: click.core.Context, url: str, token: str):
    """
    Local training process.\n
    Alias for the command is train.
    """
    ctx.obj = ModelTrainingClient(url, token)


@training.command("list")
@click.option('--output-format', '-o', 'output_format',
              default=PLAIN_TEXT_OUTPUT_FORMAT, type=click.Choice([PLAIN_TEXT_OUTPUT_FORMAT, JSON_OUTPUT_FORMAT]))
def training_list(output_format: str):
    """
    Get list of local training artifacts.\n
    Get all training artifacts:\n
        odahuflowctl local train list\n
    Get all training artifacts in json format:\n
        odahuflowctl local train -o json\n
    \f
    :param output_format: Output format
    """
    artifacts = list_local_trainings()

    if output_format == JSON_OUTPUT_FORMAT:
        click.echo(json.dumps(artifacts, indent=2))
    else:
        if not artifacts:
            click.echo('Artifacts not found')
            return

        click.echo('Training artifacts:')

        for artifact in artifacts:
            click.echo(f'* {artifact}')


@training.command('cleanup-artifacts')
def cleanup_artifacts():
    """
    Delete all training local artifacts.
    Usage example:\n
        * odahuflowctl local train cleanup\n
    \f
    """
    cleanup_local_artifacts()


@training.command('cleanup-containers')
def cleanup_containers():
    """
    Delete all training docker containers.
    Usage example:\n
        * odahuflowctl local train cleanup\n
    \f
    """
    cleanup_training_docker_containers()


@training.command()
@click.option('--train-id', '--id', help='Model training ID', required=True)
@click.option('--manifest-file', '-f', type=click.Path(), multiple=True,
              help='Path to a ODAHU-flow manifest file')
@click.option('--manifest-dir', '-d', type=click.Path(), multiple=True,
              help='Path to a directory with ODAHU-flow manifest files')
@click.option('--output-dir', '--output', type=click.Path(), help='Directory where model artifact will be saved')
@pass_obj
def run(client: ModelTrainingClient, train_id: str, manifest_file: List[str], manifest_dir: List[str],
        output_dir: str):
    """
    Start a training process locally.\n
    Usage example:\n
        * odahuflowctl local train run --id examples-git\n
    \f
    """
    entities: List[OdahuflowCloudResourceUpdatePair] = []
    for file_path in manifest_file:
        entities.extend(parse_resources_file(file_path).changes)

    for dir_path in manifest_dir:
        entities.extend(parse_resources_dir(dir_path))

    mt: Optional[ModelTraining] = None

    # find a training
    toolchains: Dict[str, ToolchainIntegration] = {}
    for entity in map(lambda x: x.resource, entities):
        if isinstance(entity, ToolchainIntegration):
            toolchains[entity.id] = entity
        elif isinstance(entity, ModelTraining) and entity.id == train_id:
            mt = entity

    if not mt:
        LOGGER.debug(f'{train_id} training not found. Trying to retrieve it from API server')
        mt = client.get(train_id)

    toolchain = toolchains.get(mt.spec.toolchain)
    if not toolchain:
        LOGGER.debug(f'{toolchain} toolchain not found. Trying to retrieve it from API server')
        toolchain = ToolchainIntegrationClient.construct_from_other(client).get(mt.spec.toolchain)

    trainer = K8sTrainer(
        model_training=mt,
        toolchain_integration=toolchain,
    )

    start_train(trainer, output_dir)
