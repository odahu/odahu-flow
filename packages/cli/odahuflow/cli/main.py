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
CLI entrypoint
"""

import click
import click_completion
from odahuflow.cli.parsers.bulk import bulk
from odahuflow.cli.parsers.completion import completion
from odahuflow.cli.parsers.config import config_group
from odahuflow.cli.parsers.connection import connection
from odahuflow.cli.parsers.deployment import deployment
from odahuflow.cli.parsers.gppi import gppi
from odahuflow.cli.parsers.local import local
from odahuflow.cli.parsers.model import model
from odahuflow.cli.parsers.packaging import packaging
from odahuflow.cli.parsers.packaging_integration import packaging_integration
from odahuflow.cli.parsers.route import route
from odahuflow.cli.parsers.security import login, logout
from odahuflow.cli.parsers.template import template
from odahuflow.cli.parsers.toolchain_integration import toolchain_integration
from odahuflow.cli.parsers.training import training
from odahuflow.cli.utils.abbr import AbbreviationGroup
from odahuflow.cli.utils.error_handler import cli_error_handler
from odahuflow.sdk.logger import configure_logging
from odahuflow.cli.version import version

COMMAND_GROUPS = [
    config_group,
    connection,
    deployment,
    model,
    packaging,
    packaging_integration,
    bulk,
    route,
    template,
    gppi,
    toolchain_integration,
    training,
    login,
    logout,
    version,
    completion,
    local,
]

# Initialize shell completion
click_completion.init()


@click.group(cls=AbbreviationGroup)
@click.option('--verbose/--non-verbose', default=False)
def base(verbose=False):
    """
    Odahuflow CLI
    """
    configure_logging(verbose)


# Initialize all commands
for group in COMMAND_GROUPS:
    base.add_command(group)


def main():
    """
    Main CLI entrypoint
    """
    with cli_error_handler():
        base()


if __name__ == '__main__':
    base()
