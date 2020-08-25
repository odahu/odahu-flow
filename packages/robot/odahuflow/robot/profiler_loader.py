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

import json
import os
import typing

from odahuflow.sdk.clients.oauth_handler import do_client_cred_authentication
from odahuflow.sdk.clients.oidc import OpenIdProviderConfiguration

API_URL_PARAM_NAME = 'API_URL'
AUTH_TOKEN_PARAM_NAME = 'AUTH_TOKEN'

CLUSTER_PROFILE = 'CLUSTER_PROFILE'

# Key in cluster_profile.json that contain creds for
TEST_SA_ADMIN = 'test'
TEST_SA_DS = 'test_data_scientist'


class ServiceAccount:

    def __init__(self, profile_json, service_account_key):
        self._client_id = profile_json['service_accounts'][service_account_key]['client_id']
        self._client_secret = profile_json['service_accounts'][service_account_key]['client_secret']

    @property
    def client_id(self):
        return self._client_id

    @property
    def client_secret(self):
        return self._client_secret


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

            test_sa_admin = ServiceAccount(data, TEST_SA_ADMIN)
            test_sa_ds = ServiceAccount(data, TEST_SA_DS)

            host_base_domain = data['dns']['domain']
            variables = {
                'HOST_BASE_DOMAIN': host_base_domain,
                'CLUSTER_NAME': data.get('cluster_name'),
                'CLUSTER_CONTEXT': data.get('cluster_context'),
                'FEEDBACK_BUCKET': data.get('data_bucket'),
                'EXAMPLES_VERSION': data.get('examples_version'),
                'CLOUD_TYPE': data['cloud']['type'],
                'EDGE_URL': os.getenv('EDGE_URL', f'https://{host_base_domain}'),
                API_URL_PARAM_NAME: os.getenv(API_URL_PARAM_NAME, f'https://{host_base_domain}'),
                'GRAFANA_URL': os.getenv('GRAFANA_URL', f'https://{host_base_domain}/grafana'),
                'PROMETHEUS_URL': os.getenv('PROMETHEUS_URL', f'https://{host_base_domain}/prometheus'),
                'ALERTMANAGER_URL': os.getenv('ALERTMANAGER_URL', f'https://{host_base_domain}/alertmanager'),
                'JUPYTERLAB_URL': os.getenv('JUPITERLAB_URL', f'https://{host_base_domain}/jupyterlab'),
                'MLFLOW_URL': os.getenv('MLFLOW_URL', f'https://{host_base_domain}/mlflow'),
                'AIRFLOW_URL': os.getenv('AIRFLOW_URL', f'https://{host_base_domain}/airflow'),
                'KIBANA_URL': os.getenv('KIBANA_URL', f'https://{host_base_domain}/kibana'),
                'IS_GPU_ENABLED': 'training_gpu' in data['node_pools'],
                'SA_CLIENT_ID': test_sa_admin.client_id,
                'SA_CLIENT_SECRET': test_sa_admin.client_secret,
                'SA_ADMIN': test_sa_admin,
                'SA_DATA_SCIENTIST': test_sa_ds,
                'ISSUER': data.get('oauth_oidc_issuer_url')
            }
        except Exception as err:
            raise Exception("Can\'t get variable from cluster profile: {}".format(err))

        try:
            client_id = test_sa_admin.client_id
            client_secret = test_sa_admin.client_secret
            issuer = data['oauth_oidc_issuer_url']
            conf = OpenIdProviderConfiguration(issuer)
            conf.fetch_configuration()

            login_result = do_client_cred_authentication(
                issue_token_url=conf.token_endpoint, client_id=client_id, client_secret=client_secret
            )

            if login_result:
                variables[AUTH_TOKEN_PARAM_NAME] = login_result.id_token
            else:
                variables[AUTH_TOKEN_PARAM_NAME] = ''
        except Exception as err:
            raise Exception("Can\'t get dex authentication data: {}".format(err))

    return variables
