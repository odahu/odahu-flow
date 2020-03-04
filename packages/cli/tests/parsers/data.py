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
Test data for entity CLI commands
"""
import typing

import click
from odahuflow.cli.parsers import connection, training, toolchain_integration, packaging, packaging_integration, \
    deployment, route
from odahuflow.sdk.clients.api import RemoteAPIClient
from odahuflow.sdk.clients.connection import ConnectionClient
from odahuflow.sdk.clients.deployment import ModelDeploymentClient, READY_STATE
from odahuflow.sdk.clients.packaging import ModelPackagingClient, SUCCEEDED_STATE
from odahuflow.sdk.clients.packaging_integration import PackagingIntegrationClient
from odahuflow.sdk.clients.route import ModelRouteClient
from odahuflow.sdk.clients.toolchain_integration import ToolchainIntegrationClient
from odahuflow.sdk.clients.training import ModelTrainingClient
from odahuflow.sdk.models import ModelTraining, Connection, ConnectionSpec, ModelTrainingSpec, ModelIdentity, \
    ToolchainIntegration, ToolchainIntegrationSpec, ModelPackaging, ModelPackagingSpec, PackagingIntegration, \
    PackagingIntegrationSpec, ModelDeployment, ModelDeploymentSpec, ModelRoute, ModelRouteSpec, \
    ModelTrainingStatus, ModelPackagingStatus, ModelDeploymentStatus, ModelRouteStatus
from odahuflow.sdk.models.base_model_ import Model

ENTITY_ID = 'entity-id'


# This interface for any test cases that tests CLI commands for entities
class EntityTestData(typing.NamedTuple):
    entity_client: RemoteAPIClient
    entity: Model
    click_group: click.Group
    # For example, ModelTraining or Connection
    kind: str


def generate_entities_for_test(metafunc, keys: typing.List[str]):
    if "entity_test_data" in metafunc.fixturenames:
        metafunc.parametrize(
            "entity_test_data",
            keys,
            indirect=True,
        )


CONNECTION = EntityTestData(
    ConnectionClient(),
    Connection(
        id=ENTITY_ID,
        spec=ConnectionSpec(
            key_secret="mock-key-secret",
            uri="mock-url",
        ),
    ),
    connection.connection,
    'Connection',
)

TRAINING = EntityTestData(
    ModelTrainingClient(),
    ModelTraining(
        id=ENTITY_ID,
        spec=ModelTrainingSpec(
            work_dir="/1/2/3",
            model=ModelIdentity(
                name="name",
                version="version"
            )
        ),
        status=ModelTrainingStatus(
            state=SUCCEEDED_STATE
        )
    ),
    training.training,
    'ModelTraining'
)

TOOLCHAIN = EntityTestData(
    ToolchainIntegrationClient(),
    ToolchainIntegration(
        id=ENTITY_ID,
        spec=ToolchainIntegrationSpec(
            default_image="mock-image",
            entrypoint="default-entrypoint",
        ),
    ),
    toolchain_integration.toolchain_integration,
    'ToolchainIntegration',
)

PACKAGING = EntityTestData(
    ModelPackagingClient(),
    ModelPackaging(
        id=ENTITY_ID,
        spec=ModelPackagingSpec(
            artifact_name='test-artifact-name',
            integration_name='test'
        ),
        status=ModelPackagingStatus(
            state=SUCCEEDED_STATE,
        )
    ),
    packaging.packaging,
    'ModelPackaging',
)

PACKAGING_INTEGRATION = EntityTestData(
    PackagingIntegrationClient(),
    PackagingIntegration(
        id=ENTITY_ID,
        spec=PackagingIntegrationSpec(
            default_image="odahu:image",
            entrypoint="some_entrypoint"
        ),
    ),
    packaging_integration.packaging_integration,
    'PackagingIntegration',
)

DEPLOYMENT = EntityTestData(
    ModelDeploymentClient(),
    ModelDeployment(
        id=ENTITY_ID,
        spec=ModelDeploymentSpec(
            image="odahu:image",
            min_replicas=0
        ),
        status=ModelDeploymentStatus(
            state=READY_STATE,
            available_replicas=1,
        )
    ),
    deployment.deployment,
    'ModelDeployment',
)

ROUTER = EntityTestData(
    ModelRouteClient(),
    ModelRoute(
        id=ENTITY_ID,
        spec=ModelRouteSpec(
            mirror="test",
        ),
        status=ModelRouteStatus(
            state=READY_STATE,
        )
    ),
    route.route,
    'ModelRoute',
)
