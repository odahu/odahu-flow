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
from odahuflow.sdk.definitions import FEEDBACK_URL
from odahuflow.sdk.models import FeedbackModelFeedback


class FeedbackClient(RemoteAPIClient):
    """
    HTTP Feedback client
    """

    def post(self, feedback: dict, model_name: str, model_version: str, request_id: str) -> FeedbackModelFeedback:
        """
        Get Feedback from API server

        :param feedback: feedback dict
        :param model_name: model name
        :param model_version: model version
        :param request_id: request id
        :return: Feedback Model
        """
        headers = {
            'model-name': model_name,
            'model-version': model_version,
            'request-id': request_id
        }

        return FeedbackModelFeedback.from_dict(
            self.query(FEEDBACK_URL, action='POST', payload=feedback, headers=headers)
        )


class AsyncFeedbackClient(AsyncRemoteAPIClient):
    """
    HTTP Feedback async client
    """

    async def post(self, feedback: dict, model_name: str, model_version: str, request_id: str) -> FeedbackModelFeedback:
        """
        Get FeedbackModelFeedbackResponse from API server

        :param feedback: feedback dict
        :param model_name: model name
        :param model_version: model version
        :param request_id: request id
        :return: Feedback Model
        """
        headers = {
            'model-name': model_name,
            'model-version': model_version,
            'request-id': request_id
        }

        return FeedbackModelFeedback.from_dict(
            await self.query(FEEDBACK_URL, action='POST', payload=feedback, headers=headers)
        )
