# coding: utf-8

from __future__ import absolute_import
from datetime import date, datetime  # noqa: F401

from typing import List, Dict  # noqa: F401

from odahuflow.sdk.models.base_model_ import Model
from odahuflow.sdk.models.connection import Connection  # noqa: F401,E501
from odahuflow.sdk.models import util


class PackagerTarget(Model):
    """NOTE: This class is auto generated by the swagger code generator program.

    Do not edit the class manually.
    """

    def __init__(self, connection: Connection = None, name: str = None):  # noqa: E501
        """PackagerTarget - a model defined in Swagger

        :param connection: The connection of this PackagerTarget.  # noqa: E501
        :type connection: Connection
        :param name: The name of this PackagerTarget.  # noqa: E501
        :type name: str
        """
        self.swagger_types = {"connection": Connection, "name": str}

        self.attribute_map = {"connection": "connection", "name": "name"}

        self._connection = connection
        self._name = name

    @classmethod
    def from_dict(cls, dikt) -> "PackagerTarget":
        """Returns the dict as a model

        :param dikt: A dict.
        :type: dict
        :return: The PackagerTarget of this PackagerTarget.  # noqa: E501
        :rtype: PackagerTarget
        """
        return util.deserialize_model(dikt, cls)

    @property
    def connection(self) -> Connection:
        """Gets the connection of this PackagerTarget.

        A Connection for this target  # noqa: E501

        :return: The connection of this PackagerTarget.
        :rtype: Connection
        """
        return self._connection

    @connection.setter
    def connection(self, connection: Connection):
        """Sets the connection of this PackagerTarget.

        A Connection for this target  # noqa: E501

        :param connection: The connection of this PackagerTarget.
        :type connection: Connection
        """

        self._connection = connection

    @property
    def name(self) -> str:
        """Gets the name of this PackagerTarget.

        Target name  # noqa: E501

        :return: The name of this PackagerTarget.
        :rtype: str
        """
        return self._name

    @name.setter
    def name(self, name: str):
        """Sets the name of this PackagerTarget.

        Target name  # noqa: E501

        :param name: The name of this PackagerTarget.
        :type name: str
        """

        self._name = name
