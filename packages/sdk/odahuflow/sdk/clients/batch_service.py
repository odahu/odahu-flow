#
#      Copyright 2021 EPAM Systems
#
#      Licensed under the Apache License, Version 2.0 (the "License");
#      you may not use this file except in compliance with the License.
#      You may obtain a copy of the License at
#
#          http://www.apache.org/licenses/LICENSE-2.0
#
#      Unless required by applicable law or agreed to in writing, software
#      distributed under the License is distributed on an "AS IS" BASIS,
#      WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
#      See the License for the specific language governing permissions and
#      limitations under the License.

import logging
import typing

from odahuflow.sdk.clients.api import RemoteAPIClient, AsyncRemoteAPIClient
from odahuflow.sdk.definitions import INFERENCE_SERVICE_URL
from odahuflow.sdk.models import InferenceService

LOGGER = logging.getLogger(__name__)


class BatchInferenceServiceClient(RemoteAPIClient):
    """
    HTTP InferenceService client
    """

    def get(self, service_id: str) -> InferenceService:
        """
        Get InferenceService from API server

        :param service_id: InferenceService ID
        :return: InferenceService
        """
        return InferenceService.from_dict(self.query(f'{INFERENCE_SERVICE_URL}/{service_id}'))

    def get_all(self) -> typing.List[InferenceService]:
        """
        Get all Services from API server

        :return: all Services
        """
        return [InferenceService.from_dict(service) for service in self.query(INFERENCE_SERVICE_URL)]

    def create(self, service: InferenceService) -> InferenceService:
        """
        Create InferenceService

        :param service: InferenceService
        :return Message from API server
        """
        return InferenceService.from_dict(self.query(INFERENCE_SERVICE_URL, action='POST', payload=service.to_dict()))

    def edit(self, service: InferenceService) -> InferenceService:
        """
        Edit InferenceService

        :param service: InferenceService
        :return Message from API server
        """
        return InferenceService.from_dict(self.query(INFERENCE_SERVICE_URL, action='PUT', payload=service.to_dict()))

    def delete(self, name: str) -> str:
        """
        Delete Services

        :param name: Name of a InferenceService
        :return Message from API server
        """
        return self.query(f'{INFERENCE_SERVICE_URL}/{name}', action='DELETE')['message']


class AsyncBatchInferenceServiceClient(AsyncRemoteAPIClient):
    """
    HTTP InferenceService async client
    """

    async def get(self, service_id: str) -> InferenceService:
        """
        Get InferenceService from API server

        :param service_id: InferenceService ID
        :return: InferenceService
        """
        return InferenceService.from_dict(await self.query(f'{INFERENCE_SERVICE_URL}/{service_id}'))

    async def get_all(self) -> typing.List[InferenceService]:
        """
        Get all Services from API server

        :return: all Services
        """
        return [InferenceService.from_dict(service) for service in await self.query(INFERENCE_SERVICE_URL)]

    async def create(self, service: InferenceService) -> InferenceService:
        """
        Create InferenceService

        :param service: InferenceService
        :return Message from API server
        """
        return InferenceService.from_dict(await self.query(INFERENCE_SERVICE_URL,
                                                           action='POST', payload=service.to_dict()))

    async def edit(self, service: InferenceService) -> InferenceService:
        """
        Edit InferenceService

        :param service: InferenceService
        :return Message from API server
        """
        return InferenceService.from_dict(await self.query(INFERENCE_SERVICE_URL, action='PUT',
                                                           payload=service.to_dict()))

    async def delete(self, name: str) -> str:
        """
        Delete Services

        :param name: Name of a InferenceService
        :return Message from API server
        """
        return (await self.query(f'{INFERENCE_SERVICE_URL}/{name}', action='DELETE'))['message']
