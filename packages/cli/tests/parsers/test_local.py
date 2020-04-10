import os
from os.path import isdir
import pdb

import pytest
from _pytest import tmpdir
from pytest_mock import mocker
import pdb
from unittest.mock import patch

from odahuflow.sdk.local.training import list_local_trainings


# def test_list_local_trainings(tmp_path):
#     # pdb.set_trace()
#     with patch('odahuflow.sdk.local.training.config.LOCAL_MODEL_OUTPUT_DIR') as mock_dir:
#         mock_dir = tmp_path
#         mock_dir.mkdir('wine-132121-fine')
#         # pdb.set_trace()
#         assert isdir(mock_dir)
#         assert len(os.listdir(mock_dir)) == 1
#         assert list_local_trainings() == ['wine-132121-fine']

def test_list_local_trainings(tmpdir):
    mock_dir = tmpdir
    mock_dir.mkdir('wine-132121-fine')
    # pdb.set_trace()
    assert isdir(mock_dir)
    assert len(os.listdir(mock_dir)) == 1
    with patch('odahuflow.sdk.local.training.config.LOCAL_MODEL_OUTPUT_DIR', mock_dir):
        assert list_local_trainings() == ['wine-132121-fine']