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
Variables loader from json Cluster Profile
"""

import os
import json

import typing
from odahuflow.robot.libraries import auth_client

API_URL_PARAM_NAME = 'API_URL'
AUTH_TOKEN_PARAM_NAME = 'AUTH_TOKEN'

CLUSTER_PROFILE = 'CLUSTER_PROFILE'


def get_variables(profile=None) -> typing.Dict[str, str]:
    """
    Gather and return all variables to robot

    :param profile: path to cluster profile
    :type profile: str
    :return: dict[str, Any] -- values for robot
    """

    # load Cluster Profile
    profile = profile or os.getenv(CLUSTER_PROFILE)

    if not profile:
        raise Exception('Can\'t get profile at path {}'.format(profile))
    if not os.path.exists(profile):
        raise Exception('Can\'t get profile - {} file not found'.format(profile))

    with open(profile, 'r') as json_file:
        data = json.load(json_file)
        variables = {}

        try:
            host_base_domain = "odahu.{}.{}".format(data.get('cluster_name'), data.get('root_domain'))
            variables = {
                'HOST_BASE_DOMAIN': host_base_domain,
                'CLUSTER_NAME': data.get('cluster_name'),
                'CLUSTER_CONTEXT': data.get('cluster_context'),
                'FEEDBACK_BUCKET': data.get('data_bucket'),
                'CLOUD_TYPE': data.get('cloud_type'),
                'EDGE_URL': os.getenv('EDGE_URL', f'https://{host_base_domain}'),
                API_URL_PARAM_NAME: os.getenv(API_URL_PARAM_NAME, f'https://{host_base_domain}'),
                'GRAFANA_URL': os.getenv('GRAFANA_URL', f'https://{host_base_domain}/grafana'),
                'PROMETHEUS_URL': os.getenv('PROMETHEUS_URL', f'https://{host_base_domain}/prometheus'),
                'ALERTMANAGER_URL': os.getenv('ALERTMANAGER_URL', f'https://{host_base_domain}/alertmanager'),
                'JUPYTERLAB_URL': os.getenv('JUPITERLAB_URL', f'https://{host_base_domain}/jupyterlab'),
                'MLFLOW_URL': os.getenv('MLFLOW_URL', f'https://{host_base_domain}/mlflow'),
                # TODO: Remove after implementation of the issue https://github.com/legion-platform/legion/issues/1008
                'CONN_DECRYPT_TOKEN': data.get('odahuflow_connection_decrypt_token'),
            }
        except Exception as err:
            raise Exception("Can\'t get variable from cluster profile: {}".format(err))

        try:
            if data.get('test_user_email') and data.get('test_user_password'):
                variables[AUTH_TOKEN_PARAM_NAME] = auth_client.init_token(
                    data['test_user_email'],
                    data['test_user_password'],
                    data['test_oauth_auth_url'],
                    data['test_oauth_client_id'],
                    data['test_oauth_client_secret'],
                    data['test_oauth_scope']
                )
            else:
                variables[AUTH_TOKEN_PARAM_NAME] = ''
        except Exception as err:
            raise Exception("Can\'t get dex authentication data: {}".format(err))

    return variables
