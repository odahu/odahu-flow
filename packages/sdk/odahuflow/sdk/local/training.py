import datetime
import json
import logging
import os
from os import listdir
from os.path import join
from pathlib import Path
from shutil import rmtree
from typing import List

import docker
from docker.models.containers import Container
from docker.types import Mount
from odahuflow.sdk import config
from odahuflow.sdk.local.docker_utils import stream_container_logs, TRAINING_DOCKER_LABELS, WORKSPACE_PATH, \
    convert_labels_to_filter, cleanup_docker_containers, raise_error_if_container_failed
from odahuflow.sdk.models import K8sTrainer

MODEL_OUTPUT_CONTAINER_PATH = '/output'
TRAINER_CONF_PATH = 'mt.json'

LOGGER = logging.getLogger(__name__)


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
            join(WORKSPACE_PATH, TRAINER_CONF_PATH),
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


def start_train(trainer: K8sTrainer, output_dir: str):
    create_mt_config_file(trainer)

    if not output_dir:
        # For example, 01-Mar-2020-17-38-04
        suffix = datetime.datetime.now().strftime("%d-%b-%Y-%H-%M-%S")
        output_dir = join(
            config.LOCAL_MODEL_OUTPUT_DIR,
            f'{trainer.model_training.spec.model.name}-{trainer.model_training.spec.model.version}-{suffix}'
        )

        LOGGER.debug(f'Output model training directory is not provided. Generate the directory: {output_dir}')

    os.makedirs(output_dir, exist_ok=True)

    launch_training_container(trainer, output_dir)

    launch_gppi_validation_container(trainer, output_dir)

    print(f'Model {Path(output_dir).name} was saved in the {output_dir} directory')


def list_local_trainings() -> List[str]:
    if not os.path.exists(config.LOCAL_MODEL_OUTPUT_DIR):
        return []

    return listdir(config.LOCAL_MODEL_OUTPUT_DIR)


def cleanup_local_artifacts():
    if os.path.exists(config.LOCAL_MODEL_OUTPUT_DIR):
        rmtree(config.LOCAL_MODEL_OUTPUT_DIR)

        LOGGER.debug(f'Remove the {config.LOCAL_MODEL_OUTPUT_DIR} directory')


def cleanup_training_docker_containers():
    cleanup_docker_containers(convert_labels_to_filter(TRAINING_DOCKER_LABELS))


if __name__ == '__main__':
    cleanup_training_docker_containers()
