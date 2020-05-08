import abc
import shutil
import subprocess
from pathlib import Path
from typing import Optional

import pydantic


class Dependency(abc.ABC):

    @abc.abstractmethod
    def install(self, artifact_dir: Path):
        pass

    # TODO: rename
    @abc.abstractmethod
    def dump_to_file(self, artifact_dir: Path):
        pass


class CondaDependencies(pydantic.BaseModel, Dependency):
    file_name: str = "conda.yaml"
    conda_env_name: str = "odahu_model"
    # TODO: looks weird....
    source_file_name: Optional[Path] = None

    def install(self, artifact_dir: Path):
        subprocess.run(
            ["conda", "env", "update", "-n", self.conda_env_name, "-f", self.file_name],
            cwd=str(artifact_dir.absolute()),
        )

    def dump_to_file(self, artifact_dir: Path):
        shutil.copyfile(str(self.source_file_name), (artifact_dir / self.file_name))


class Dependencies(pydantic.BaseModel, Dependency):
    conda: Optional[CondaDependencies]

    def install(self, artifact_dir: Path):
        for dep in [self.conda]:
            if dep:
                dep.install(artifact_dir)

    def dump_to_file(self, artifact_dir: Path):
        for dep in [self.conda]:
            if dep:
                dep.dump_to_file(artifact_dir=artifact_dir)
