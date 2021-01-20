# coding: utf-8

from __future__ import absolute_import
from datetime import date, datetime  # noqa: F401

from typing import List, Dict  # noqa: F401

from odahuflow.sdk.models.base_model_ import Model
from odahuflow.sdk.models import util


class JWKS(Model):
    """NOTE: This class is auto generated by the swagger code generator program.

    Do not edit the class manually.
    """

    def __init__(
        self, enabled: bool = None, issuer: str = None, url: str = None
    ):  # noqa: E501
        """JWKS - a model defined in Swagger

        :param enabled: The enabled of this JWKS.  # noqa: E501
        :type enabled: bool
        :param issuer: The issuer of this JWKS.  # noqa: E501
        :type issuer: str
        :param url: The url of this JWKS.  # noqa: E501
        :type url: str
        """
        self.swagger_types = {"enabled": bool, "issuer": str, "url": str}

        self.attribute_map = {"enabled": "enabled", "issuer": "issuer", "url": "url"}

        self._enabled = enabled
        self._issuer = issuer
        self._url = url

    @classmethod
    def from_dict(cls, dikt) -> "JWKS":
        """Returns the dict as a model

        :param dikt: A dict.
        :type: dict
        :return: The JWKS of this JWKS.  # noqa: E501
        :rtype: JWKS
        """
        return util.deserialize_model(dikt, cls)

    @property
    def enabled(self) -> bool:
        """Gets the enabled of this JWKS.

        Model authorization enabled  # noqa: E501

        :return: The enabled of this JWKS.
        :rtype: bool
        """
        return self._enabled

    @enabled.setter
    def enabled(self, enabled: bool):
        """Sets the enabled of this JWKS.

        Model authorization enabled  # noqa: E501

        :param enabled: The enabled of this JWKS.
        :type enabled: bool
        """

        self._enabled = enabled

    @property
    def issuer(self) -> str:
        """Gets the issuer of this JWKS.

        Issuer claim value  # noqa: E501

        :return: The issuer of this JWKS.
        :rtype: str
        """
        return self._issuer

    @issuer.setter
    def issuer(self, issuer: str):
        """Sets the issuer of this JWKS.

        Issuer claim value  # noqa: E501

        :param issuer: The issuer of this JWKS.
        :type issuer: str
        """

        self._issuer = issuer

    @property
    def url(self) -> str:
        """Gets the url of this JWKS.

        JWKS URL  # noqa: E501

        :return: The url of this JWKS.
        :rtype: str
        """
        return self._url

    @url.setter
    def url(self, url: str):
        """Sets the url of this JWKS.

        JWKS URL  # noqa: E501

        :param url: The url of this JWKS.
        :type url: str
        """

        self._url = url
