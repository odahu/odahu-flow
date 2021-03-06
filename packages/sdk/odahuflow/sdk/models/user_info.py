# coding: utf-8

from __future__ import absolute_import
from datetime import date, datetime  # noqa: F401

from typing import List, Dict  # noqa: F401

from odahuflow.sdk.models.base_model_ import Model
from odahuflow.sdk.models import util


class UserInfo(Model):
    """NOTE: This class is auto generated by the swagger code generator program.

    Do not edit the class manually.
    """

    def __init__(self, email: str=None, username: str=None):  # noqa: E501
        """UserInfo - a model defined in Swagger

        :param email: The email of this UserInfo.  # noqa: E501
        :type email: str
        :param username: The username of this UserInfo.  # noqa: E501
        :type username: str
        """
        self.swagger_types = {
            'email': str,
            'username': str
        }

        self.attribute_map = {
            'email': 'email',
            'username': 'username'
        }

        self._email = email
        self._username = username

    @classmethod
    def from_dict(cls, dikt) -> 'UserInfo':
        """Returns the dict as a model

        :param dikt: A dict.
        :type: dict
        :return: The UserInfo of this UserInfo.  # noqa: E501
        :rtype: UserInfo
        """
        return util.deserialize_model(dikt, cls)

    @property
    def email(self) -> str:
        """Gets the email of this UserInfo.


        :return: The email of this UserInfo.
        :rtype: str
        """
        return self._email

    @email.setter
    def email(self, email: str):
        """Sets the email of this UserInfo.


        :param email: The email of this UserInfo.
        :type email: str
        """

        self._email = email

    @property
    def username(self) -> str:
        """Gets the username of this UserInfo.


        :return: The username of this UserInfo.
        :rtype: str
        """
        return self._username

    @username.setter
    def username(self, username: str):
        """Sets the username of this UserInfo.


        :param username: The username of this UserInfo.
        :type username: str
        """

        self._username = username
