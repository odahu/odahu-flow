import shutil
from os.path import dirname, join

from odahuflow.gppi.python.creator import create_python_library


def test_create_wheel_library(tmpdir):
    output_artifact = tmpdir / 'output'
    output_artifact.mkdir()
    output_artifact /= 'library.tar.gz'

    library_tmp = tmpdir / 'library'
    library = join(dirname(__file__), 'library')

    shutil.copytree(library, str(library_tmp))

    create_python_library(library_tmp, output_artifact)

    assert output_artifact.exists()
    assert output_artifact.isfile()
