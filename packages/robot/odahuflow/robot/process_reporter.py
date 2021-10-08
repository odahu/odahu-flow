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
Listener that kills all active processes and outputs their stdout/stderr streams.
"""

import signal
import sys
import logging
import contextlib

from odahuflow.robot.libraries.framework_extensions import get_imported_library_instance

ROBOT_LISTENER_API_VERSION = 3
KILL_SIGNAL = signal.SIGKILL


LOGGER = logging.getLogger('listener.process_reporter')
LOGGER.setLevel(logging.DEBUG)

handler = logging.StreamHandler(sys.__stdout__)
handler.setLevel(logging.DEBUG)
handler.setFormatter(logging.Formatter('\n%(asctime)s - %(levelname)s - %(message)s'))
LOGGER.addHandler(handler)


class ExecutionTimeoutException(Exception):
    pass


@contextlib.contextmanager
def with_execution_time_limit(limit):
    """
    Send alarm after limit
    :param limit: limit alarm
    """
    def sig_handler(_1, _2):
        raise ExecutionTimeoutException()

    signal.signal(signal.SIGALRM, sig_handler)
    signal.alarm(limit)

    try:
        yield
    finally:
        signal.alarm(0)


def report_process_output(process_name, stream_name, stream):
    """
    Report process'es stream

    :param process_name: name of process
    :type process_name: str
    :param stream_name: name of stream, e.g. stdout
    :type stream_name: str
    :param stream: stream
    :type stream: bytes or str
    :return: None
    """
    if not stream:
        LOGGER.info(f'Process {process_name!r} has no {stream_name}')
    else:
        if isinstance(stream, bytes):
            stream = stream.decode('utf-8')
        LOGGER.info(f'Process {process_name!r} has {stream_name} output:\n{stream}')


def kill_and_report_process(popen_object):
    """
    Kill and report process stderr and stdout

    :param popen_object: instance of process
    :type popen_object: :py:class:`subprocess.Popen`
    :return: None
    """
    try:
        LOGGER.debug(f'Killing process #{popen_object.pid}')
        popen_object.kill()

        try:
            report_process_output(popen_object.args, 'stdout', popen_object.stdout.read())
            report_process_output(popen_object.args, 'stderr', popen_object.stderr.read())
        except Exception as gather_exception:
            LOGGER.error(f'Cannot gather process #{popen_object.pid} logs: {gather_exception!r}')
    except ProcessLookupError:
        LOGGER.error(f'Cannot find process by id #{popen_object.pid}')
    except Exception as kill_exception:
        LOGGER.error(f'Cannot kill process: {kill_exception!r}')


def end_test(test, result):  # pylint: disable=W0613
    """
    Listen Robot's "end of test" event

    :param test: test
    :param result: test result
    :return: None
    """
    if not result.passed:
        process_lib = get_imported_library_instance('Process')
        if process_lib:
            all_processes = process_lib._processes._connections
            all_results = process_lib._results

            active_processes = [process for process in all_processes
                                if process in all_results and all_results[process].rc is None]

            if active_processes:
                LOGGER.info(f'Some hanging processes have been detected for failed test {result.name!r}')

                for process in active_processes:
                    LOGGER.info(f'Killing active process {process.args!r} for test {result.name!r} '
                                f'because of {result.message!r}')
                    try:
                        with with_execution_time_limit(10):
                            kill_and_report_process(process)
                    except ExecutionTimeoutException:
                        LOGGER.error('Operation timed out')
