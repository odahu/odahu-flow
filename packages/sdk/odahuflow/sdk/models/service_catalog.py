# coding: utf-8

from __future__ import absolute_import
from datetime import date, datetime  # noqa: F401

from typing import List, Dict  # noqa: F401

from odahuflow.sdk.models.base_model_ import Model
from odahuflow.sdk.models.auth_config import AuthConfig  # noqa: F401,E501
from odahuflow.sdk.models import util


class ServiceCatalog(Model):
    """NOTE: This class is auto generated by the swagger code generator program.

    Do not edit the class manually.
    """

    def __init__(self, auth: AuthConfig=None, base_url: str=None, debug: bool=None, edge_host: str=None, edge_url: str=None, fetch_timeout: int=None, workers_count: int=None):  # noqa: E501
        """ServiceCatalog - a model defined in Swagger

        :param auth: The auth of this ServiceCatalog.  # noqa: E501
        :type auth: AuthConfig
        :param base_url: The base_url of this ServiceCatalog.  # noqa: E501
        :type base_url: str
        :param debug: The debug of this ServiceCatalog.  # noqa: E501
        :type debug: bool
        :param edge_host: The edge_host of this ServiceCatalog.  # noqa: E501
        :type edge_host: str
        :param edge_url: The edge_url of this ServiceCatalog.  # noqa: E501
        :type edge_url: str
        :param fetch_timeout: The fetch_timeout of this ServiceCatalog.  # noqa: E501
        :type fetch_timeout: int
        :param workers_count: The workers_count of this ServiceCatalog.  # noqa: E501
        :type workers_count: int
        """
        self.swagger_types = {
            'auth': AuthConfig,
            'base_url': str,
            'debug': bool,
            'edge_host': str,
            'edge_url': str,
            'fetch_timeout': int,
            'workers_count': int
        }

        self.attribute_map = {
            'auth': 'auth',
            'base_url': 'baseUrl',
            'debug': 'debug',
            'edge_host': 'edgeHost',
            'edge_url': 'edgeURL',
            'fetch_timeout': 'fetchTimeout',
            'workers_count': 'workersCount'
        }

        self._auth = auth
        self._base_url = base_url
        self._debug = debug
        self._edge_host = edge_host
        self._edge_url = edge_url
        self._fetch_timeout = fetch_timeout
        self._workers_count = workers_count

    @classmethod
    def from_dict(cls, dikt) -> 'ServiceCatalog':
        """Returns the dict as a model

        :param dikt: A dict.
        :type: dict
        :return: The ServiceCatalog of this ServiceCatalog.  # noqa: E501
        :rtype: ServiceCatalog
        """
        return util.deserialize_model(dikt, cls)

    @property
    def auth(self) -> AuthConfig:
        """Gets the auth of this ServiceCatalog.

        Auth configures connection parameters to ODAHU API Server  # noqa: E501

        :return: The auth of this ServiceCatalog.
        :rtype: AuthConfig
        """
        return self._auth

    @auth.setter
    def auth(self, auth: AuthConfig):
        """Sets the auth of this ServiceCatalog.

        Auth configures connection parameters to ODAHU API Server  # noqa: E501

        :param auth: The auth of this ServiceCatalog.
        :type auth: AuthConfig
        """

        self._auth = auth

    @property
    def base_url(self) -> str:
        """Gets the base_url of this ServiceCatalog.

        BaseURL is a prefix to service catalog web server endpoints  # noqa: E501

        :return: The base_url of this ServiceCatalog.
        :rtype: str
        """
        return self._base_url

    @base_url.setter
    def base_url(self, base_url: str):
        """Sets the base_url of this ServiceCatalog.

        BaseURL is a prefix to service catalog web server endpoints  # noqa: E501

        :param base_url: The base_url of this ServiceCatalog.
        :type base_url: str
        """

        self._base_url = base_url

    @property
    def debug(self) -> bool:
        """Gets the debug of this ServiceCatalog.

        enabled Debug increase logger verbosity and format. Default: false  # noqa: E501

        :return: The debug of this ServiceCatalog.
        :rtype: bool
        """
        return self._debug

    @debug.setter
    def debug(self, debug: bool):
        """Sets the debug of this ServiceCatalog.

        enabled Debug increase logger verbosity and format. Default: false  # noqa: E501

        :param debug: The debug of this ServiceCatalog.
        :type debug: bool
        """

        self._debug = debug

    @property
    def edge_host(self) -> str:
        """Gets the edge_host of this ServiceCatalog.

        ServiceCatalog set EdgeHost as Host header in requests to ML servers  # noqa: E501

        :return: The edge_host of this ServiceCatalog.
        :rtype: str
        """
        return self._edge_host

    @edge_host.setter
    def edge_host(self, edge_host: str):
        """Sets the edge_host of this ServiceCatalog.

        ServiceCatalog set EdgeHost as Host header in requests to ML servers  # noqa: E501

        :param edge_host: The edge_host of this ServiceCatalog.
        :type edge_host: str
        """

        self._edge_host = edge_host

    @property
    def edge_url(self) -> str:
        """Gets the edge_url of this ServiceCatalog.

        ServiceCatalog uses EdgeURL to call MLServer by adding ModelRoute prefix to EdgeURL path  # noqa: E501

        :return: The edge_url of this ServiceCatalog.
        :rtype: str
        """
        return self._edge_url

    @edge_url.setter
    def edge_url(self, edge_url: str):
        """Sets the edge_url of this ServiceCatalog.

        ServiceCatalog uses EdgeURL to call MLServer by adding ModelRoute prefix to EdgeURL path  # noqa: E501

        :param edge_url: The edge_url of this ServiceCatalog.
        :type edge_url: str
        """

        self._edge_url = edge_url

    @property
    def fetch_timeout(self) -> int:
        """Gets the fetch_timeout of this ServiceCatalog.

        FetchTimeout configures how often new events will be fetched. Default 5 seconds.  # noqa: E501

        :return: The fetch_timeout of this ServiceCatalog.
        :rtype: int
        """
        return self._fetch_timeout

    @fetch_timeout.setter
    def fetch_timeout(self, fetch_timeout: int):
        """Sets the fetch_timeout of this ServiceCatalog.

        FetchTimeout configures how often new events will be fetched. Default 5 seconds.  # noqa: E501

        :param fetch_timeout: The fetch_timeout of this ServiceCatalog.
        :type fetch_timeout: int
        """

        self._fetch_timeout = fetch_timeout

    @property
    def workers_count(self) -> int:
        """Gets the workers_count of this ServiceCatalog.

        WorkersCount configures how many workers will process events. Default: 4  # noqa: E501

        :return: The workers_count of this ServiceCatalog.
        :rtype: int
        """
        return self._workers_count

    @workers_count.setter
    def workers_count(self, workers_count: int):
        """Sets the workers_count of this ServiceCatalog.

        WorkersCount configures how many workers will process events. Default: 4  # noqa: E501

        :param workers_count: The workers_count of this ServiceCatalog.
        :type workers_count: int
        """

        self._workers_count = workers_count
