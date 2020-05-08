from pathlib import Path
from typing import Any, Dict

import pydantic
import yaml
from odahuflow.gppi.model.dependencies import Dependencies

META_FILE_NAME = "odahuflow.project.yaml"


class ToolchainMeta(pydantic.BaseModel):
    name: str
    version: str


class ModelMeta(pydantic.BaseModel):
    name: str
    version: str


class Meta(pydantic.BaseModel):
    toolchain: ToolchainMeta
    model: ModelMeta
    output: Dict[str, Any]
    dependencies: Dependencies

    def dump_to_file(self, output_dir: Path):
        meta_file = output_dir / META_FILE_NAME
        with meta_file.open('w') as f:
            yaml.safe_dump(self.dict(), f)

    @staticmethod
    def read_from_file(output_dir: Path) -> 'Meta':
        meta_file = output_dir / META_FILE_NAME
        with meta_file.open('r') as f:
            return Meta(**yaml.safe_load(f))
