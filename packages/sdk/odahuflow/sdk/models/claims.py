# coding: utf-8

from __future__ import absolute_import
from datetime import date, datetime  # noqa: F401

from typing import List, Dict  # noqa: F401

from odahuflow.sdk.models.base_model_ import Model
from odahuflow.sdk.models import util


class Claims(Model):
    """NOTE: This class is auto generated by the swagger code generator program.

    Do not edit the class manually.
    """

    def __init__(self, email: str = None, name: str = None):  # noqa: E501
        """Claims - a model defined in Swagger

        :param email: The email of this Claims.  # noqa: E501
        :type email: str
        :param name: The name of this Claims.  # noqa: E501
        :type name: str
        """
        self.swagger_types = {"email": str, "name": str}

        self.attribute_map = {"email": "email", "name": "name"}

        self._email = email
        self._name = name

    @classmethod
    def from_dict(cls, dikt) -> "Claims":
        """Returns the dict as a model

        :param dikt: A dict.
        :type: dict
        :return: The Claims of this Claims.  # noqa: E501
        :rtype: Claims
        """
        return util.deserialize_model(dikt, cls)

    @property
    def email(self) -> str:
        """Gets the email of this Claims.


        :return: The email of this Claims.
        :rtype: str
        """
        return self._email

    @email.setter
    def email(self, email: str):
        """Sets the email of this Claims.


        :param email: The email of this Claims.
        :type email: str
        """

        self._email = email

    @property
    def name(self) -> str:
        """Gets the name of this Claims.


        :return: The name of this Claims.
        :rtype: str
        """
        return self._name

    @name.setter
    def name(self, name: str):
        """Sets the name of this Claims.


        :param name: The name of this Claims.
        :type name: str
        """

        self._name = name
