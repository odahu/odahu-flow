# coding: utf-8

from __future__ import absolute_import
from datetime import date, datetime  # noqa: F401

from typing import List, Dict  # noqa: F401

from odahuflow.sdk.models.base_model_ import Model
from odahuflow.sdk.models import util


class InferenceJobStatus(Model):
    """NOTE: This class is auto generated by the swagger code generator program.

    Do not edit the class manually.
    """

    def __init__(self, message: str=None, pod_name: str=None, reason: str=None, state: str=None):  # noqa: E501
        """InferenceJobStatus - a model defined in Swagger

        :param message: The message of this InferenceJobStatus.  # noqa: E501
        :type message: str
        :param pod_name: The pod_name of this InferenceJobStatus.  # noqa: E501
        :type pod_name: str
        :param reason: The reason of this InferenceJobStatus.  # noqa: E501
        :type reason: str
        :param state: The state of this InferenceJobStatus.  # noqa: E501
        :type state: str
        """
        self.swagger_types = {
            'message': str,
            'pod_name': str,
            'reason': str,
            'state': str
        }

        self.attribute_map = {
            'message': 'message',
            'pod_name': 'podName',
            'reason': 'reason',
            'state': 'state'
        }

        self._message = message
        self._pod_name = pod_name
        self._reason = reason
        self._state = state

    @classmethod
    def from_dict(cls, dikt) -> 'InferenceJobStatus':
        """Returns the dict as a model

        :param dikt: A dict.
        :type: dict
        :return: The InferenceJobStatus of this InferenceJobStatus.  # noqa: E501
        :rtype: InferenceJobStatus
        """
        return util.deserialize_model(dikt, cls)

    @property
    def message(self) -> str:
        """Gets the message of this InferenceJobStatus.

        Message is any message from runtime service about status of InferenceJob  # noqa: E501

        :return: The message of this InferenceJobStatus.
        :rtype: str
        """
        return self._message

    @message.setter
    def message(self, message: str):
        """Sets the message of this InferenceJobStatus.

        Message is any message from runtime service about status of InferenceJob  # noqa: E501

        :param message: The message of this InferenceJobStatus.
        :type message: str
        """

        self._message = message

    @property
    def pod_name(self) -> str:
        """Gets the pod_name of this InferenceJobStatus.

        PodName is a name of Pod in Kubernetes that is running under the hood of InferenceJob  # noqa: E501

        :return: The pod_name of this InferenceJobStatus.
        :rtype: str
        """
        return self._pod_name

    @pod_name.setter
    def pod_name(self, pod_name: str):
        """Sets the pod_name of this InferenceJobStatus.

        PodName is a name of Pod in Kubernetes that is running under the hood of InferenceJob  # noqa: E501

        :param pod_name: The pod_name of this InferenceJobStatus.
        :type pod_name: str
        """

        self._pod_name = pod_name

    @property
    def reason(self) -> str:
        """Gets the reason of this InferenceJobStatus.

        Reason is a reason of some InferenceJob state that was retrieved from runtime service. for example reason of failure  # noqa: E501

        :return: The reason of this InferenceJobStatus.
        :rtype: str
        """
        return self._reason

    @reason.setter
    def reason(self, reason: str):
        """Sets the reason of this InferenceJobStatus.

        Reason is a reason of some InferenceJob state that was retrieved from runtime service. for example reason of failure  # noqa: E501

        :param reason: The reason of this InferenceJobStatus.
        :type reason: str
        """

        self._reason = reason

    @property
    def state(self) -> str:
        """Gets the state of this InferenceJobStatus.

        State describes current state of InferenceJob  # noqa: E501

        :return: The state of this InferenceJobStatus.
        :rtype: str
        """
        return self._state

    @state.setter
    def state(self, state: str):
        """Sets the state of this InferenceJobStatus.

        State describes current state of InferenceJob  # noqa: E501

        :param state: The state of this InferenceJobStatus.
        :type state: str
        """

        self._state = state
