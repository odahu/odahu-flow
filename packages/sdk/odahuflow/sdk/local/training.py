import json
import logging
import os
import re
import uuid
from collections import namedtuple
from os import listdir
from pathlib import PurePosixPath
from shutil import rmtree
from typing import List

import docker
from docker.models.containers import Container
from docker.types import Mount

from odahuflow.sdk import config
from odahuflow.sdk.gppi.executor import PROJECT_FILE
from odahuflow.sdk.local.docker_utils import stream_container_logs, TRAINING_DOCKER_LABELS, WORKSPACE_PATH, \
    convert_labels_to_filter, cleanup_docker_containers, raise_error_if_container_failed
from odahuflow.sdk.models import K8sTrainer

MODEL_OUTPUT_CONTAINER_PATH = '/output'
TRAINER_CONF_PATH = 'mt.json'

LOGGER = logging.getLogger(__name__)


# Context object for compiling result model directory name from Go template
TemplateContext = namedtuple('TemplateContext', ['Name', 'Version', 'RandomUUID'])

DEFAULT_MODEL_DIR_TEMPLATE = '{{ .Name }}-{{ .Version }}-{{ .RandomUUID }}'


def compile_artifact_name_template(go_template: str, context: TemplateContext) -> str:
    """
    Converts artifact name template from Go to Python format and compiles it
    :param go_template: result model directory name template, e.g. {{ .Name }}-{{ .Version}}-{{ .RandomUUID }}
    :param context: contains values to put into template
    :return: compiled directory name
    """

    # Go -> Python template conversion
    py_template = re.sub(r'{{\s*?\.([a-zA-Z_0-9]+)\s*?}}',
                         lambda x: f'{{ctx.{x.group(1)}}}',
                         go_template)
    return py_template.format(ctx=context)


def launch_training_container(trainer: K8sTrainer, output_dir: str) -> None:
    client = docker.from_env()
    container: Container = client.containers.run(
        trainer.model_training.spec.image or trainer.toolchain_integration.spec.default_image,
        stderr=True,
        stdout=True,
        working_dir=WORKSPACE_PATH,
        command=[
            trainer.toolchain_integration.spec.entrypoint,
            "--verbose",
            "--mt",
            # specified path for Docker Linux Container ignoring OS
            str(PurePosixPath(WORKSPACE_PATH, TRAINER_CONF_PATH)),
            "--target",
            MODEL_OUTPUT_CONTAINER_PATH,
        ],
        mounts=[
            Mount(WORKSPACE_PATH, os.getcwd(), type="bind"),
            Mount(MODEL_OUTPUT_CONTAINER_PATH, output_dir, type="bind")
        ],
        detach=True,
        labels=TRAINING_DOCKER_LABELS,
    )

    container_info = client.api.inspect_container(container.id)
    LOGGER.debug(f'Container info:\n{json.dumps(container_info, indent=2)}')

    print(f'Training docker image {container.id} has started. Stream logs:')
    stream_container_logs(container)

    raise_error_if_container_failed(container.id)


def launch_gppi_validation_container(trainer: K8sTrainer, output_dir: str) -> None:
    client = docker.from_env()
    container: Container = client.containers.run(
        trainer.model_training.spec.image or trainer.toolchain_integration.spec.default_image,
        stderr=True,
        stdout=True,
        working_dir=WORKSPACE_PATH,
        command=[
            "odahuflowctl",
            "--verbose",
            "gppi",
            "--gppi-model-path",
            MODEL_OUTPUT_CONTAINER_PATH,
            "--env-name",
            "odahu_model",
            "test"
        ],
        mounts=[
            Mount(WORKSPACE_PATH, os.getcwd(), type="bind"),
            Mount(MODEL_OUTPUT_CONTAINER_PATH, output_dir, type="bind")
        ],
        detach=True,
        labels=TRAINING_DOCKER_LABELS,
    )

    container_info = client.api.inspect_container(container.id)
    LOGGER.debug(f'Container info:\n{json.dumps(container_info, indent=2)}')

    print(f'Validation docker image {container.id} has started. Stream logs:')
    stream_container_logs(container)

    raise_error_if_container_failed(container.id)


def create_mt_config_file(trainer: K8sTrainer) -> None:
    with open(TRAINER_CONF_PATH, 'w') as f:
        trainer_dict = trainer.to_dict()
        json.dump(trainer_dict, f)

        LOGGER.debug(f"Saved the trainer configuration:\n{json.dumps(trainer_dict, indent=2)}")


def start_train(trainer: K8sTrainer, output_dir: str) -> None:
    """
    :param trainer: container object for all configuration objects of a training
    :param output_dir: path to directory to save result artifact (relative or absolute)
    """
    create_mt_config_file(trainer)

    if not output_dir:
        output_dir = config.LOCAL_MODEL_OUTPUT_DIR
        LOGGER.debug(f'Output directory for model training is not provided. Using default: {output_dir}')

    model_dir_name_template = trainer.model_training.spec.model.artifact_name_template or DEFAULT_MODEL_DIR_TEMPLATE

    # Removing .zip extension if there's one
    model_dir_name_template = model_dir_name_template.strip()
    if model_dir_name_template.endswith('.zip'):
        model_dir_name_template = model_dir_name_template[:-4]

    model_dir_name = compile_artifact_name_template(
        go_template=model_dir_name_template,
        context=TemplateContext(Name=trainer.model_training.spec.model.name,
                                Version=trainer.model_training.spec.model.version,
                                RandomUUID=uuid.uuid4())
    )

    model_dir_path = os.path.abspath(os.path.join(output_dir, model_dir_name))
    os.makedirs(model_dir_path, exist_ok=True)

    launch_training_container(trainer, model_dir_path)
    launch_gppi_validation_container(trainer, model_dir_path)

    print(f'Model {model_dir_name} was saved in the {output_dir} directory')


def list_local_trainings() -> List[str]:
    if not os.path.exists(config.LOCAL_MODEL_OUTPUT_DIR):
        return []

    def is_training_artifact(name: str) -> bool:
        full_path = os.path.join(config.LOCAL_MODEL_OUTPUT_DIR, name)
        return os.path.isdir(full_path) and PROJECT_FILE in os.listdir(full_path)

    return sorted(filter(is_training_artifact, listdir(config.LOCAL_MODEL_OUTPUT_DIR)))


def cleanup_local_artifacts():
    if os.path.exists(config.LOCAL_MODEL_OUTPUT_DIR):
        rmtree(config.LOCAL_MODEL_OUTPUT_DIR)

        LOGGER.debug(f'Remove the {config.LOCAL_MODEL_OUTPUT_DIR} directory')


def cleanup_training_docker_containers():
    cleanup_docker_containers(convert_labels_to_filter(TRAINING_DOCKER_LABELS))


if __name__ == '__main__':
    cleanup_training_docker_containers()
