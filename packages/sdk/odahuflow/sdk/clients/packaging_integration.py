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
from odahuflow.sdk.definitions import PACKING_INTEGRATION_URL
from odahuflow.sdk.models import PackagingIntegration


class PackagingIntegrationClient(RemoteAPIClient):
    """
    HTTP packaging integration client
    """

    def get(self, name: str) -> PackagingIntegration:
        """
        Get Packaging Integration from API server

        :param name: Packaging Integration name
        :type name: str
        :return: Packaging Integration
        """
        return PackagingIntegration.from_dict(self.query(f'{PACKING_INTEGRATION_URL}/{name}'))

    def get_all(self) -> typing.List[PackagingIntegration]:
        """
        Get all Packaging Integrations from API server

        :return: all Packaging Integrations
        """
        return [PackagingIntegration.from_dict(mr) for mr in self.query(PACKING_INTEGRATION_URL)]

    def create(self, mr: PackagingIntegration) -> PackagingIntegration:
        """
        Create Packaging Integration

        :param mr: Packaging Integration
        :return Message from API server
        """
        return PackagingIntegration.from_dict(
            self.query(PACKING_INTEGRATION_URL, action='POST', payload=mr.to_dict())
        )

    def edit(self, mr: PackagingIntegration) -> PackagingIntegration:
        """
        Edit Packaging Integration

        :param mr: Packaging Integration
        :return Message from API server
        """
        return PackagingIntegration.from_dict(
            self.query(PACKING_INTEGRATION_URL, action='PUT', payload=mr.to_dict())
        )

    def delete(self, name: str) -> str:
        """
        Delete Packaging Integrations

        :param name: Name of a Packaging Integration
        :return Message from API server
        """
        return self.query(f'{PACKING_INTEGRATION_URL}/{name}', action='DELETE')


class AsyncPackagingIntegrationClient(AsyncRemoteAPIClient):
    """
    HTTP packaging integration async client
    """

    async def get(self, name: str) -> PackagingIntegration:
        """
        Get Packaging Integration from API server

        :param name: Packaging Integration name
        :type name: str
        :return: Packaging Integration
        """
        return PackagingIntegration.from_dict(await self.query(f'{PACKING_INTEGRATION_URL}/{name}'))

    async def get_all(self) -> typing.List[PackagingIntegration]:
        """
        Get all Packaging Integrations from API server

        :return: all Packaging Integrations
        """
        return [PackagingIntegration.from_dict(mr) for mr in await self.query(PACKING_INTEGRATION_URL)]

    async def create(self, mr: PackagingIntegration) -> PackagingIntegration:
        """
        Create Packaging Integration

        :param mr: Packaging Integration
        :return Message from API server
        """
        return PackagingIntegration.from_dict(
            await self.query(PACKING_INTEGRATION_URL, action='POST', payload=mr.to_dict())
        )

    async def edit(self, mr: PackagingIntegration) -> PackagingIntegration:
        """
        Edit Packaging Integration

        :param mr: Packaging Integration
        :return Message from API server
        """
        return PackagingIntegration.from_dict(
            await self.query(PACKING_INTEGRATION_URL, action='PUT', payload=mr.to_dict())
        )

    async def delete(self, name: str) -> str:
        """
        Delete Packaging Integrations

        :param name: Name of a Packaging Integration
        :return Message from API server
        """
        return await self.query(f'{PACKING_INTEGRATION_URL}/{name}', action='DELETE')
