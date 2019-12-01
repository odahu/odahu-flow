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
