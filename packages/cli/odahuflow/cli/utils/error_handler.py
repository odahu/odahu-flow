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
from typing import Optional

import click

from odahuflow.sdk.logger import is_verbose_enabled

LOG = logging.getLogger(__name__)

TIMEOUT_ERROR_MESSAGE = 'Time out: operation has not been confirmed'
IGNORE_NOT_FOUND_ERROR_MESSAGE = '{kind} {id} was not found. Ignore'
ID_AND_FILE_GIVEN_ERROR_MESSAGE = 'You should provide an ID or ' \
                                   'file parameter, not both.'
ID_AND_FILE_MISSED_ERROR_MESSAGE = 'You should provide an ID or ' \
                                   'file parameter.'


def check_id_or_file_params_present(
        entity_id: Optional[str], file: Optional[str]
) -> None:
    """
    Verify that only one parameter is present
    """
    if not entity_id and not file:
        raise ValueError(ID_AND_FILE_MISSED_ERROR_MESSAGE)

    if entity_id and file:
        raise ValueError(ID_AND_FILE_GIVEN_ERROR_MESSAGE)


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
