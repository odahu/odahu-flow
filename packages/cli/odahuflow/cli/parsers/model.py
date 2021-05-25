#
#    Copyright 2019 EPAM Systems
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
EDGE commands for odahuflow cli
"""
import json

import click

from odahuflow.cli.utils import click_utils
from odahuflow.sdk import config
from odahuflow.sdk.clients.model import ModelClient


@click.group(cls=click_utils.BetterHelpGroup)
def model():
    """
    Allow you to perform actions on deployed models
    """
    pass


@model.command()
@click.option('--model-route', '--mr', default=config.MODEL_ROUTE_NAME, type=str, help='Name of Model Route')
@click.option('--model-deployment', '--md', default=config.MODEL_DEPLOYMENT_NAME, type=str,
              help='Name of Model Deployment')
@click.option('--base-url', default=config.API_URL, type=str, help='Base model server url')
@click.option('--url-prefix', type=str, help='Url prefix of model server')
@click.option('--token', type=str, default=config.API_TOKEN, help='Model jwt token')
@click.option('--json', 'json_input', type=str, help='Json parameter. For example: --json {"x": 2}')
@click.option('--json-file', '--file', type=click.Path(exists=True), help='Path to json file')
def invoke(json_input, model_route: str, model_deployment: str, url_prefix: str,
           base_url: str, token: str, json_file: str):
    """
    Invoke model endpoint.
    \f
    :param json_input:
    :param client:
    :return: None
    """
    if json_file:
        with open(json_file) as f:
            json_input = f.read()

    if not (url_prefix or model_route or model_deployment):
        raise ValueError(
            'Cannot create a model url. Specify one of the options: --url-prefix/--model-route/--model-deployment'
        )

    client = ModelClient(base_url, model_route, model_deployment, url_prefix, token)

    result = client.invoke(**json.loads(json_input))

    click.echo(json.dumps(result))


@model.command()
@click.option('--model-route', '--mr', default=config.MODEL_ROUTE_NAME, type=str, help='Name of Model Route')
@click.option('--model-deployment', '--md', default=config.MODEL_DEPLOYMENT_NAME, type=str,
              help='Name of Model Deployment')
@click.option('--base-url', default=config.API_URL, type=str, help='Base model server url')
@click.option('--url-prefix', type=str, help='Url prefix of model server')
@click.option('--token', type=str, default=config.API_TOKEN, help='Model jwt token')
def info(model_route: str, model_deployment: str, url_prefix: str, base_url: str, token: str):
    """
    Get model information.
    \f
    :param client: Model HTTP Client
    """
    if not (url_prefix or model_route or model_deployment):
        raise ValueError(
            'Cannot create a model url. Specify one of the options: --url-prefix/--model-route/--model-deployment'
        )

    client = ModelClient(base_url, model_route, model_deployment, url_prefix, token)

    result = client.info()

    print(result)
