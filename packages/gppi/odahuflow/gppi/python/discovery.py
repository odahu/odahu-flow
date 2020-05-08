from typing import List

from odahuflow.gppi.model.discovery import ModelDiscovery
from odahuflow.gppi.model.model import Model
from odahuflow.gppi.python.stub_model import StubPythonModel


class PythonModelDiscovery(ModelDiscovery):

    @property
    def models(self) -> List[Model]:
        # mgr = extension.ExtensionManager(
        #     namespace='odahuflow.models',
        #     invoke_on_load=True,
        # )

        return [StubPythonModel()]

    def get_model(self, entrypoint_model_name: str) -> Model:
        # mgr = extension.ExtensionManager(
        #     namespace='odahuflow.models',
        #     invoke_on_load=True,
        # )

        if entrypoint_model_name == "stub":
            return StubPythonModel()

        raise ValueError(f"{entrypoint_model_name} model not found")