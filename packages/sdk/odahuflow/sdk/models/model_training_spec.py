# coding: utf-8

from __future__ import absolute_import
from datetime import date, datetime  # noqa: F401

from typing import List, Dict  # noqa: F401

from odahuflow.sdk.models.base_model_ import Model
from odahuflow.sdk.models.data_binding_dir import DataBindingDir  # noqa: F401,E501
from odahuflow.sdk.models.environment_variable import EnvironmentVariable  # noqa: F401,E501
from odahuflow.sdk.models.model_identity import ModelIdentity  # noqa: F401,E501
from odahuflow.sdk.models.resource_requirements import ResourceRequirements  # noqa: F401,E501
from odahuflow.sdk.models import util


class ModelTrainingSpec(Model):
    """NOTE: This class is auto generated by the swagger code generator program.

    Do not edit the class manually.
    """

    def __init__(self, args: List[str]=None, data: List[DataBindingDir]=None, entrypoint: str=None, envs: List[EnvironmentVariable]=None, hyper_parameters: Dict[str, str]=None, image: str=None, model: ModelIdentity=None, output_connection: str=None, reference: str=None, resources: ResourceRequirements=None, toolchain: str=None, vcs_name: str=None, work_dir: str=None):  # noqa: E501
        """ModelTrainingSpec - a model defined in Swagger

        :param args: The args of this ModelTrainingSpec.  # noqa: E501
        :type args: List[str]
        :param data: The data of this ModelTrainingSpec.  # noqa: E501
        :type data: List[DataBindingDir]
        :param entrypoint: The entrypoint of this ModelTrainingSpec.  # noqa: E501
        :type entrypoint: str
        :param envs: The envs of this ModelTrainingSpec.  # noqa: E501
        :type envs: List[EnvironmentVariable]
        :param hyper_parameters: The hyper_parameters of this ModelTrainingSpec.  # noqa: E501
        :type hyper_parameters: Dict[str, str]
        :param image: The image of this ModelTrainingSpec.  # noqa: E501
        :type image: str
        :param model: The model of this ModelTrainingSpec.  # noqa: E501
        :type model: ModelIdentity
        :param output_connection: The output_connection of this ModelTrainingSpec.  # noqa: E501
        :type output_connection: str
        :param reference: The reference of this ModelTrainingSpec.  # noqa: E501
        :type reference: str
        :param resources: The resources of this ModelTrainingSpec.  # noqa: E501
        :type resources: ResourceRequirements
        :param toolchain: The toolchain of this ModelTrainingSpec.  # noqa: E501
        :type toolchain: str
        :param vcs_name: The vcs_name of this ModelTrainingSpec.  # noqa: E501
        :type vcs_name: str
        :param work_dir: The work_dir of this ModelTrainingSpec.  # noqa: E501
        :type work_dir: str
        """
        self.swagger_types = {
            'args': List[str],
            'data': List[DataBindingDir],
            'entrypoint': str,
            'envs': List[EnvironmentVariable],
            'hyper_parameters': Dict[str, str],
            'image': str,
            'model': ModelIdentity,
            'output_connection': str,
            'reference': str,
            'resources': ResourceRequirements,
            'toolchain': str,
            'vcs_name': str,
            'work_dir': str
        }

        self.attribute_map = {
            'args': 'args',
            'data': 'data',
            'entrypoint': 'entrypoint',
            'envs': 'envs',
            'hyper_parameters': 'hyperParameters',
            'image': 'image',
            'model': 'model',
            'output_connection': 'outputConnection',
            'reference': 'reference',
            'resources': 'resources',
            'toolchain': 'toolchain',
            'vcs_name': 'vcsName',
            'work_dir': 'workDir'
        }

        self._args = args
        self._data = data
        self._entrypoint = entrypoint
        self._envs = envs
        self._hyper_parameters = hyper_parameters
        self._image = image
        self._model = model
        self._output_connection = output_connection
        self._reference = reference
        self._resources = resources
        self._toolchain = toolchain
        self._vcs_name = vcs_name
        self._work_dir = work_dir

    @classmethod
    def from_dict(cls, dikt) -> 'ModelTrainingSpec':
        """Returns the dict as a model

        :param dikt: A dict.
        :type: dict
        :return: The ModelTrainingSpec of this ModelTrainingSpec.  # noqa: E501
        :rtype: ModelTrainingSpec
        """
        return util.deserialize_model(dikt, cls)

    @property
    def args(self) -> List[str]:
        """Gets the args of this ModelTrainingSpec.


        :return: The args of this ModelTrainingSpec.
        :rtype: List[str]
        """
        return self._args

    @args.setter
    def args(self, args: List[str]):
        """Sets the args of this ModelTrainingSpec.


        :param args: The args of this ModelTrainingSpec.
        :type args: List[str]
        """

        self._args = args

    @property
    def data(self) -> List[DataBindingDir]:
        """Gets the data of this ModelTrainingSpec.

        Input data for a training  # noqa: E501

        :return: The data of this ModelTrainingSpec.
        :rtype: List[DataBindingDir]
        """
        return self._data

    @data.setter
    def data(self, data: List[DataBindingDir]):
        """Sets the data of this ModelTrainingSpec.

        Input data for a training  # noqa: E501

        :param data: The data of this ModelTrainingSpec.
        :type data: List[DataBindingDir]
        """

        self._data = data

    @property
    def entrypoint(self) -> str:
        """Gets the entrypoint of this ModelTrainingSpec.

        Model training file. It can be python\\bash script or jupiter notebook  # noqa: E501

        :return: The entrypoint of this ModelTrainingSpec.
        :rtype: str
        """
        return self._entrypoint

    @entrypoint.setter
    def entrypoint(self, entrypoint: str):
        """Sets the entrypoint of this ModelTrainingSpec.

        Model training file. It can be python\\bash script or jupiter notebook  # noqa: E501

        :param entrypoint: The entrypoint of this ModelTrainingSpec.
        :type entrypoint: str
        """

        self._entrypoint = entrypoint

    @property
    def envs(self) -> List[EnvironmentVariable]:
        """Gets the envs of this ModelTrainingSpec.

        Custom environment variables that should be set before entrypoint invocation.  # noqa: E501

        :return: The envs of this ModelTrainingSpec.
        :rtype: List[EnvironmentVariable]
        """
        return self._envs

    @envs.setter
    def envs(self, envs: List[EnvironmentVariable]):
        """Sets the envs of this ModelTrainingSpec.

        Custom environment variables that should be set before entrypoint invocation.  # noqa: E501

        :param envs: The envs of this ModelTrainingSpec.
        :type envs: List[EnvironmentVariable]
        """

        self._envs = envs

    @property
    def hyper_parameters(self) -> Dict[str, str]:
        """Gets the hyper_parameters of this ModelTrainingSpec.

        Model training hyperParameters in parameter:value format  # noqa: E501

        :return: The hyper_parameters of this ModelTrainingSpec.
        :rtype: Dict[str, str]
        """
        return self._hyper_parameters

    @hyper_parameters.setter
    def hyper_parameters(self, hyper_parameters: Dict[str, str]):
        """Sets the hyper_parameters of this ModelTrainingSpec.

        Model training hyperParameters in parameter:value format  # noqa: E501

        :param hyper_parameters: The hyper_parameters of this ModelTrainingSpec.
        :type hyper_parameters: Dict[str, str]
        """

        self._hyper_parameters = hyper_parameters

    @property
    def image(self) -> str:
        """Gets the image of this ModelTrainingSpec.

        Train image  # noqa: E501

        :return: The image of this ModelTrainingSpec.
        :rtype: str
        """
        return self._image

    @image.setter
    def image(self, image: str):
        """Sets the image of this ModelTrainingSpec.

        Train image  # noqa: E501

        :param image: The image of this ModelTrainingSpec.
        :type image: str
        """

        self._image = image

    @property
    def model(self) -> ModelIdentity:
        """Gets the model of this ModelTrainingSpec.

        Model Identity  # noqa: E501

        :return: The model of this ModelTrainingSpec.
        :rtype: ModelIdentity
        """
        return self._model

    @model.setter
    def model(self, model: ModelIdentity):
        """Sets the model of this ModelTrainingSpec.

        Model Identity  # noqa: E501

        :param model: The model of this ModelTrainingSpec.
        :type model: ModelIdentity
        """

        self._model = model

    @property
    def output_connection(self) -> str:
        """Gets the output_connection of this ModelTrainingSpec.

        Name of Connection to storage where training output artifact will be stored. Permitted connection types are defined by specific toolchain  # noqa: E501

        :return: The output_connection of this ModelTrainingSpec.
        :rtype: str
        """
        return self._output_connection

    @output_connection.setter
    def output_connection(self, output_connection: str):
        """Sets the output_connection of this ModelTrainingSpec.

        Name of Connection to storage where training output artifact will be stored. Permitted connection types are defined by specific toolchain  # noqa: E501

        :param output_connection: The output_connection of this ModelTrainingSpec.
        :type output_connection: str
        """

        self._output_connection = output_connection

    @property
    def reference(self) -> str:
        """Gets the reference of this ModelTrainingSpec.

        VCS Reference  # noqa: E501

        :return: The reference of this ModelTrainingSpec.
        :rtype: str
        """
        return self._reference

    @reference.setter
    def reference(self, reference: str):
        """Sets the reference of this ModelTrainingSpec.

        VCS Reference  # noqa: E501

        :param reference: The reference of this ModelTrainingSpec.
        :type reference: str
        """

        self._reference = reference

    @property
    def resources(self) -> ResourceRequirements:
        """Gets the resources of this ModelTrainingSpec.

        Resources for model container The same format like k8s uses for pod resources.  # noqa: E501

        :return: The resources of this ModelTrainingSpec.
        :rtype: ResourceRequirements
        """
        return self._resources

    @resources.setter
    def resources(self, resources: ResourceRequirements):
        """Sets the resources of this ModelTrainingSpec.

        Resources for model container The same format like k8s uses for pod resources.  # noqa: E501

        :param resources: The resources of this ModelTrainingSpec.
        :type resources: ResourceRequirements
        """

        self._resources = resources

    @property
    def toolchain(self) -> str:
        """Gets the toolchain of this ModelTrainingSpec.

        IntegrationName of toolchain  # noqa: E501

        :return: The toolchain of this ModelTrainingSpec.
        :rtype: str
        """
        return self._toolchain

    @toolchain.setter
    def toolchain(self, toolchain: str):
        """Sets the toolchain of this ModelTrainingSpec.

        IntegrationName of toolchain  # noqa: E501

        :param toolchain: The toolchain of this ModelTrainingSpec.
        :type toolchain: str
        """

        self._toolchain = toolchain

    @property
    def vcs_name(self) -> str:
        """Gets the vcs_name of this ModelTrainingSpec.

        Name of Connection resource. Must exists  # noqa: E501

        :return: The vcs_name of this ModelTrainingSpec.
        :rtype: str
        """
        return self._vcs_name

    @vcs_name.setter
    def vcs_name(self, vcs_name: str):
        """Sets the vcs_name of this ModelTrainingSpec.

        Name of Connection resource. Must exists  # noqa: E501

        :param vcs_name: The vcs_name of this ModelTrainingSpec.
        :type vcs_name: str
        """

        self._vcs_name = vcs_name

    @property
    def work_dir(self) -> str:
        """Gets the work_dir of this ModelTrainingSpec.

        Directory with model scripts/files in a git repository  # noqa: E501

        :return: The work_dir of this ModelTrainingSpec.
        :rtype: str
        """
        return self._work_dir

    @work_dir.setter
    def work_dir(self, work_dir: str):
        """Sets the work_dir of this ModelTrainingSpec.

        Directory with model scripts/files in a git repository  # noqa: E501

        :param work_dir: The work_dir of this ModelTrainingSpec.
        :type work_dir: str
        """

        self._work_dir = work_dir
