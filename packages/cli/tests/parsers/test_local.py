from os.path import isdir
from unittest.mock import patch
from odahuflow.sdk.local.training import list_local_trainings


def test_list_local_trainings(tmpdir):
    mock_dir = tmpdir
    folders = ['wine-1@-12',
               '&wine-1-12',
               '[wine-1@-12',
               '2Wine-1@-12',
               'Wine-1@-123',
               '@wine-1@-12',
               'Zine-1@-12',
               'zine-1@-12',
               'Awine-1@']
    for folder in folders:
        mock_dir.mkdir(folder)
    assert isdir(mock_dir)
    with patch('odahuflow.sdk.local.training.config.LOCAL_MODEL_OUTPUT_DIR', mock_dir):
        # sorting based on ASCII sort order https://en.wikipedia.org/wiki/ASCII
        assert list_local_trainings() == ['&wine-1-12',
                                          '2Wine-1@-12',
                                          '@wine-1@-12',
                                          'Awine-1@',
                                          'Wine-1@-123',
                                          'Zine-1@-12',
                                          '[wine-1@-12',
                                          'wine-1@-12',
                                          'zine-1@-12']