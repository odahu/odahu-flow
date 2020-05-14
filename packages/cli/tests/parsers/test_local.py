from unittest.mock import patch
from odahuflow.sdk.local.training import list_local_trainings


def test_list_local_trainings(tmpdir):
    folders = ['wine-1@-12',
               '&wine-1-12',
               '[wine-1@-12',
               '2Wine-1@-12',
               'Wine-1@-123',
               '@wine-1@-12',
               'zine-1@-12',
               'Awine-1@']
    for folder in folders:
        tmpdir.mkdir(folder)
    with patch('odahuflow.sdk.local.training.config.LOCAL_MODEL_OUTPUT_DIR', tmpdir):
        assert list_local_trainings() == ['&wine-1-12',
                                          '2Wine-1@-12',
                                          '@wine-1@-12',
                                          'Awine-1@',
                                          'Wine-1@-123',
                                          '[wine-1@-12',
                                          'wine-1@-12',
                                          'zine-1@-12']