import abc
from typing import Any, Optional, Tuple, List, Union

import pandas as pd
from odahuflow.gppi.model.meta import Meta
from odahuflow.gppi.model.schema import ModelSchemas

ModelInputPredict = Union[List[List[Any]], pd.Dataframe]
ModelOutputPredict = Tuple[List[List[Any]], Tuple[str, ...]]


class Model(abc.ABC):

    def __init__(self, meta: Meta, schemas: Optional[ModelSchemas] = None):
        self._meta = meta
        self._schemas = schemas

    def info(self) -> Meta:
        """
        Get input and output schemas
        :return: OpenAPI specifications. Each specification is assigned as (input / output)
        """
        return self._meta

    @property
    @abc.abstractmethod
    def raw_model(self) -> Any:
        """
        Extract raw model
        :param data:
        """
        pass

    @property
    def schemas(self) -> Optional[ModelSchemas]:
        return self._schemas

    @abc.abstractmethod
    def predict(
            self,
            input_matrix: ModelInputPredict,
            provided_columns_names: Optional[List[str]] = None,
    ) -> ModelOutputPredict:
        """
        Make a prediction on a Matrix of values

        :param input_matrix: data for prediction
        :param provided_columns_names: Name of columns for provided matrix
        :return: result matrix and result column names
        """
        pass

    def verify(self):
        """
        Check GPPI correctness.
        By default, it has no implementation.
        Must raise an exception if verifying will fail
        """
        pass
