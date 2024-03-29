# coding: utf-8

from __future__ import absolute_import
from datetime import date, datetime  # noqa: F401

from typing import List, Dict  # noqa: F401

from odahuflow.sdk.models.base_model_ import Model
from odahuflow.sdk.models.local_model_source import LocalModelSource  # noqa: F401,E501
from odahuflow.sdk.models.remote_model_source import RemoteModelSource  # noqa: F401,E501
from odahuflow.sdk.models import util


class ModelSource(Model):
    """NOTE: This class is auto generated by the swagger code generator program.

    Do not edit the class manually.
    """

    def __init__(self, local: LocalModelSource=None, remote: RemoteModelSource=None):  # noqa: E501
        """ModelSource - a model defined in Swagger

        :param local: The local of this ModelSource.  # noqa: E501
        :type local: LocalModelSource
        :param remote: The remote of this ModelSource.  # noqa: E501
        :type remote: RemoteModelSource
        """
        self.swagger_types = {
            'local': LocalModelSource,
            'remote': RemoteModelSource
        }

        self.attribute_map = {
            'local': 'local',
            'remote': 'remote'
        }

        self._local = local
        self._remote = remote

    @classmethod
    def from_dict(cls, dikt) -> 'ModelSource':
        """Returns the dict as a model

        :param dikt: A dict.
        :type: dict
        :return: The ModelSource of this ModelSource.  # noqa: E501
        :rtype: ModelSource
        """
        return util.deserialize_model(dikt, cls)

    @property
    def local(self) -> LocalModelSource:
        """Gets the local of this ModelSource.

        Local does not fetch model and assume that model is embedded into container  # noqa: E501

        :return: The local of this ModelSource.
        :rtype: LocalModelSource
        """
        return self._local

    @local.setter
    def local(self, local: LocalModelSource):
        """Sets the local of this ModelSource.

        Local does not fetch model and assume that model is embedded into container  # noqa: E501

        :param local: The local of this ModelSource.
        :type local: LocalModelSource
        """

        self._local = local

    @property
    def remote(self) -> RemoteModelSource:
        """Gets the remote of this ModelSource.

        Remote fetch model from remote model registry using ODAHU connections mechanism  # noqa: E501

        :return: The remote of this ModelSource.
        :rtype: RemoteModelSource
        """
        return self._remote

    @remote.setter
    def remote(self, remote: RemoteModelSource):
        """Sets the remote of this ModelSource.

        Remote fetch model from remote model registry using ODAHU connections mechanism  # noqa: E501

        :param remote: The remote of this ModelSource.
        :type remote: RemoteModelSource
        """

        self._remote = remote
