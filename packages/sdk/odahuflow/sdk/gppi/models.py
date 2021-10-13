import typing

import pydantic


class OdahuflowProjectManifestBinaries(pydantic.BaseModel):
    """
    Odahuflow Project Manifest's Binaries description
    """

    type: str
    dependencies: str
    conda_path: typing.Optional[str]


class OdahuflowProjectManifestModel(pydantic.BaseModel):
    """
    Odahuflow Project Manifest's Model description
    """

    name: str
    version: str
    workDir: str
    entrypoint: str


class OdahuflowProjectManifestTrainingIntegration(pydantic.BaseModel):
    """
    Odahuflow Project Manifest's Training Integration description
    """

    name: str
    version: str


class OdahuflowProjectManifestOutput(pydantic.BaseModel):
    """
    Odahuflow Project Manifest's Output description
    """

    run_id: str


class OdahuflowProjectManifest(pydantic.BaseModel):
    """
    Odahuflow Project Manifest description class
    """

    binaries: OdahuflowProjectManifestBinaries
    model: typing.Optional[OdahuflowProjectManifestModel]
    odahuflowVersion: typing.Optional[str]
    training_integration: typing.Optional[OdahuflowProjectManifestTrainingIntegration]
    output: typing.Optional[OdahuflowProjectManifestOutput]
