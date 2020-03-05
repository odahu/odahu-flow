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
from click.testing import CliRunner
from odahuflow.cli.parsers import template


def test_get_all_template_names(cli_runner: CliRunner):
    result = cli_runner.invoke(template.template, ['all'])

    assert result.exit_code == 0
    assert result.output
    assert "* deployment\n" in result.output
    assert "* mlflow_training\n" in result.stdout
    assert "* gcp_connection\n" in result.stdout


def test_generate_template_by_name(cli_runner: CliRunner):
    result = cli_runner.invoke(template.template, ['generate', '--name', 'deployment'])

    assert result.exit_code == 0
    assert result.output
    assert "image: <model/image:tag>" in result.output
    assert "kind: ModelDeployment" in result.output


def test_generate_not_present_template(cli_runner: CliRunner):
    result = cli_runner.invoke(template.template, ['generate', '--name', 'not_present'])

    assert result.exit_code != 0
    assert result.exception
    assert "Cannot find not_present template" in str(result.exception)
