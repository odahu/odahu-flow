from typing import Optional, List, Any

from odahuflow.gppi.model.meta import Meta
from odahuflow.gppi.model.model import Model, ModelInputPredict, ModelOutputPredict


class StubPythonModel(Model):
    """
    Example model
    """
    def __init__(self):
        super().__init__(Meta())

    @property
    def raw_model(self) -> Any:
        return {}

    def predict(self, input_matrix: ModelInputPredict,
                provided_columns_names: Optional[List[str]] = None) -> ModelOutputPredict:
        return input_matrix

