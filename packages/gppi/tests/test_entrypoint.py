import numpy as np
import pandas as pd
from odahuflow.gppi.model.schema import extract_df_properties, SchemaProp


def test_extract_df_properties():
    df = pd.DataFrame(
        {
            'A': 1.,
            'B': pd.Timestamp('20130102'),
            'C': pd.Series(1, index=list(range(4)), dtype='float32'),
            'D': np.array([3] * 4, dtype='int32'),
            'F': 'foo'
        }
    )

    # For now, we assume that the order of columns will be the same as in the input DataFrame
    assert extract_df_properties(df) == [
        SchemaProp(example=0, name='A', required=True, type='number'),
        SchemaProp(example=None, name='B', required=True, type=None),
        SchemaProp(example=0, name='C', required=True, type='number'),
        SchemaProp(example=0, name='D', required=True, type='integer'),
        SchemaProp(example='', name='F', required=True, type='string')
    ]


def test_extract_empty_df_properties():
    assert extract_df_properties(pd.DataFrame({})) == []
