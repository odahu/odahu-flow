# coding: utf-8

from __future__ import absolute_import
from datetime import date, datetime  # noqa: F401

from typing import List, Dict  # noqa: F401

from odahuflow.sdk.models.base_model_ import Model
from odahuflow.sdk.models.auth_config import AuthConfig  # noqa: F401,E501
from odahuflow.sdk.models import util


class PackagerConfig(Model):
    """NOTE: This class is auto generated by the swagger code generator program.

    Do not edit the class manually.
    """

    def __init__(self, auth: AuthConfig=None, model_training_id: str=None, mp_file: str=None, output_training_dir: str=None):  # noqa: E501
        """PackagerConfig - a model defined in Swagger

        :param auth: The auth of this PackagerConfig.  # noqa: E501
        :type auth: AuthConfig
        :param model_training_id: The model_training_id of this PackagerConfig.  # noqa: E501
        :type model_training_id: str
        :param mp_file: The mp_file of this PackagerConfig.  # noqa: E501
        :type mp_file: str
        :param output_training_dir: The output_training_dir of this PackagerConfig.  # noqa: E501
        :type output_training_dir: str
        """
        self.swagger_types = {
            'auth': AuthConfig,
            'model_training_id': str,
            'mp_file': str,
            'output_training_dir': str
        }

        self.attribute_map = {
            'auth': 'auth',
            'model_training_id': 'modelTrainingId',
            'mp_file': 'mpFile',
            'output_training_dir': 'outputTrainingDir'
        }

        self._auth = auth
        self._model_training_id = model_training_id
        self._mp_file = mp_file
        self._output_training_dir = output_training_dir

    @classmethod
    def from_dict(cls, dikt) -> 'PackagerConfig':
        """Returns the dict as a model

        :param dikt: A dict.
        :type: dict
        :return: The PackagerConfig of this PackagerConfig.  # noqa: E501
        :rtype: PackagerConfig
        """
        return util.deserialize_model(dikt, cls)

    @property
    def auth(self) -> AuthConfig:
        """Gets the auth of this PackagerConfig.


        :return: The auth of this PackagerConfig.
        :rtype: AuthConfig
        """
        return self._auth

    @auth.setter
    def auth(self, auth: AuthConfig):
        """Sets the auth of this PackagerConfig.


        :param auth: The auth of this PackagerConfig.
        :type auth: AuthConfig
        """

        self._auth = auth

    @property
    def model_training_id(self) -> str:
        """Gets the model_training_id of this PackagerConfig.

        ID of the model packaging  # noqa: E501

        :return: The model_training_id of this PackagerConfig.
        :rtype: str
        """
        return self._model_training_id

    @model_training_id.setter
    def model_training_id(self, model_training_id: str):
        """Sets the model_training_id of this PackagerConfig.

        ID of the model packaging  # noqa: E501

        :param model_training_id: The model_training_id of this PackagerConfig.
        :type model_training_id: str
        """

        self._model_training_id = model_training_id

    @property
    def mp_file(self) -> str:
        """Gets the mp_file of this PackagerConfig.

        The path to the configuration file for a user packager.  # noqa: E501

        :return: The mp_file of this PackagerConfig.
        :rtype: str
        """
        return self._mp_file

    @mp_file.setter
    def mp_file(self, mp_file: str):
        """Sets the mp_file of this PackagerConfig.

        The path to the configuration file for a user packager.  # noqa: E501

        :param mp_file: The mp_file of this PackagerConfig.
        :type mp_file: str
        """

        self._mp_file = mp_file

    @property
    def output_training_dir(self) -> str:
        """Gets the output_training_dir of this PackagerConfig.

        The path to the dir when a user packager will save their result.  # noqa: E501

        :return: The output_training_dir of this PackagerConfig.
        :rtype: str
        """
        return self._output_training_dir

    @output_training_dir.setter
    def output_training_dir(self, output_training_dir: str):
        """Sets the output_training_dir of this PackagerConfig.

        The path to the dir when a user packager will save their result.  # noqa: E501

        :param output_training_dir: The output_training_dir of this PackagerConfig.
        :type output_training_dir: str
        """

        self._output_training_dir = output_training_dir
