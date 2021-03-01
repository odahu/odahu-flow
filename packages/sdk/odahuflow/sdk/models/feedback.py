from typing import Dict, Any, Type, TypeVar

from odahuflow.sdk.models.base_model_ import Model
from odahuflow.sdk.models import util

T = TypeVar('T')


class FeedbackModel(Model):
    """NOTE: This class is created manually."""

    def __init__(
            self,
            feedback: Dict[Any, Any] = None,
            request_id: str = None,
            model_name: str = None,
            model_version: str = None
    ):
        """
        :param feedback: model feedback
        :type feedback: dict
        :param request_id: feedback request id
        :type request_id: str
        :param model_name: model name for feedback
        :type model_name: str
        :param model_version: model version for feedback
        :type model_version: str
        """

        self._feedback = feedback
        self._request_id = request_id
        self._model_name = model_name
        self._model_version = model_version

    @classmethod
    def from_dict(cls: Type[T], dikt) -> 'FeedbackModel':
        """Returns the dict as a model

        :param dikt: feedback response data
        :type dikt : dict
        :return: FeedbackModel instance
        :rtype: FeedbackModel
        """
        data = dikt.get('message')
        return FeedbackModel(
            feedback=data.get('Payload', {}).get('json'),
            request_id=data.get('RequestID '),
            model_name=data.get('ModelName'),
            model_version=data.get('ModelVersion')
        )

    def to_dict(self):
        """Returns the model as a dict

        :return: model data in dict
        :rtype: dict
        """
        return {
            'RequestID': self._request_id,
            'ModelName': self._model_name,
            'ModelVersion': self._model_version,
            'Feedback': self._feedback
        }

    @property
    def feedback(self) -> dict:
        """Gets the feedback of this FeedbackModel.

        :return: The feedback
        :rtype: dict
        """
        return self._feedback

    @property
    def request_id(self) -> str:
        """Gets the request ID of this FeedbackModel.

        :return: Requet ID
        :rtype: str
        """
        return self._request_id

    @property
    def model_version(self) -> str:
        """Gets the model version of this FeedbackModel.

        :return: Model version
        :rtype: str
        """
        return self._model_version

    @property
    def model_name(self) -> str:
        """Gets the model name of this FeedbackModel.

        :return: Model name
        :rtype: str
        """
        return self._model_name
