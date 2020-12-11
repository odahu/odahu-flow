import json
import logging
from os.path import join
from pathlib import PurePosixPath
from typing import Dict, Any

import docker
from docker.models.containers import Container
from docker.types import Mount
from odahuflow.sdk import config
from odahuflow.sdk.local.docker_utils import stream_container_logs, \
    convert_labels_to_filter, cleanup_docker_containers, PACKAGING_DOCKER_LABELS, raise_error_if_container_failed
from odahuflow.sdk.logger import is_verbose_enabled
from odahuflow.sdk.models import K8sPackager

PACKAGER_CONF_FILE_PATH = 'mp.json'
PACKAGER_RESULT_FILE_PATH = 'result.json'
ARTIFACT_PATH = '/trained_artifact'

LOGGER = logging.getLogger(__name__)


def create_mp_config_file(config_dir: str, packager: K8sPackager) -> None:
    with open(join(config_dir, PACKAGER_CONF_FILE_PATH), 'w') as f:
        packager_dict = packager.to_dict()
        json.dump(packager_dict, f)

        LOGGER.debug(f"Saved the packaging configuration:\n{json.dumps(packager_dict, indent=2)}")


def read_mp_result_file(config_dir: str) -> Dict[str, Any]:
    with open(join(config_dir, PACKAGER_RESULT_FILE_PATH)) as f:
        return json.load(f)


def start_package(packager: K8sPackager, artifact_path: str) -> Dict[str, Any]:
    if not artifact_path:
        artifact_path = config.LOCAL_MODEL_OUTPUT_DIR

    # removing .zip extension if there's one
    artifact_name = packager.model_packaging.spec.artifact_name.strip()
    if artifact_name.endswith('.zip'):
        packager.model_packaging.spec.artifact_name = artifact_name[:-4]

    # make full artifact path with artifact name
    artifact_path = join(
        artifact_path,
        packager.model_packaging.spec.artifact_name
    )

    create_mp_config_file(artifact_path, packager)

    client = docker.from_env()
    container: Container = client.containers.run(
        packager.model_packaging.spec.image or packager.packaging_integration.spec.default_image,
        stderr=True,
        stdout=True,
        working_dir=ARTIFACT_PATH,
        command=[
            packager.packaging_integration.spec.entrypoint,
            ARTIFACT_PATH,
            # specified path for Docker Linux Container ignoring OS
            str(PurePosixPath(ARTIFACT_PATH, PACKAGER_CONF_FILE_PATH)),
        ] + [v for v in ("--verbose",) if is_verbose_enabled()],
        mounts=[
            Mount(ARTIFACT_PATH, artifact_path, type="bind"),
            Mount("/var/run/docker.sock", "/var/run/docker.sock", type="bind")
        ],
        detach=True,
        labels=PACKAGING_DOCKER_LABELS,
    )

    container_info = client.api.inspect_container(container.id)
    LOGGER.debug(f'Container info:\n{json.dumps(container_info, indent=2)}')

    print(f'Packaging docker image {container.id} has started. Stream logs:')
    stream_container_logs(container)

    raise_error_if_container_failed(container.id)

    return read_mp_result_file(artifact_path)


def cleanup_packaging_docker_containers():
    cleanup_docker_containers(convert_labels_to_filter(PACKAGING_DOCKER_LABELS))
