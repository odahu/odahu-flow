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
from collections.abc import AsyncIterable
import logging
from typing import List, Iterator

from odahuflow.sdk.clients.api import RemoteAPIClient, AsyncRemoteAPIClient
from odahuflow.sdk.definitions import MODEL_TRAINING_URL
from odahuflow.sdk.models import ModelTraining

LOGGER = logging.getLogger(__name__)
TRAINING_SUCCESS_STATE = "succeeded"
TRAINING_FAILED_STATE = "failed"


class ModelTrainingClient(RemoteAPIClient):
    """
    Model training client
    """

    def get(self, name: str) -> ModelTraining:
        """
        Get Model Training from API server

        :param name: Model Training name
        :type name: str
        :return: Model Training
        """
        return ModelTraining.from_dict(self.query(f'{MODEL_TRAINING_URL}/{name}'))

    def get_all(self) -> List[ModelTraining]:
        """
        Get all Model Trainings from API server

        :return: all Model Trainings
        """
        return [ModelTraining.from_dict(mt) for mt in self.query(MODEL_TRAINING_URL)]

    def create(self, mt: ModelTraining) -> ModelTraining:
        """
        Create Model Training

        :param mt: Model Training
        :return Message from API server
        """
        return ModelTraining.from_dict(
            self.query(MODEL_TRAINING_URL, action='POST', payload=mt.to_dict())
        )

    def edit(self, mt: ModelTraining) -> ModelTraining:
        """
        Edit Model Training

        :param mt: Model Training
        :return Message from API server
        """
        return ModelTraining.from_dict(
            self.query(MODEL_TRAINING_URL, action='PUT', payload=mt.to_dict())
        )

    def delete(self, name: str) -> str:
        """
        Delete Model Trainings

        :param name: Name of a Model Training
        :return Message from API server
        """
        return self.query(f'{MODEL_TRAINING_URL}/{name}', action='DELETE')['message']

    def log(self, name: str, follow: bool = False) -> Iterator[str]:
        """
        Stream logs from training

        :param follow: follow stream
        :param name: Name of a Model Training
        :return Message from API server
        """
        return self.stream(f'{MODEL_TRAINING_URL}/{name}/log', 'GET', params={'follow': follow})


class AsyncModelTrainingClient(AsyncRemoteAPIClient):
    """
    Model training async client
    """

    async def get(self, name: str) -> ModelTraining:
        """
        Get Model Training from API server

        :param name: Model Training name
        :type name: str
        :return: Model Training
        """
        return ModelTraining.from_dict(await self.query(f'{MODEL_TRAINING_URL}/{name}'))

    async def get_all(self) -> List[ModelTraining]:
        """
        Get all Model Trainings from API server

        :return: all Model Trainings
        """
        return [ModelTraining.from_dict(mt) for mt in await self.query(MODEL_TRAINING_URL)]

    async def create(self, mt: ModelTraining) -> ModelTraining:
        """
        Create Model Training

        :param mt: Model Training
        :return Message from API server
        """
        return ModelTraining.from_dict(
            await self.query(MODEL_TRAINING_URL, action='POST', payload=mt.to_dict())
        )

    async def edit(self, mt: ModelTraining) -> ModelTraining:
        """
        Edit Model Training

        :param mt: Model Training
        :return Message from API server
        """
        return ModelTraining.from_dict(
            await self.query(MODEL_TRAINING_URL, action='PUT', payload=mt.to_dict())
        )

    async def delete(self, name: str) -> str:
        """
        Delete Model Trainings

        :param name: Name of a Model Training
        :return Message from API server
        """
        return (await self.query(f'{MODEL_TRAINING_URL}/{name}', action='DELETE'))['message']

    async def log(self, name: str, follow: bool = False) -> AsyncIterable:
        """
        Stream logs from training

        :param follow: follow stream
        :param name: Name of a Model Training
        :return Message from API server
        """
        params = {'follow': 'true' if follow else 'false'}
        async for chunk in self.stream(f'{MODEL_TRAINING_URL}/{name}/log', 'GET', params=params):
            yield chunk
