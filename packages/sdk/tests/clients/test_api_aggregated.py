import json
import logging
import tempfile
from io import StringIO

import pytest
import yaml
from odahuflow.sdk.clients.api_aggregated import parse_resources_file, parse_stream, parse_resources_dir
from odahuflow.sdk.models import Connection, ConnectionSpec

NOT_EXISTED_FILE = "not-existed-file"
NOT_VALID_YAML_OR_JSON = "{not': valid-yaml-or-json"
CONNECTION_RES = Connection(id='test', spec=ConnectionSpec(type="git"))

logging.basicConfig(level=logging.DEBUG)

# Because of fixtures are not recognized by pylint as special kind of variables
# pylint: disable=redefined-outer-name


@pytest.fixture()
def conn_manifest_file() -> str:
    """
    Returns path to a temporary Connection file
    """
    with tempfile.NamedTemporaryFile(mode="w") as temp_file:
        temp_file.write(json.dumps({**CONNECTION_RES.to_dict(), **{"kind": "Connection"}}))
        temp_file.flush()

        yield temp_file.name


@pytest.fixture()
def conn_manifest_yaml_file() -> str:
    """
    Returns path to a temporary Connection YAML file
    """
    with tempfile.NamedTemporaryFile(mode="w") as temp_file:
        temp_file.write(yaml.dump({**CONNECTION_RES.to_dict(), **{"kind": "Connection"}}))
        temp_file.flush()

        yield temp_file.name


@pytest.fixture()
def conn_manifest_dir() -> str:
    """
    Returns path to a temporary dir with Connection files
    """
    with tempfile.TemporaryDirectory() as temp_dir:
        with tempfile.NamedTemporaryFile(mode="w", dir=temp_dir) as temp_file:
            temp_file.write(json.dumps({**CONNECTION_RES.to_dict(), **{"kind": "Connection"}}))
            temp_file.flush()

            yield temp_dir


def test_parse_not_existed_resource_file():
    with pytest.raises(FileNotFoundError, match=fr".*'{NOT_EXISTED_FILE}' not found.*"):
        parse_resources_file(NOT_EXISTED_FILE)


def test_parse_not_existed_resources_file_with_one_item():
    with pytest.raises(FileNotFoundError, match=fr".*'{NOT_EXISTED_FILE}' not found.*"):
        parse_resources_file(NOT_EXISTED_FILE)


def test_parse_stream_not_valid_json_or_yaml():
    with pytest.raises(Exception, match=r'^not valid JSON or YAML$'):
        parse_stream(StringIO(NOT_VALID_YAML_OR_JSON))


def test_parse_stream_json_empty_array():
    assert len(parse_stream(StringIO("[]")).changes) == 0


def test_parse_stream_yaml_empty_array():
    assert len(parse_stream(StringIO("---\n")).changes) == 0


def test_parse_stream_empty():
    """An empty file treats like empty array"""
    assert len(parse_stream(StringIO("")).changes) == 0


def test_parse_resources_file(conn_manifest_file: str):
    result_conn = parse_resources_file(conn_manifest_file)

    assert len(result_conn.changes) == 1
    assert result_conn.changes[0].resource_id == CONNECTION_RES.id
    assert result_conn.changes[0].resource == CONNECTION_RES


def test_parse_resources_yaml_file(conn_manifest_yaml_file: str):
    result_conn = parse_resources_file(conn_manifest_yaml_file)

    assert len(result_conn.changes) == 1
    assert result_conn.changes[0].resource_id == CONNECTION_RES.id
    assert result_conn.changes[0].resource == CONNECTION_RES


def test_parse_dir_resources_file(conn_manifest_dir: str):
    result_conn = parse_resources_dir(conn_manifest_dir)

    assert len(result_conn) == 1
    assert result_conn[0].resource_id == CONNECTION_RES.id
    assert result_conn[0].resource == CONNECTION_RES
