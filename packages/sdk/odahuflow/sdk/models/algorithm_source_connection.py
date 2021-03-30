# coding: utf-8

from __future__ import absolute_import
from datetime import date, datetime  # noqa: F401

from typing import List, Dict  # noqa: F401

from odahuflow.sdk.models.base_model_ import Model
from odahuflow.sdk.models.connection import Connection  # noqa: F401,E501
from odahuflow.sdk.models import util


class AlgorithmSourceConnection(Model):
    """NOTE: This class is auto generated by the swagger code generator program.

    Do not edit the class manually.
    """

    def __init__(self, conn: Connection=None, path: str=None):  # noqa: E501
        """AlgorithmSourceConnection - a model defined in Swagger

        :param conn: The conn of this AlgorithmSourceConnection.  # noqa: E501
        :type conn: Connection
        :param path: The path of this AlgorithmSourceConnection.  # noqa: E501
        :type path: str
        """
        self.swagger_types = {
            'conn': Connection,
            'path': str
        }

        self.attribute_map = {
            'conn': 'conn',
            'path': 'path'
        }

        self._conn = conn
        self._path = path

    @classmethod
    def from_dict(cls, dikt) -> 'AlgorithmSourceConnection':
        """Returns the dict as a model

        :param dikt: A dict.
        :type: dict
        :return: The AlgorithmSourceConnection of this AlgorithmSourceConnection.  # noqa: E501
        :rtype: AlgorithmSourceConnection
        """
        return util.deserialize_model(dikt, cls)

    @property
    def conn(self) -> Connection:
        """Gets the conn of this AlgorithmSourceConnection.

        Connection specific for Algorithm  # noqa: E501

        :return: The conn of this AlgorithmSourceConnection.
        :rtype: Connection
        """
        return self._conn

    @conn.setter
    def conn(self, conn: Connection):
        """Sets the conn of this AlgorithmSourceConnection.

        Connection specific for Algorithm  # noqa: E501

        :param conn: The conn of this AlgorithmSourceConnection.
        :type conn: Connection
        """

        self._conn = conn

    @property
    def path(self) -> str:
        """Gets the path of this AlgorithmSourceConnection.

        Remote path for object storage  # noqa: E501

        :return: The path of this AlgorithmSourceConnection.
        :rtype: str
        """
        return self._path

    @path.setter
    def path(self, path: str):
        """Sets the path of this AlgorithmSourceConnection.

        Remote path for object storage  # noqa: E501

        :param path: The path of this AlgorithmSourceConnection.
        :type path: str
        """

        self._path = path