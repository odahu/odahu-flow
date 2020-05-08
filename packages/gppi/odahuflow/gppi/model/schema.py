from typing import Optional, List, Union, Type, Tuple, Any

import pandas as pd
import pandas.api.types as pdt
from pydantic import BaseModel


class SchemaProp(BaseModel):
    name: str
    type: Optional[str]
    example: Optional[Union[str, int, None]] = None
    required: bool = False


class ModelSchemas(BaseModel):
    input: List[SchemaProp]
    output: List[SchemaProp]

    @staticmethod
    def from_df(input: Optional[pd.DataFrame], output: Optional[pd.DataFrame]) -> 'ModelSchemas':
        return ModelSchemas(
            input=extract_df_properties(input),
            output=extract_df_properties(output),
        )


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


def extract_df_properties(df: pd.DataFrame) -> List[SchemaProp]:
    """
    Extract OpenAPI specification for pd.DataFrame columns

    :param df: pandas DataFrame
    :return: OpenAPI specification for parameters (each columns is parameter)
    """
    if df is None:
        return []

    props = []

    for column_name, column_type in df.dtypes.items():
        open_api_type, example = _type_to_open_api_format(column_type)

        props.append(SchemaProp(
            name=column_name,
            type=open_api_type,
            example=example,
            required=True,
        ))

    return props
