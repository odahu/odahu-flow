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
from os import makedirs
from pathlib import Path

import requests


class ExamplesLoader:
    def __init__(self, examples_http_url: str, examples_version: str) -> None:
        self._examples_http_url = examples_http_url
        self._examples_version = examples_version

    def download_file(self, remote_file_path: str, result_file_path: str):
        """
        Download file from the ODAHU example repository and save it to a file

        :param remote_file_path: File path in the example repository
        :param result_file_path: Save a content of remote file by this path
        """
        # TODO: replace by urljoin
        file_http_url = "/".join(
            [
                self._examples_http_url,
                self._examples_version,
                remote_file_path,
            ]
        )

        # TODO: replace with stream downloading
        resp = requests.get(file_http_url)
        resp.raise_for_status()

        # Create all intermediate-level directories
        makedirs(Path(result_file_path).parent, exist_ok=True)

        with open(result_file_path, "w") as f:
            f.write(resp.text)
