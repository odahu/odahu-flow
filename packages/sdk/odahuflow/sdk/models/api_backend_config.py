# coding: utf-8

from __future__ import absolute_import
from datetime import date, datetime  # noqa: F401

from typing import List, Dict  # noqa: F401

from odahuflow.sdk.models.base_model_ import Model
from odahuflow.sdk.models.api_local_backend_config import APILocalBackendConfig  # noqa: F401,E501
from odahuflow.sdk.models import util


class APIBackendConfig(Model):
    """NOTE: This class is auto generated by the swagger code generator program.

    Do not edit the class manually.
    """

    def __init__(self, local: APILocalBackendConfig=None, type: str=None):  # noqa: E501
        """APIBackendConfig - a model defined in Swagger

        :param local: The local of this APIBackendConfig.  # noqa: E501
        :type local: APILocalBackendConfig
        :param type: The type of this APIBackendConfig.  # noqa: E501
        :type type: str
        """
        self.swagger_types = {
            'local': APILocalBackendConfig,
            'type': str
        }

        self.attribute_map = {
            'local': 'local',
            'type': 'type'
        }

        self._local = local
        self._type = type

    @classmethod
    def from_dict(cls, dikt) -> 'APIBackendConfig':
        """Returns the dict as a model

        :param dikt: A dict.
        :type: dict
        :return: The APIBackendConfig of this APIBackendConfig.  # noqa: E501
        :rtype: APIBackendConfig
        """
        return util.deserialize_model(dikt, cls)

    @property
    def local(self) -> APILocalBackendConfig:
        """Gets the local of this APIBackendConfig.

        Local backend  # noqa: E501

        :return: The local of this APIBackendConfig.
        :rtype: APILocalBackendConfig
        """
        return self._local

    @local.setter
    def local(self, local: APILocalBackendConfig):
        """Sets the local of this APIBackendConfig.

        Local backend  # noqa: E501

        :param local: The local of this APIBackendConfig.
        :type local: APILocalBackendConfig
        """

        self._local = local

    @property
    def type(self) -> str:
        """Gets the type of this APIBackendConfig.

        Type of the backend. Available values:    * local    * config  # noqa: E501

        :return: The type of this APIBackendConfig.
        :rtype: str
        """
        return self._type

    @type.setter
    def type(self, type: str):
        """Sets the type of this APIBackendConfig.

        Type of the backend. Available values:    * local    * config  # noqa: E501

        :param type: The type of this APIBackendConfig.
        :type type: str
        """

        self._type = type