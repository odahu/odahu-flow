#
#    Copyright 2021 EPAM Systems
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

from odahuflow.sdk.clients.api import RemoteAPIClient, AsyncRemoteAPIClient
from odahuflow.sdk.definitions import USER_INFO_URL
from odahuflow.sdk.models import UserInfo


class UserInfoClient(RemoteAPIClient):
    """
    HTTP user info client
    """

    def get(self) -> UserInfo:
        """
        Get User Info from API server

        :return: UserInfo
        """
        return UserInfo.from_dict(self.query(USER_INFO_URL))


class AsyncUserInfoClient(AsyncRemoteAPIClient):
    """
    HTTP user info client
    """

    async def get(self) -> UserInfo:
        """
        Get User Info from API server

        :return: User Info
        """
        return UserInfo.from_dict(await self.query(USER_INFO_URL))
