# coding: utf-8

from __future__ import absolute_import
from datetime import date, datetime  # noqa: F401

from typing import List, Dict  # noqa: F401

from odahuflow.sdk.models.base_model_ import Model
from odahuflow.sdk.models import util


class DataBindingDir(Model):
    """NOTE: This class is auto generated by the swagger code generator program.

    Do not edit the class manually.
    """

    def __init__(self, connection: str=None, local_path: str=None, remote_path: str=None):  # noqa: E501
        """DataBindingDir - a model defined in Swagger

        :param connection: The connection of this DataBindingDir.  # noqa: E501
        :type connection: str
        :param local_path: The local_path of this DataBindingDir.  # noqa: E501
        :type local_path: str
        :param remote_path: The remote_path of this DataBindingDir.  # noqa: E501
        :type remote_path: str
        """
        self.swagger_types = {
            'connection': str,
            'local_path': str,
            'remote_path': str
        }

        self.attribute_map = {
            'connection': 'connection',
            'local_path': 'localPath',
            'remote_path': 'remotePath'
        }

        self._connection = connection
        self._local_path = local_path
        self._remote_path = remote_path

    @classmethod
    def from_dict(cls, dikt) -> 'DataBindingDir':
        """Returns the dict as a model

        :param dikt: A dict.
        :type: dict
        :return: The DataBindingDir of this DataBindingDir.  # noqa: E501
        :rtype: DataBindingDir
        """
        return util.deserialize_model(dikt, cls)

    @property
    def connection(self) -> str:
        """Gets the connection of this DataBindingDir.

        Connection name for data  # noqa: E501

        :return: The connection of this DataBindingDir.
        :rtype: str
        """
        return self._connection

    @connection.setter
    def connection(self, connection: str):
        """Sets the connection of this DataBindingDir.

        Connection name for data  # noqa: E501

        :param connection: The connection of this DataBindingDir.
        :type connection: str
        """

        self._connection = connection

    @property
    def local_path(self) -> str:
        """Gets the local_path of this DataBindingDir.

        Local training path  # noqa: E501

        :return: The local_path of this DataBindingDir.
        :rtype: str
        """
        return self._local_path

    @local_path.setter
    def local_path(self, local_path: str):
        """Sets the local_path of this DataBindingDir.

        Local training path  # noqa: E501

        :param local_path: The local_path of this DataBindingDir.
        :type local_path: str
        """

        self._local_path = local_path

    @property
    def remote_path(self) -> str:
        """Gets the remote_path of this DataBindingDir.

        Overwrite remote data path in connection  # noqa: E501

        :return: The remote_path of this DataBindingDir.
        :rtype: str
        """
        return self._remote_path

    @remote_path.setter
    def remote_path(self, remote_path: str):
        """Sets the remote_path of this DataBindingDir.

        Overwrite remote data path in connection  # noqa: E501

        :param remote_path: The remote_path of this DataBindingDir.
        :type remote_path: str
        """

        self._remote_path = remote_path
