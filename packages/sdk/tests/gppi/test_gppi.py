import logging
import os

import pytest
from odahuflow.sdk.gppi.executor import GPPITrainedModelBinary, VALIDATION_FAILED_EXCEPTION_MESSAGE

GPPI_MODEL_RESOURCES = 'packages/sdk/tests/gppi/resources'

logging.basicConfig(level=logging.INFO)


class TestGPPITrainedModelBinary:

    def test_validate_ok(self):
        model_path = os.path.join(GPPI_MODEL_RESOURCES, 'mlflow_gppi_correct')
        gppi_model = GPPITrainedModelBinary(model_path)
        gppi_model.self_check()

    def test_validate_lib_missed(self):

        model_path = os.path.join(GPPI_MODEL_RESOURCES, 'mlflow_gppi_lib_missed')
        print(model_path)
        gppi_model = GPPITrainedModelBinary(model_path)
        with pytest.raises(Exception) as exc_info:
            gppi_model.self_check()
        assert VALIDATION_FAILED_EXCEPTION_MESSAGE in str(exc_info.value)
