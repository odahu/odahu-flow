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
from __future__ import print_function

import glob
import os
import time
from argparse import Namespace

import legion.config
import legion.containers.docker
import legion.edi.deploy as deploy
import legion.io
import legion.model.client
import legion.model.model_id
from legion.model import ModelClient
from legion.utils import TemporaryFolder

import docker
import docker.errors
import docker.models.containers
import docker.models.images
import numpy
import pandas
import requests
import unittest2


class TestDeploy(unittest2.TestCase):
    MODEL_ID = 'temp'
    MODEL_VERSION = '1.8'

    def setUp(self):
        common_arguments = Namespace(docker_network=None)
        self.client = legion.containers.docker.build_docker_client(common_arguments)
        self.network = legion.containers.docker.find_network(self.client, common_arguments)
        self.wheel_path = self._get_latest_bdist()
        legion.model.model_id.init(self.MODEL_ID)
        self.deployed_containers_ids = []

    def tearDown(self):
        legion.model.model_id._model_id = None
        legion.model.model_id._model_initialized_from_function = False
        for container_id in self.deployed_containers_ids:
            try:
                container = self.client.containers.get(container_id)
                container.stop()
                container.remove()
            except docker.errors.NotFound:
                pass

    def _get_latest_bdist(self):
        dist_dir = os.path.abspath(os.path.join(os.path.dirname(__file__), '..', 'dist'))
        if not os.path.exists(dist_dir):
            raise Exception('Cannot find dist dir: %s' % dist_dir)

        list_of_files = glob.glob('%s/*.whl' % (dist_dir,))
        latest_file = max(list_of_files, key=os.path.getctime)
        return latest_file

    def test_stack_is_running(self):
        containers = legion.containers.docker.get_stack_containers_and_images(self.client, self.network)
        self.assertTrue(len(containers['services']) > 0, 'Cannot found any service container')
        for container in containers['services']:
            container_required = container.labels.get('com.epam.legion.container_required', 'true').lower() \
                                 in ('1', 'yes', 'true')
            if container_required:
                container_name = container.labels.get('com.epam.legion.container_description', container.image.tags[0])
                self.assertEqual(container.status, 'running', 'Wrong status of required container %s' % container_name)

    def _build_model(self, path, version):
        def prepare(x):
            x['additional'] = x['d_int']
            return x

        def apply(x):
            assert type(x) == pandas.DataFrame

            assert x['d_int'].dtype == numpy.int
            assert x['d_float'].dtype == numpy.float

        df = pandas.DataFrame([{
            'd_int': 1,
            'd_float': 1.0,
        }])

        return legion.io.export(path,
                                apply,
                                prepare,
                                input_data_frame=df,
                                version=version)

    def _build_summation_model(self, path, version):
        def prepare(x):
            return x

        def apply(x):
            return {'x': x['a'] + x['b']}

        df = pandas.DataFrame([{
            'a': 1,
            'b': 1,
        }])

        return legion.io.export(path,
                                apply,
                                prepare,
                                input_data_frame=df,
                                use_df=False,
                                version=version)

    def test_model_image_build(self, remove_image=True, summation_model=False):
        self.test_stack_is_running()

        with TemporaryFolder('legion-model-image-build') as temp_folder:
            path = os.path.join(temp_folder.path, 'temp.model')

            if summation_model:
                self._build_summation_model(path, self.MODEL_VERSION)
            else:
                self._build_model(path, self.MODEL_VERSION)

            args = Namespace(
                model_file=path,
                model_id=None,
                base_docker_image=None,
                docker_network=None,
                python_package=self.wheel_path,
                python_package_version=None,
                python_repository=None,
                docker_image_tag=None,
                serving=deploy.VALID_SERVING_WORKERS[1],
                push_to_registry=None
            )

            image = deploy.build_model(args)
            self.assertIsInstance(image, docker.models.images.Image)

        if remove_image:
            try:
                if image:
                    self.client.images.remove(image.short_id)
            except Exception:
                pass
        else:
            return image

    def _build_image_deploy_and_test(self, image, deploy_args, undeploy=False):
        container = deploy.deploy_model(deploy_args)

        self.deployed_containers_ids.append(container.id)

        self.assertIsInstance(container, docker.models.containers.Container)
        time.sleep(3)
        container = self.client.containers.get(container.id)
        logs = container.logs().decode('utf-8')

        print('--- CONTAINER LOGS ---')
        print(logs)

        if undeploy:
            container.stop()
            container.remove()

            self.client.images.remove(image.short_id)

        self.assertEqual(container.status, 'running', 'Wrong status after deploy')

        return None if undeploy else container

    def test_model_image_deploy_by_model_id(self):
        image = self.test_model_image_build(False)
        args = Namespace(
            model_id=self.MODEL_ID,
            docker_image=None,
            docker_network=None,
            grafana_server=None,
            grafana_user=None,
            grafana_password=None,
            expose_model_port=None
        )

        self._build_image_deploy_and_test(image, args, True)

    def test_model_image_deploy_and_query_exposed(self):
        image = self.test_model_image_build(False, True)
        model_port = 9009

        args = Namespace(
            model_id=self.MODEL_ID,
            docker_image=None,
            docker_network=None,
            grafana_server=None,
            grafana_user=None,
            grafana_password=None,
            expose_model_port=model_port
        )

        container = self._build_image_deploy_and_test(image, args, False)
        ports_information = [item for sublist in container.attrs['NetworkSettings']['Ports'].values()
                             for item in sublist]
        ports_information = [int(x['HostPort']) for x in ports_information]
        self.assertTrue(model_port in ports_information, 'Port not binded')

        values = {
            'a': 10,
            'b': 20
        }

        client = ModelClient(self.MODEL_ID, 'http://{}:{}'.format('localhost', model_port))

        model_information = client.info()
        self.assertEqual(model_information['version'], self.MODEL_VERSION, 'Incorrect model version')
        self.assertEqual(model_information['use_df'], False, 'Incorrect model use_df field')
        self.assertDictEqual(model_information['input_params'],
                             {'b': {'numpy_type': 'int64', 'type': 'Integer'},
                              'a': {'numpy_type': 'int64', 'type': 'Integer'}},
                             'Incorrect model input_params')

        invoke_result = client.invoke(**values)
        self.assertEqual(invoke_result['x'], 30, 'Incorrect model result')

        args = Namespace(
            model_id=self.MODEL_ID,
            docker_network=None,
            grafana_server=None,
            grafana_user=None,
            grafana_password=None,
        )

        deploy.undeploy_model(args)
        time.sleep(3)

        self.client.images.remove(image.short_id)

    def test_model_image_deploy_by_image_id(self):
        image = self.test_model_image_build(False)
        args = Namespace(
            model_id=None,
            docker_image=image.id,
            docker_network=None,
            grafana_server=None,
            grafana_user=None,
            grafana_password=None,
            expose_model_port=None
        )

        self._build_image_deploy_and_test(image, args, True)

    def test_model_deploy_undeploy(self):
        image = self.test_model_image_build(False)

        args = Namespace(
            model_id=self.MODEL_ID,
            docker_image=None,
            docker_network=None,
            grafana_server=None,
            grafana_user=None,
            grafana_password=None,
            expose_model_port=None,
            push_to_registry=None
        )

        container = self._build_image_deploy_and_test(image, args, False)

        args = Namespace(
            model_id=self.MODEL_ID,
            docker_network=None,
            grafana_server=None,
            grafana_user=None,
            grafana_password=None,
        )

        deploy.undeploy_model(args)
        time.sleep(3)

        self.client.images.remove(image.short_id)

        self.assertRaises(docker.errors.NotFound, self.client.containers.get, container.id)
        self.assertRaises(docker.errors.NotFound, self.client.images.get, image.short_id)


if __name__ == '__main__':
    unittest2.main()
