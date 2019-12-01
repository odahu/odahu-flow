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
Template CLI commands
"""
import click
from odahuflow.sdk.clients.templates import get_odahuflow_template_content
from odahuflow.sdk.clients.templates import get_odahuflow_template_names


@click.group()
def template():
    """
    Allow you to perform actions on odahuflow template files
    """
    pass


@template.command(name='all')
def template_all():
    """
    Get all template names.\n
    Usage example:\n
        * odahuflowctl template all\n
    """
    nl = "\n * "
    click.echo(f'Templates:{nl}{nl.join(get_odahuflow_template_names())}')


@template.command()
@click.option('--name', help='Template name', required=True)
def generate(name: str):
    """
    Generate a template by name.
    Format of templates is YAML.\n
    To find all templates execute the following command:\n
        * odahuflowctl template all\n
    Usage example:\n
        * odahuflowctl template generate --name deployment
    """
    click.echo(get_odahuflow_template_content(name))
