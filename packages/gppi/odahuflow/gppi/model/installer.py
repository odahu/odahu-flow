import abc
import logging
import tempfile
from pathlib import Path

from odahuflow.gppi.model.meta import Meta
from odahuflow.gppi.utils.zip import zip_dir

LOG = logging.getLogger(__name__)

MODEL_DIR_NAME = 'odahuflow_model'


class ModelCreator(abc.ABC):

    def __init__(self, meta: Meta, output_artifact: Path):
        self.meta = meta
        self.output_artifact = output_artifact

    @abc.abstractmethod
    def create_language_specific_part(self, artifact_dir: Path):
        pass

    def create(self):
        with tempfile.TemporaryDirectory() as temp_artifact_dir:
            temp_artifact_dir = Path(temp_artifact_dir)
            model_dir = temp_artifact_dir / MODEL_DIR_NAME
            model_dir.mkdir()

            self.create_language_specific_part(temp_artifact_dir)

            self.meta.dump_to_file(temp_artifact_dir)

            self.meta.dependencies.dump_to_file(temp_artifact_dir)

            zip_dir(temp_artifact_dir, self.output_artifact)


class ModelInstaller(abc.ABC):

    @abc.abstractmethod
    def install(self):
        pass
