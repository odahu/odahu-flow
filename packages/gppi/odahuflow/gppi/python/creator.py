import glob
import logging
import shutil
import subprocess
from pathlib import Path

from odahuflow.gppi.model.installer import ModelCreator
from odahuflow.gppi.model.meta import Meta

LOG = logging.getLogger(__name__)


class PythonModelCreator(ModelCreator):

    def __init__(self, meta: Meta, output_artifact: Path, library: Path):
        super().__init__(meta, output_artifact)
        self.library = library

    def create_language_specific_part(self, artifact_dir: Path):
        shutil.copyfile(str(self.library), str(artifact_dir))


def create_python_library(source_code: Path, output: Path):
    """
    Create a python package using source distribution format
    :param source_code: path to the source directory. This dir must contain setup.py file.
    :param output: path to a file where the result package will be saved
    """
    dist_dir = source_code / "dist"
    if dist_dir.exists():
        LOG.debug(f'Detected the dist dir {dist_dir}. Removed it.')

        shutil.rmtree(dist_dir)

    # build the source distribution of the package
    # the package will be stored in the source_code/dist package
    subprocess.run(['python', 'setup.py', 'sdist'], cwd=source_code)

    # the result file name contains package name, version and so on.
    files = glob.glob(str(dist_dir / '*.tar.gz'))

    if not files:
        raise ValueError(f"The result dir {source_code} doesn't contain the package")
    if len(files) > 1:
        raise ValueError(f"The result dir {source_code} contain than one package")

    package = files[0]
    shutil.copyfile(package, output)
