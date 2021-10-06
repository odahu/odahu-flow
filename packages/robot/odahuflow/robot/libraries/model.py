#
#    Copyright 2018 EPAM Systems
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
Robot test library - model API
"""
import requests

from odahuflow.sdk.containers import headers as odahuflow_headers


class Model:
    """
    Model API class
    """

    def __init__(self):
        """
        Init model
        """
        self._last_response_id = None
        self._last_response = None

    @staticmethod
    def get_model_info(edge, token, model_name, model_version=None):
        """
        Invoke model through API

        :param model_name: model name
        :param model_version: model version
        :param edge: edge url
        :param token: model API JWT token
        :return: dict -- response
        """
        headers = {"Authorization": f"Bearer {token}"}
        if model_version:
            url = f'{edge}/api/model/{model_name}/{model_version}/info'
        else:
            url = f'{edge}/api/model/{model_name}/info'

        print(f'Requesting {url} in GET mode')

        response = requests.get(
            url,
            headers=headers
        )

        if response.status_code != 200:
            raise Exception(f'Returned wrong status code: {response.status_code}')

        return response.json()

    def invoke_model_feedback(self, md_name, model_name, model_version, edge, token, request_id, **payload):
        """
        Invoke model through API

        :param model_name: model name
        :param model_version: model version
        :param edge: edge url
        :param token: model API JWT token
        :param request_id: request ID
        :param payload: payload dict
        :return: dict -- response
        """
        headers = {
            'Authorization': f'Bearer {token}',
            odahuflow_headers.MODEL_REQUEST_ID: request_id,
            odahuflow_headers.MODEL_NAME: model_name,
            odahuflow_headers.MODEL_VERSION: model_version,
        }

        url = f'{edge}/api/v1/feedback/model/{md_name}/api/model'

        print(f'Requesting {url} with data = {payload} in POST mode')

        response = requests.post(
            url,
            json=payload,
            headers=headers
        )

        if not response.ok:
            raise Exception(f'Returned wrong status code: {response.status_code}, body: {response.text}')

        return response.json()

    def invoke_model_api(self, md_name, edge, token, request_id=None, **payload):
        """
        Invoke model through API

        :param md_name: model deployment name
        :param edge: edge url
        :param token: model API JWT token
        :param request_id: (Optional) request ID
        :param payload: payload dict
        :return: dict -- response
        """
        headers = {'Authorization': f'Bearer {token}'}
        if request_id:
            headers[odahuflow_headers.MODEL_REQUEST_ID] = request_id

        url = f'{edge}/model/{md_name}/api/model/invoke'

        print(f'Requesting {url} with data = {payload} in POST mode')

        payload = {
            'data': [list(payload.values())],
            'columns': list(payload.keys()),
        }

        response = requests.post(
            url,
            json=payload,
            headers=headers
        )

        if not response.ok:
            raise Exception(f'Returned wrong status code: {response.status_code}, body: {response.text},'
                            f' payload: {payload}')

        self._last_response_id = response.headers.get(odahuflow_headers.MODEL_REQUEST_ID)
        self._last_response = response.json()
        return self._last_response

    def get_model_api_last_response(self):
        """
        Get model last response

        :return: dict -- last response
        """
        return self._last_response

    def get_model_api_last_response_id(self):
        """
        Get last model response ID

        :return: str -- last response ID
        """
        return self._last_response_id
