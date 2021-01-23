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
Version file
"""
import click
from odahuflow.sdk.version import __version__ as __sdk_version__

__version__ = '1.4.0-rc2'


@click.command()
def version():
    """
    Show version of cli and sdk
    """
    click.echo(f'cli version: {__version__}\n'
               f'sdk version: {__sdk_version__}')
