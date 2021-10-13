#
#    Copyright 2017 EPAM Systems
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
API client
"""
import typing

from odahuflow.sdk.clients.api import RemoteAPIClient, AsyncRemoteAPIClient
from odahuflow.sdk.definitions import TRAINING_INTEGRATION_URL
from odahuflow.sdk.models import TrainingIntegration


class TrainingIntegrationClient(RemoteAPIClient):
    """
    HTTP training integration client
    """

    def get(self, name: str) -> TrainingIntegration:
        """
        Get Training Integration from API server

        :param name: Training Integration name
        :type name: str
        :return: Training Integration
        """
        return TrainingIntegration.from_dict(self.query(f'{TRAINING_INTEGRATION_URL}/{name}'))

    def get_all(self) -> typing.List[TrainingIntegration]:
        """
        Get all Training Integrations from API server

        :return: all Training Integrations
        """
        return [TrainingIntegration.from_dict(ti) for ti in self.query(TRAINING_INTEGRATION_URL)]

    def create(self, ti: TrainingIntegration) -> TrainingIntegration:
        """
        Create Training Integration

        :param ti: Training Integration
        :return Message from API server
        """
        return TrainingIntegration.from_dict(
            self.query(TRAINING_INTEGRATION_URL, action='POST', payload=ti.to_dict())
        )

    def edit(self, ti: TrainingIntegration) -> TrainingIntegration:
        """
        Edit Training Integration

        :param ti: Training Integration
        :return Message from API server
        """
        return TrainingIntegration.from_dict(
            self.query(TRAINING_INTEGRATION_URL, action='PUT', payload=ti.to_dict())
        )

    def delete(self, name: str) -> str:
        """
        Delete Training Integrations

        :param name: Name of a Training Integration
        :return Message from API server
        """
        return self.query(f'{TRAINING_INTEGRATION_URL}/{name}', action='DELETE')


class AsyncTrainingIntegrationClient(AsyncRemoteAPIClient):
    """
    HTTP Training integration async client
    """
    async def get(self, name: str) -> TrainingIntegration:
        """
        Get Training Integration from API server

        :param name: Training Integration name
        :type name: str
        :return: Training Integration
        """
        return TrainingIntegration.from_dict(await self.query(f'{TRAINING_INTEGRATION_URL}/{name}'))

    async def get_all(self) -> typing.List[TrainingIntegration]:
        """
        Get all Training Integrations from API server

        :return: all Training Integrations
        """
        return [TrainingIntegration.from_dict(ti) for ti in await self.query(TRAINING_INTEGRATION_URL)]

    async def create(self, ti: TrainingIntegration) -> TrainingIntegration:
        """
        Create Training Integration

        :param ti: Training Integration
        :return Message from API server
        """
        return TrainingIntegration.from_dict(
            await self.query(TRAINING_INTEGRATION_URL, action='POST', payload=ti.to_dict())
        )

    async def edit(self, ti: TrainingIntegration) -> TrainingIntegration:
        """
        Edit Training Integration

        :param ti: Training Integration
        :return Message from API server
        """
        return TrainingIntegration.from_dict(
            await self.query(TRAINING_INTEGRATION_URL, action='PUT', payload=ti.to_dict())
        )

    async def delete(self, name: str) -> str:
        """
        Delete Training Integrations

        :param name: Name of a Training Integration
        :return Message from API server
        """
        return await self.query(f'{TRAINING_INTEGRATION_URL}/{name}', action='DELETE')
