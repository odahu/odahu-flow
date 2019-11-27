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
import logging
import typing

from odahuflow.sdk.clients.api import RemoteAPIClient, AsyncRemoteAPIClient
from odahuflow.sdk.definitions import CONNECTION_URL
from odahuflow.sdk.models import Connection

LOGGER = logging.getLogger(__name__)


class ConnectionClient(RemoteAPIClient):
    """
    HTTP connection client
    """

    def get(self, conn_id: str) -> Connection:
        """
        Get Connection from API server

        :param conn_id: Connection ID
        :return: Connection
        """
        return Connection.from_dict(self.query(f'{CONNECTION_URL}/{conn_id}'))

    # TODO: Remove after implementation of the issue https://github.com/legion-platform/legion/issues/1008
    def get_decrypted(self, conn_id: str, decrypt_token: str) -> Connection:
        """
        Get decrypted connection from API server

        :param decrypt_token: Token for getting a decrypted connection
        :param conn_id: Connection ID
        :return: Connection
        """
        return Connection.from_dict(self.query(
            f'{CONNECTION_URL}/{conn_id}/decrypted',
            payload={'token': decrypt_token}
        ))

    def get_all(self) -> typing.List[Connection]:
        """
        Get all Connections from API server

        :return: all Connections
        """
        return [Connection.from_dict(conn) for conn in self.query(CONNECTION_URL)]

    def create(self, conn: Connection) -> Connection:
        """
        Create Connection

        :param conn: Connection
        :return Message from API server
        """
        return Connection.from_dict(self.query(CONNECTION_URL, action='POST', payload=conn.to_dict()))

    def edit(self, conn: Connection) -> Connection:
        """
        Edit Connection

        :param conn: Connection
        :return Message from API server
        """
        return Connection.from_dict(self.query(CONNECTION_URL, action='PUT', payload=conn.to_dict()))

    def delete(self, name: str) -> str:
        """
        Delete Connections

        :param name: Name of a Connection
        :return Message from API server
        """
        return self.query(f'{CONNECTION_URL}/{name}', action='DELETE')['message']


class AsyncConnectionClient(AsyncRemoteAPIClient):
    """
    HTTP connection async client
    """

    async def get(self, conn_id: str) -> Connection:
        """
        Get Connection from API server

        :param conn_id: Connection ID
        :return: Connection
        """
        return Connection.from_dict(await self.query(f'{CONNECTION_URL}/{conn_id}'))

    # TODO: Remove after implementation of the issue https://github.com/legion-platform/legion/issues/1008
    async def get_decrypted(self, conn_id: str, decrypt_token: str) -> Connection:
        """
        Get decrypted connection from API server

        :param decrypt_token: Token for getting a decrypted connection
        :param conn_id: Connection ID
        :return: Connection
        """
        return Connection.from_dict(await self.query(
            f'{CONNECTION_URL}/{conn_id}/decrypted',
            payload={'token': decrypt_token}
        ))

    async def get_all(self) -> typing.List[Connection]:
        """
        Get all Connections from API server

        :return: all Connections
        """
        return [Connection.from_dict(conn) for conn in await self.query(CONNECTION_URL)]

    async def create(self, conn: Connection) -> Connection:
        """
        Create Connection

        :param conn: Connection
        :return Message from API server
        """
        return Connection.from_dict(await self.query(CONNECTION_URL, action='POST', payload=conn.to_dict()))

    async def edit(self, conn: Connection) -> Connection:
        """
        Edit Connection

        :param conn: Connection
        :return Message from API server
        """
        return Connection.from_dict(await self.query(CONNECTION_URL, action='PUT', payload=conn.to_dict()))

    async def delete(self, name: str) -> str:
        """
        Delete Connections

        :param name: Name of a Connection
        :return Message from API server
        """
        return (await self.query(f'{CONNECTION_URL}/{name}', action='DELETE'))['message']
