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
from typing import Dict, List

import docker
from docker.models.containers import Container

WORKSPACE_PATH = '/workspace'
TRAINING_DOCKER_LABELS: Dict[str, str] = {
    "component": "odahu",
    "api": "training"
}
PACKAGING_DOCKER_LABELS: Dict[str, str] = {
    "component": "odahu",
    "api": "packaging"
}


def raise_error_if_container_failed(container_id: str) -> None:
    client = docker.from_env()

    container_info = client.api.inspect_container(container_id)

    container_exit_code = container_info.get("State", {}).get("ExitCode", 0)
    if container_exit_code != 0:
        raise Exception(f'Container finished with {container_exit_code} error code')


def stream_container_logs(container: Container) -> None:
    """
    Stream logs of the Docker container to stdout
    :param container: Docker container
    """
    logs = container.logs(stream=True, follow=True)
    for log in logs:
        for line in log.splitlines():
            print(f'[Container {container.id[:5]}] {line.decode()}')


def convert_labels_to_filter(labels: Dict[str, str]) -> List[str]:
    """
    The docker client has the following label filter format: "key=value"
    :param labels: container labels
    :return: filters
    """
    return list(map(lambda item: f'{item[0]}={item[1]}', labels.items()))


def cleanup_docker_containers(labels: List[str]):
    client = docker.from_env()
    containers = client.containers.list(filters={
        "status": "exited",
        "label": labels,
    })

    if not containers:
        print('There are no containers')
        return

    for container in containers:
        print(f'{container.id} container has been deleted')
        container.remove()
