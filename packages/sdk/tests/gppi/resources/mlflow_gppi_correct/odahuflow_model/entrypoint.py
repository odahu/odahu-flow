import functools
import json
import os
from typing import Optional, List, Dict, Union, Any, Tuple, Type

# MLFlow packages
import mlflow.models
import mlflow.pyfunc
# Third-party modules (is provided by MLFlow)
import numpy as np
import pandas as pd
import pandas.api.types as pdt

# Storage of loaded prediction function
MODEL_FLAVOR = None

# Path to model's root
MODEL_LOCATION = os.getenv('MODEL_LOCATION', '.')

# Optional. Examples of input and output pandas DataFrames
MODEL_INPUT_SAMPLE_FILE = os.path.join(MODEL_LOCATION, 'head_input.pkl')
MODEL_OUTPUT_SAMPLE_FILE = os.path.join(MODEL_LOCATION, 'head_output.pkl')


def _type_to_open_api_format(t: Type) -> Tuple[Optional[str], Optional[Any]]:
    """
    Convert type of column to OpenAPI type name and example

    :param t: object's type
    :return: name for OpenAPI
    """
    if isinstance(t, (str, bytes, bytearray)):
        return 'string', ''
    if isinstance(t, bool):
        return 'boolean', False
    if isinstance(t, int):
        return 'integer', 0
    if isinstance(t, float):
        return 'number', 0

    if pdt.is_integer_dtype(t):
        return 'integer', 0

    if pdt.is_float_dtype(t):
        return 'number', 0

    if pdt.is_string_dtype(t):
        return 'string', ''

    if pdt.is_bool_dtype(t) or pdt.is_complex_dtype(t):
        return 'string', ''

    return None, None


class NumpyEncoder(json.JSONEncoder):
    """
    Converts Numpy objects to Python's core objects
    """

    def default(self, o):
        if isinstance(o, np.generic):
            return o.item()
        return json.JSONEncoder.default(self, o)


def init() -> str:
    """
    Initialize model and return prediction type

    :return: prediction type (matrix or objects)
    """
    model = mlflow.models.Model.load(MODEL_LOCATION)
    if mlflow.pyfunc.FLAVOR_NAME not in model.flavors:
        raise ValueError('{} not in model\'s flavors'.format(mlflow.pyfunc.FLAVOR_NAME))

    global MODEL_FLAVOR
    MODEL_FLAVOR = mlflow.pyfunc.load_model(MODEL_LOCATION)
    return 'matrix'


def predict_on_matrix(input_matrix: List[List[Any]], provided_columns_names: Optional[List[str]] = None) \
        -> Tuple[List[List[Any]], Tuple[str, ...]]:
    """
    Make prediction on a Matrix of values

    :param input_matrix: data for prediction
    :param provided_columns_names: Name of columns for provided matrix
    :return: result matrix and result column names
    """
    if provided_columns_names:
        input_matrix = pd.DataFrame(input_matrix, columns=provided_columns_names)
    else:
        input_matrix = pd.DataFrame(input_matrix)

    input_sample = _input_df_sample()
    output_sample = _output_df_sample()

    if provided_columns_names and input_sample is not None:
        input_matrix = input_matrix.reindex(columns=input_sample.columns)

    result = MODEL_FLAVOR.predict(input_matrix)

    result_columns = []
    if output_sample is not None:
        result_columns = output_sample.columns

    # Register column names, overwrite if we've a sample
    if hasattr(result, 'columns'):
        result_columns = result.columns

    # TODO: think about better approach
    if isinstance(result, pd.DataFrame):
        output_matrix = result.to_numpy().tolist()
    elif isinstance(result, np.ndarray):
        output_matrix = result.tolist()
    else:
        output_matrix = result

    return output_matrix, tuple(result_columns)


@functools.lru_cache()
def _input_df_sample() -> Optional[pd.DataFrame]:
    """
    Internal function for getting input DataFrame sample

    :return: input sample if provided
    """
    if os.path.exists(MODEL_INPUT_SAMPLE_FILE):
        return pd.read_pickle(MODEL_INPUT_SAMPLE_FILE)
    else:
        return None


@functools.lru_cache()
def _output_df_sample() -> Optional[pd.DataFrame]:
    """
    Internal function for getting output DataFrame sample

    :return: input sample if provided
    """
    if os.path.exists(MODEL_OUTPUT_SAMPLE_FILE):
        return pd.read_pickle(MODEL_OUTPUT_SAMPLE_FILE)
    else:
        return None


def _extract_df_properties(df: pd.DataFrame) -> List[Dict[str, Union[Union[str, None, bool], Any]]]:
    """
    Extract OpenAPI specification for pd.DataFrame columns

    :param df: pandas DataFrame
    :return: OpenAPI specification for parameters (each columns is parameter)
    """
    if df is None:
        return []

    props = []

    for pos, column in enumerate(df.columns):
        open_api_type, example = _type_to_open_api_format(df.dtypes.array[pos])

        props.append({'name': column, 'type': open_api_type, 'example': example, 'required': True})

    return props


@functools.lru_cache()
def info() -> Tuple[List[Dict[str, Any]], List[Dict[str, Any]]]:
    """
    Get input and output schemas

    :return: OpenAPI specifications. Each specification is assigned as (input / output)
    """
    input_sample = _input_df_sample()
    output_sample = _output_df_sample()

    return _extract_df_properties(input_sample), _extract_df_properties(output_sample)


def get_output_json_serializer() -> type:
    """
    Returns JSON serializer to be used in output

    :return: JSON serializer
    """
    return NumpyEncoder
