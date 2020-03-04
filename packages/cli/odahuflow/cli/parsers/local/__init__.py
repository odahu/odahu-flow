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
"""
Group of local commands
"""
import click
from odahuflow.cli.parsers.local.packaging import packaging
from odahuflow.cli.parsers.local.training import training
from odahuflow.cli.utils.abbr import AbbreviationGroup


@click.group(cls=AbbreviationGroup)
def local():
    """
    Train and package locally
    """
    pass


LOCAL_GROUPS = [
    training,
    packaging,
]

# Initialize local groups
for group in LOCAL_GROUPS:
    local.add_command(group)
