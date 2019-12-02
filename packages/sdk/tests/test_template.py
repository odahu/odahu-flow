import pytest
from odahuflow.sdk.clients.templates import get_odahuflow_template_names, get_odahuflow_template_content


def test_list_of_templates():
    assert {"deployment", "mlflow_training", "gcp_connection"}.issubset(get_odahuflow_template_names())


def test_get_template_by_name():
    content = get_odahuflow_template_content("deployment")

    assert content
    assert "image: <model/image:tag>" in content
    assert "kind: ModelDeployment" in content


def test_get_template_not_found():
    with pytest.raises(ValueError):
        get_odahuflow_template_content("not_present")
