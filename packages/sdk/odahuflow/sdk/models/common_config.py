# coding: utf-8

from __future__ import absolute_import
from datetime import date, datetime  # noqa: F401

from typing import List, Dict  # noqa: F401

from odahuflow.sdk.models.base_model_ import Model
from odahuflow.sdk.models.external_url import ExternalUrl  # noqa: F401,E501
from odahuflow.sdk.models import util


class CommonConfig(Model):
    """NOTE: This class is auto generated by the swagger code generator program.

    Do not edit the class manually.
    """

    def __init__(self, database_connection_string: str=None, external_urls: List[ExternalUrl]=None, graceful_timeout: str=None, launch_period: str=None, oauth_oidc_token_endpoint: str=None, resource_gpu_name: str=None, version: str=None):  # noqa: E501
        """CommonConfig - a model defined in Swagger

        :param database_connection_string: The database_connection_string of this CommonConfig.  # noqa: E501
        :type database_connection_string: str
        :param external_urls: The external_urls of this CommonConfig.  # noqa: E501
        :type external_urls: List[ExternalUrl]
        :param graceful_timeout: The graceful_timeout of this CommonConfig.  # noqa: E501
        :type graceful_timeout: str
        :param launch_period: The launch_period of this CommonConfig.  # noqa: E501
        :type launch_period: str
        :param oauth_oidc_token_endpoint: The oauth_oidc_token_endpoint of this CommonConfig.  # noqa: E501
        :type oauth_oidc_token_endpoint: str
        :param resource_gpu_name: The resource_gpu_name of this CommonConfig.  # noqa: E501
        :type resource_gpu_name: str
        :param version: The version of this CommonConfig.  # noqa: E501
        :type version: str
        """
        self.swagger_types = {
            'database_connection_string': str,
            'external_urls': List[ExternalUrl],
            'graceful_timeout': str,
            'launch_period': str,
            'oauth_oidc_token_endpoint': str,
            'resource_gpu_name': str,
            'version': str
        }

        self.attribute_map = {
            'database_connection_string': 'databaseConnectionString',
            'external_urls': 'externalUrls',
            'graceful_timeout': 'gracefulTimeout',
            'launch_period': 'launchPeriod',
            'oauth_oidc_token_endpoint': 'oauthOidcTokenEndpoint',
            'resource_gpu_name': 'resourceGpuName',
            'version': 'version'
        }

        self._database_connection_string = database_connection_string
        self._external_urls = external_urls
        self._graceful_timeout = graceful_timeout
        self._launch_period = launch_period
        self._oauth_oidc_token_endpoint = oauth_oidc_token_endpoint
        self._resource_gpu_name = resource_gpu_name
        self._version = version

    @classmethod
    def from_dict(cls, dikt) -> 'CommonConfig':
        """Returns the dict as a model

        :param dikt: A dict.
        :type: dict
        :return: The CommonConfig of this CommonConfig.  # noqa: E501
        :rtype: CommonConfig
        """
        return util.deserialize_model(dikt, cls)

    @property
    def database_connection_string(self) -> str:
        """Gets the database_connection_string of this CommonConfig.

        Database connection string  # noqa: E501

        :return: The database_connection_string of this CommonConfig.
        :rtype: str
        """
        return self._database_connection_string

    @database_connection_string.setter
    def database_connection_string(self, database_connection_string: str):
        """Sets the database_connection_string of this CommonConfig.

        Database connection string  # noqa: E501

        :param database_connection_string: The database_connection_string of this CommonConfig.
        :type database_connection_string: str
        """

        self._database_connection_string = database_connection_string

    @property
    def external_urls(self) -> List[ExternalUrl]:
        """Gets the external_urls of this CommonConfig.

        The collection of external urls, for example: metrics, edge, service catalog and so on  # noqa: E501

        :return: The external_urls of this CommonConfig.
        :rtype: List[ExternalUrl]
        """
        return self._external_urls

    @external_urls.setter
    def external_urls(self, external_urls: List[ExternalUrl]):
        """Sets the external_urls of this CommonConfig.

        The collection of external urls, for example: metrics, edge, service catalog and so on  # noqa: E501

        :param external_urls: The external_urls of this CommonConfig.
        :type external_urls: List[ExternalUrl]
        """

        self._external_urls = external_urls

    @property
    def graceful_timeout(self) -> str:
        """Gets the graceful_timeout of this CommonConfig.

        Graceful shutdown timeout  # noqa: E501

        :return: The graceful_timeout of this CommonConfig.
        :rtype: str
        """
        return self._graceful_timeout

    @graceful_timeout.setter
    def graceful_timeout(self, graceful_timeout: str):
        """Sets the graceful_timeout of this CommonConfig.

        Graceful shutdown timeout  # noqa: E501

        :param graceful_timeout: The graceful_timeout of this CommonConfig.
        :type graceful_timeout: str
        """

        self._graceful_timeout = graceful_timeout

    @property
    def launch_period(self) -> str:
        """Gets the launch_period of this CommonConfig.

        How often launch new training  # noqa: E501

        :return: The launch_period of this CommonConfig.
        :rtype: str
        """
        return self._launch_period

    @launch_period.setter
    def launch_period(self, launch_period: str):
        """Sets the launch_period of this CommonConfig.

        How often launch new training  # noqa: E501

        :param launch_period: The launch_period of this CommonConfig.
        :type launch_period: str
        """

        self._launch_period = launch_period

    @property
    def oauth_oidc_token_endpoint(self) -> str:
        """Gets the oauth_oidc_token_endpoint of this CommonConfig.

        OpenID token url  # noqa: E501

        :return: The oauth_oidc_token_endpoint of this CommonConfig.
        :rtype: str
        """
        return self._oauth_oidc_token_endpoint

    @oauth_oidc_token_endpoint.setter
    def oauth_oidc_token_endpoint(self, oauth_oidc_token_endpoint: str):
        """Sets the oauth_oidc_token_endpoint of this CommonConfig.

        OpenID token url  # noqa: E501

        :param oauth_oidc_token_endpoint: The oauth_oidc_token_endpoint of this CommonConfig.
        :type oauth_oidc_token_endpoint: str
        """

        self._oauth_oidc_token_endpoint = oauth_oidc_token_endpoint

    @property
    def resource_gpu_name(self) -> str:
        """Gets the resource_gpu_name of this CommonConfig.

        Kubernetes can consume the GPU resource in the <vendor>.com/gpu format. For example, amd.com/gpu or nvidia.com/gpu.  # noqa: E501

        :return: The resource_gpu_name of this CommonConfig.
        :rtype: str
        """
        return self._resource_gpu_name

    @resource_gpu_name.setter
    def resource_gpu_name(self, resource_gpu_name: str):
        """Sets the resource_gpu_name of this CommonConfig.

        Kubernetes can consume the GPU resource in the <vendor>.com/gpu format. For example, amd.com/gpu or nvidia.com/gpu.  # noqa: E501

        :param resource_gpu_name: The resource_gpu_name of this CommonConfig.
        :type resource_gpu_name: str
        """

        self._resource_gpu_name = resource_gpu_name

    @property
    def version(self) -> str:
        """Gets the version of this CommonConfig.

        Version of ODAHU platform  # noqa: E501

        :return: The version of this CommonConfig.
        :rtype: str
        """
        return self._version

    @version.setter
    def version(self, version: str):
        """Sets the version of this CommonConfig.

        Version of ODAHU platform  # noqa: E501

        :param version: The version of this CommonConfig.
        :type version: str
        """

        self._version = version
