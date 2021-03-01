# coding: utf-8

from __future__ import absolute_import
from datetime import date, datetime  # noqa: F401

from typing import List, Dict  # noqa: F401

from odahuflow.sdk.models.base_model_ import Model
from odahuflow.sdk.models import util


class FeedbackModelFeedbackRequest(Model):
    """NOTE: This class is manually edited.

    Requires replacement with auto generated model.
    """

    def __init__(self, feedback: Dict = None):  # noqa: E501
        """FeedbackModelFeedbackRequest - a model defined in Swagger

        """
        self.swagger_types = {
            'feedback': dict
        }

        self.attribute_map = {
            'feedback': 'feedback'
        }

        self._feedback = feedback

    @property
    def feedback(self) -> dict:
        """Gets the feedback.


        :return: The feedback for request
        :rtype: dict
        """
        return self._feedback

    @feedback.setter
    def feedback(self, feedback: dict):
        """Sets the feedback for request.


        :param feedback: The feedback for request.
        :type feedback: dict
        """

        self._feedback = feedback

    @classmethod
    def from_dict(cls, dikt) -> 'FeedbackModelFeedbackRequest':
        """Returns the dict as a model

        :param dikt: A dict.
        :type: dict
        :return: The feedback.ModelFeedbackRequest of this FeedbackModelFeedbackRequest.  # noqa: E501
        :rtype: FeedbackModelFeedbackRequest
        """
        return util.deserialize_model(dikt, cls)
