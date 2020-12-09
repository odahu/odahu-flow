# coding: utf-8

from __future__ import absolute_import
from datetime import date, datetime  # noqa: F401

from typing import List, Dict  # noqa: F401

from odahuflow.sdk.models.base_model_ import Model
from odahuflow.sdk.models.api_backend_config import APIBackendConfig  # noqa: F401,E501
from odahuflow.sdk.models import util


class APIConfig(Model):
    """NOTE: This class is auto generated by the swagger code generator program.

    Do not edit the class manually.
    """

    def __init__(
        self, backend: APIBackendConfig = None, port: int = None
    ):  # noqa: E501
        """APIConfig - a model defined in Swagger

        :param backend: The backend of this APIConfig.  # noqa: E501
        :type backend: APIBackendConfig
        :param port: The port of this APIConfig.  # noqa: E501
        :type port: int
        """
        self.swagger_types = {"backend": APIBackendConfig, "port": int}

        self.attribute_map = {"backend": "backend", "port": "port"}

        self._backend = backend
        self._port = port

    @classmethod
    def from_dict(cls, dikt) -> "APIConfig":
        """Returns the dict as a model

        :param dikt: A dict.
        :type: dict
        :return: The APIConfig of this APIConfig.  # noqa: E501
        :rtype: APIConfig
        """
        return util.deserialize_model(dikt, cls)

    @property
    def backend(self) -> APIBackendConfig:
        """Gets the backend of this APIConfig.


        :return: The backend of this APIConfig.
        :rtype: APIBackendConfig
        """
        return self._backend

    @backend.setter
    def backend(self, backend: APIBackendConfig):
        """Sets the backend of this APIConfig.


        :param backend: The backend of this APIConfig.
        :type backend: APIBackendConfig
        """

        self._backend = backend

    @property
    def port(self) -> int:
        """Gets the port of this APIConfig.

        API HTTP port  # noqa: E501

        :return: The port of this APIConfig.
        :rtype: int
        """
        return self._port

    @port.setter
    def port(self, port: int):
        """Sets the port of this APIConfig.

        API HTTP port  # noqa: E501

        :param port: The port of this APIConfig.
        :type port: int
        """

        self._port = port
