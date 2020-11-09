import base64
from typing import Any, Dict, List, Tuple
from unittest.mock import Mock

from dataclasses import dataclass
import pytest
import yaml
from click.testing import CliRunner
from odahuflow.cli.parsers.local.packaging import run
from odahuflow.sdk.clients.api import WrongHttpStatusCode
from odahuflow.sdk.local import packaging as packaging_sdk
from odahuflow.sdk.models import ConnectionSpec, K8sPackager, ModelPackaging, ModelPackagingSpec, ModelPackagingStatus, \
    PackagingIntegration, \
    PackagerTarget, \
    PackagingIntegrationSpec, Schema, Target, \
    Connection, TargetSchema
from pytest_mock import MockFixture

PLAIN_VALUE = 'Value'

conn1 = Connection(
    id='conn1',
    spec=ConnectionSpec(
        key_id=base64.b64encode(bytes(PLAIN_VALUE, 'utf-8')).decode('utf-8'),
        key_secret=base64.b64encode(bytes(PLAIN_VALUE, 'utf-8')).decode('utf-8'),
        public_key=base64.b64encode(bytes(PLAIN_VALUE, 'utf-8')).decode('utf-8'),
        password=base64.b64encode(bytes(PLAIN_VALUE, 'utf-8')).decode('utf-8'),
        type='gcr'
    )
)
conn2 = Connection(
    id='conn2',
    spec=ConnectionSpec(
        key_id=base64.standard_b64encode(bytes(PLAIN_VALUE, 'utf-8')).decode('utf-8'),
        key_secret=base64.b64encode(bytes(PLAIN_VALUE, 'utf-8')).decode('utf-8'),
        public_key=base64.b64encode(bytes(PLAIN_VALUE, 'utf-8')).decode('utf-8'),
        password=base64.b64encode(bytes(PLAIN_VALUE, 'utf-8')).decode('utf-8'),
        type='gcr'
    )
)
conn_d1 = Connection(
    id='conn1',
    spec=ConnectionSpec(
        key_id=PLAIN_VALUE,
        key_secret=PLAIN_VALUE,
        public_key=PLAIN_VALUE,
        password=PLAIN_VALUE,
        type='gcr'
    )
)
conn_d2 = Connection(
    id='conn2',
    spec=ConnectionSpec(
        key_id=PLAIN_VALUE,
        key_secret=PLAIN_VALUE,
        public_key=PLAIN_VALUE,
        password=PLAIN_VALUE,
        type='gcr'
    )
)
target_pull = Target("conn1", "docker-pull")
target_push = Target("conn2", "docker-push")
pi_target_pull = PackagerTarget(conn_d1, "docker-pull")
pi_target_push = PackagerTarget(conn_d2, "docker-push")


pack1 = ModelPackaging(
    id='pack1',
    spec=ModelPackagingSpec(
        integration_name="pi",
        targets=[target_pull, target_push]
    ),
    status=ModelPackagingStatus()
)

pi = PackagingIntegration(id="pi", spec=PackagingIntegrationSpec(
    schema=Schema(
        targets=[
            TargetSchema(
                connection_types=['gcr', 'ecr', 'docker'], name='docker-push', required=True
            ),
            TargetSchema(
                connection_types=['gcr', 'ecr', 'docker'], name='docker-pull', required=False
            )
        ]
    )
))


@dataclass
class I:
    cmd: List[str]
    local: List[Any]
    remote: List[Any]


@dataclass
class E:
    targets: List[PackagerTarget]
    exit_code: int = 0
    exc: BaseException = None
    output_subs: str = ""


@dataclass
class Case:
    input: I
    expected: E


test_cases: List[Case] = [
    Case(  # all manifests are locally stored, default run without options
        input=I(
            cmd=["--pack-id", "pack1"],
            local=[conn1, conn2, pi, pack1], remote=[]),
        expected=E(targets=[])
    ),
    Case(  # no global targets disable
        input=I(
            cmd=["--pack-id", "pack1", "--no-disable-package-targets"],
            local=[conn1, conn2, pi, pack1], remote=[]),
        expected=E(targets=[pi_target_pull, pi_target_push])
    ),
    Case(  # no global targets disable, but specific target is disabled
        input=I(
            cmd=[
                "--pack-id", "pack1", "--no-disable-package-targets",
                "--disable-target", "docker-pull"
            ],
            local=[conn1, conn2, pi, pack1], remote=[]),
        expected=E(targets=[pi_target_push])
    ),
    Case(  # manifests on the server
        input=I(
            cmd=["--pack-id", "pack1"],
            local=[], remote=[conn1, conn2, pi, pack1]),
        expected=E(targets=[])
    ),
    Case(  # some manifests are local stored, some in the server
        input=I(
            cmd=["--pack-id", "pack1"],
            local=[pi, conn2], remote=[conn1, pack1]),
        expected=E(targets=[])
    ),
    Case(  # pack not found
        input=I(
            cmd=["--pack-id", "pack1"],
            local=[], remote=[conn1, conn2, pi]),
        expected=E(targets=[], exit_code=1, exc=WrongHttpStatusCode(404, {"message": f"Not found {pack1.id}"}))
    ),
    Case(  # no global targets disable, conn missed
        input=I(
            cmd=["--pack-id", "pack1", "--no-disable-package-targets"],
            local=[conn1, pi, pack1], remote=[]),
        expected=E(targets=[], exit_code=1, exc=SystemExit(1,),
                   output_subs="\"conn2\" connection of \"docker-push\" target is not found")
    ),
]


@pytest.mark.parametrize("test_case", test_cases)
def test_run_targets(cli_runner: CliRunner, test_case: Case, mocker: MockFixture):

    input_d, expected = test_case.input, test_case.expected
    command_args, local, remote = input_d.cmd, input_d.local, input_d.remote

    # Prepare remote API client mocks
    api_client = Mock()

    def api_retrieve(id_: str):
        for en in remote:
            if id_ == en.id:
                return en
        raise WrongHttpStatusCode(404, {"message": f"Not found {id_}"})

    api_client.get = Mock(side_effect=api_retrieve)
    api_client.get_decrypted = Mock(side_effect=api_retrieve)

    mocker.patch(
        "odahuflow.cli.parsers.local.packaging.PackagingIntegrationClient.construct_from_other",
        new=Mock(return_value=api_client)
    )
    mocker.patch(
        "odahuflow.cli.parsers.local.packaging.ConnectionClient.construct_from_other",
        new=Mock(return_value=api_client)
    )

    # Prepare entities that are stored locally
    with cli_runner.isolated_filesystem():
        with open('manifest.yaml', 'w') as f:
            docs = []
            for en in local:
                d = en.to_dict()
                d["kind"] = en.__class__.__name__
                docs.append(d)
            yaml.safe_dump_all(docs, f)
        if "--manifest-file" not in command_args:
            command_args += ["--manifest-file", "manifest.yaml"]

        # Run
        m: Mock = mocker.patch("odahuflow.cli.parsers.local.packaging.start_package")
        result = cli_runner.invoke(run, command_args, obj=api_client)

        # Assert expectations
        if expected.exit_code == 0:
            assert result.exit_code == 0
            assert result.exception is None

            m.assert_called_once()
            mp_json_actual, _ = m.call_args[0]  # type: K8sPackager, Any

            assert mp_json_actual.targets == expected.targets

        else:
            assert result.exit_code == 1
            assert isinstance(result.exception, type(expected.exc))
            assert str(expected.exc) == str(result.exception)

        if expected.output_subs:
            assert expected.output_subs in result.output
