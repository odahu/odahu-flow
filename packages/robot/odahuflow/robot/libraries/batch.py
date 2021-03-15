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
import json
import os

from odahuflow.robot.cloud import object_storage
from odahuflow.sdk.clients.api_aggregated import parse_resources_file_with_one_item
from odahuflow.sdk.clients.batch_job import InferenceJob


class BatchUtils:

    def __init__(self, cloud_type, bucket, cluster_name):
        """
        Init client
        """
        self._client = object_storage.build_client(cloud_type, bucket, cluster_name)

    def check_batch_job_response(self, manifest_path: str, expected_output_path: str) -> bool:
        """
        Checks equity of output file in .spec.output_destination.path.<base_name> and expected file at
        `expected_output_path` location. Where <base-name> is base name of `expected_output_path`
        :param manifest_path:
        :param expected_output_path:
        :return:
        """
        with open(expected_output_path, "r") as f:
            exp = json.load(f)
        with open(manifest_path, "r") as f:
            job: InferenceJob = parse_resources_file_with_one_item(manifest_path).resource

            base_name = os.path.basename(expected_output_path)
            output_file_path = os.path.join(job.spec.output_destination.path, base_name)
            act = json.loads(self._client.read_file(output_file_path))
        return exp == act
