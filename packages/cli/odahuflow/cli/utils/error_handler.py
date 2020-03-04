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
import logging
import sys
from contextlib import contextmanager

import click
from odahuflow.sdk.logger import is_verbose_enabled

LOG = logging.getLogger(__name__)

CLI_ERROR_MESSAGE = 'Error: {}.\nFor more information rerun command with' \
                    ' --verbose flag'


@contextmanager
def cli_error_handler():
    """
    Print error message to stdout and exit if any error occurs
    """
    try:
        yield
    except Exception as e:
        # If the verbose flag is enabled than the stacktrace will be logged.
        LOG.exception('Exception occurs during CLI invocation')
        # This message always appears in stdout.
        # The verbose flag does not affect this.
        click.echo(f'Error: {str(e)}')

        if not is_verbose_enabled():
            click.echo('For more information rerun command with --verbose flag')

        # We can't reraise, for example, click.ClickException.
        # Because this function doesn't executes in click context,
        # so click library doesn't handle the exception.
        sys.exit(1)
