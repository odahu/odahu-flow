# coding: utf-8

from __future__ import absolute_import
from datetime import date, datetime  # noqa: F401

from typing import List, Dict  # noqa: F401

from odahuflow.sdk.models.base_model_ import Model
from odahuflow.sdk.models.inference_service_spec import InferenceServiceSpec  # noqa: F401,E501
from odahuflow.sdk.models.inference_service_status import InferenceServiceStatus  # noqa: F401,E501
from odahuflow.sdk.models import util


class InferenceService(Model):
    """NOTE: This class is auto generated by the swagger code generator program.

    Do not edit the class manually.
    """

    def __init__(self, created_at: str=None, id: str=None, spec: InferenceServiceSpec=None, status: InferenceServiceStatus=None, updated_at: str=None):  # noqa: E501
        """InferenceService - a model defined in Swagger

        :param created_at: The created_at of this InferenceService.  # noqa: E501
        :type created_at: str
        :param id: The id of this InferenceService.  # noqa: E501
        :type id: str
        :param spec: The spec of this InferenceService.  # noqa: E501
        :type spec: InferenceServiceSpec
        :param status: The status of this InferenceService.  # noqa: E501
        :type status: InferenceServiceStatus
        :param updated_at: The updated_at of this InferenceService.  # noqa: E501
        :type updated_at: str
        """
        self.swagger_types = {
            'created_at': str,
            'id': str,
            'spec': InferenceServiceSpec,
            'status': InferenceServiceStatus,
            'updated_at': str
        }

        self.attribute_map = {
            'created_at': 'createdAt',
            'id': 'id',
            'spec': 'spec',
            'status': 'status',
            'updated_at': 'updatedAt'
        }

        self._created_at = created_at
        self._id = id
        self._spec = spec
        self._status = status
        self._updated_at = updated_at

    @classmethod
    def from_dict(cls, dikt) -> 'InferenceService':
        """Returns the dict as a model

        :param dikt: A dict.
        :type: dict
        :return: The InferenceService of this InferenceService.  # noqa: E501
        :rtype: InferenceService
        """
        return util.deserialize_model(dikt, cls)

    @property
    def created_at(self) -> str:
        """Gets the created_at of this InferenceService.

        When resource was created. Managed by system. Cannot be overridden by User  # noqa: E501

        :return: The created_at of this InferenceService.
        :rtype: str
        """
        return self._created_at

    @created_at.setter
    def created_at(self, created_at: str):
        """Sets the created_at of this InferenceService.

        When resource was created. Managed by system. Cannot be overridden by User  # noqa: E501

        :param created_at: The created_at of this InferenceService.
        :type created_at: str
        """

        self._created_at = created_at

    @property
    def id(self) -> str:
        """Gets the id of this InferenceService.


        :return: The id of this InferenceService.
        :rtype: str
        """
        return self._id

    @id.setter
    def id(self, id: str):
        """Sets the id of this InferenceService.


        :param id: The id of this InferenceService.
        :type id: str
        """

        self._id = id

    @property
    def spec(self) -> InferenceServiceSpec:
        """Gets the spec of this InferenceService.


        :return: The spec of this InferenceService.
        :rtype: InferenceServiceSpec
        """
        return self._spec

    @spec.setter
    def spec(self, spec: InferenceServiceSpec):
        """Sets the spec of this InferenceService.


        :param spec: The spec of this InferenceService.
        :type spec: InferenceServiceSpec
        """

        self._spec = spec

    @property
    def status(self) -> InferenceServiceStatus:
        """Gets the status of this InferenceService.


        :return: The status of this InferenceService.
        :rtype: InferenceServiceStatus
        """
        return self._status

    @status.setter
    def status(self, status: InferenceServiceStatus):
        """Sets the status of this InferenceService.


        :param status: The status of this InferenceService.
        :type status: InferenceServiceStatus
        """

        self._status = status

    @property
    def updated_at(self) -> str:
        """Gets the updated_at of this InferenceService.

        When resource was updated. Managed by system. Cannot be overridden by User  # noqa: E501

        :return: The updated_at of this InferenceService.
        :rtype: str
        """
        return self._updated_at

    @updated_at.setter
    def updated_at(self, updated_at: str):
        """Sets the updated_at of this InferenceService.

        When resource was updated. Managed by system. Cannot be overridden by User  # noqa: E501

        :param updated_at: The updated_at of this InferenceService.
        :type updated_at: str
        """

        self._updated_at = updated_at
