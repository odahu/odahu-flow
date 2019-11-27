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
import logging
import os.path
import subprocess

LOGGER = logging.getLogger(__name__)


def setup_logging(verbose: bool = False) -> None:
    """
    Setup logging instance

    :param verbose: use verbose output
    :type verbose: bool
    """
    log_level = logging.DEBUG if verbose else logging.INFO

    logging.basicConfig(format='[odahuflow][%(levelname)5s] %(asctime)s - %(message)s', datefmt='%d-%b-%y %H:%M:%S',
                        level=log_level)


def run(*args: str, cwd=None, stream_output: bool = True):
    """
    Run system command and stream / capture stdout and stderr

    :param args: Commands to run (list). E.g. ['cp, '-R', '/var/www', '/www']
    :type args: :py:class:`typing.List[str]`
    :param cwd: (Optional). Program working directory. Current is used if not provided.
    :type cwd: str
    :param stream_output: (Optional). Flag that enables streaming child process output to stdout and stderr.
    :return: typing.Union[int, typing.Tuple[int, str, str]] -- exit_code (for stream_output mode)
             or exit_code + stdout + stderr.
    """
    args_line = ' '.join(args)
    logging.info(f'Running command "{args_line}"')

    cmd_env = os.environ.copy()
    if stream_output:
        child = subprocess.Popen(args, env=cmd_env, cwd=cwd, universal_newlines=True,
                                 stdin=subprocess.PIPE)
        exit_code = child.wait()
        if exit_code != 0:
            raise Exception("Non-zero exitcode: %s" % exit_code)
        return exit_code
    else:
        child = subprocess.Popen(
            args, env=cmd_env, stdout=subprocess.PIPE, stdin=subprocess.PIPE, stderr=subprocess.PIPE,
            cwd=cwd, universal_newlines=True)
        (stdout, stderr) = child.communicate()
        exit_code = child.wait()
        if exit_code != 0:
            raise Exception("Non-zero exit code: %s\n\nSTDOUT:\n%s\n\nSTDERR:%s\n ======" %
                            (exit_code, stdout, stderr))
        return exit_code, stdout, stderr
