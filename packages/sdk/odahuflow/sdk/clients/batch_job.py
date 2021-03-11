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
from odahuflow.sdk.definitions import INFERENCE_JOB_URL
from odahuflow.sdk.models import InferenceJob

LOGGER = logging.getLogger(__name__)


SUCCESS_STATE = "succeeded"
FAILED_STATE = "failed"


class BatchInferenceJobClient(RemoteAPIClient):
    """
    HTTP InferenceJob client
    """

    def get(self, job_id: str) -> InferenceJob:
        """
        Get InferenceJob from API server

        :param job_id: InferenceJob ID
        :return: InferenceJob
        """
        return InferenceJob.from_dict(self.query(f'{INFERENCE_JOB_URL}/{job_id}'))

    def get_all(self) -> typing.List[InferenceJob]:
        """
        Get all Jobs from API server

        :return: all Jobs
        """
        return [InferenceJob.from_dict(job) for job in self.query(INFERENCE_JOB_URL)]

    def create(self, job: InferenceJob) -> InferenceJob:
        """
        Create InferenceJob

        :param job: InferenceJob
        :return Message from API server
        """
        return InferenceJob.from_dict(self.query(INFERENCE_JOB_URL, action='POST', payload=job.to_dict()))

    def delete(self, name: str) -> str:
        """
        Delete Jobs

        :param name: Name of a InferenceJob
        :return Message from API server
        """
        return self.query(f'{INFERENCE_JOB_URL}/{name}', action='DELETE')['message']

    def edit(self, job: InferenceJob) -> InferenceJob:
        raise NotImplementedError


class AsyncBatchInferenceJobClient(AsyncRemoteAPIClient):
    """
    HTTP InferenceJob async client
    """

    async def get(self, job_id: str) -> InferenceJob:
        """
        Get InferenceJob from API server

        :param job_id: InferenceJob ID
        :return: InferenceJob
        """
        return InferenceJob.from_dict(await self.query(f'{INFERENCE_JOB_URL}/{job_id}'))

    async def get_all(self) -> typing.List[InferenceJob]:
        """
        Get all Jobs from API server

        :return: all Jobs
        """
        return [InferenceJob.from_dict(job) for job in await self.query(INFERENCE_JOB_URL)]

    async def create(self, job: InferenceJob) -> InferenceJob:
        """
        Create InferenceJob

        :param job: InferenceJob
        :return Message from API server
        """
        return InferenceJob.from_dict(await self.query(INFERENCE_JOB_URL, action='POST', payload=job.to_dict()))

    async def delete(self, name: str) -> str:
        """
        Delete Jobs

        :param name: Name of a InferenceJob
        :return Message from API server
        """
        return (await self.query(f'{INFERENCE_JOB_URL}/{name}', action='DELETE'))['message']

    async def edit(self, job: InferenceJob) -> InferenceJob:
        raise NotImplementedError
