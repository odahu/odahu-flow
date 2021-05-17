# coding: utf-8

from __future__ import absolute_import
from datetime import date, datetime  # noqa: F401

from typing import List, Dict  # noqa: F401

from odahuflow.sdk.models.base_model_ import Model
from odahuflow.sdk.models import util


class ConnectionStatus(Model):
    """NOTE: This class is auto generated by the swagger code generator program.

    Do not edit the class manually.
    """

    def __init__(self, secret_name: str=None, service_account: str=None):  # noqa: E501
        """ConnectionStatus - a model defined in Swagger

        :param secret_name: The secret_name of this ConnectionStatus.  # noqa: E501
        :type secret_name: str
        :param service_account: The service_account of this ConnectionStatus.  # noqa: E501
        :type service_account: str
        """
        self.swagger_types = {
            'secret_name': str,
            'service_account': str
        }

        self.attribute_map = {
            'secret_name': 'secretName',
            'service_account': 'serviceAccount'
        }

        self._secret_name = secret_name
        self._service_account = service_account

    @classmethod
    def from_dict(cls, dikt) -> 'ConnectionStatus':
        """Returns the dict as a model

        :param dikt: A dict.
        :type: dict
        :return: The ConnectionStatus of this ConnectionStatus.  # noqa: E501
        :rtype: ConnectionStatus
        """
        return util.deserialize_model(dikt, cls)

    @property
    def secret_name(self) -> str:
        """Gets the secret_name of this ConnectionStatus.

        Kubernetes secret name  # noqa: E501

        :return: The secret_name of this ConnectionStatus.
        :rtype: str
        """
        return self._secret_name

    @secret_name.setter
    def secret_name(self, secret_name: str):
        """Sets the secret_name of this ConnectionStatus.

        Kubernetes secret name  # noqa: E501

        :param secret_name: The secret_name of this ConnectionStatus.
        :type secret_name: str
        """

        self._secret_name = secret_name

    @property
    def service_account(self) -> str:
        """Gets the service_account of this ConnectionStatus.

        Kubernetes service account  # noqa: E501

        :return: The service_account of this ConnectionStatus.
        :rtype: str
        """
        return self._service_account

    @service_account.setter
    def service_account(self, service_account: str):
        """Sets the service_account of this ConnectionStatus.

        Kubernetes service account  # noqa: E501

        :param service_account: The service_account of this ConnectionStatus.
        :type service_account: str
        """

        self._service_account = service_account
