# coding: utf-8

from __future__ import absolute_import
from datetime import date, datetime  # noqa: F401

from typing import List, Dict  # noqa: F401

from odahuflow.sdk.models.base_model_ import Model
from odahuflow.sdk.models import util


class FeedbackModelFeedbackResponse(Model):
    """NOTE: This class is manually edited.

    Requires replacement with auto generated model.
    """

    def __init__(self, response):  # noqa: E501
        """FeedbackModelFeedbackResponse - a model defined in Swagger

        """
        self.swagger_types = {
            'response': dict
        }

        self.attribute_map = {
            'response': 'response'
        }

        self._response = response

    @property
    def response(self) -> dict:
        """Gets the response of feedback request.


        :return: The response of feedback request
        :rtype: dict
        """
        return self._response

    @response.setter
    def response(self, response: dict):
        """Sets the response of request.


        :param response: The response of request.
        :type response: dict
        """

        self._response = response

    @classmethod
    def from_dict(cls, dikt) -> 'FeedbackModelFeedbackResponse':
        """Returns the dict as a model

        :param dikt: A dict.
        :type: dict
        :return: The feedback.ModelFeedbackResponse of this FeedbackModelFeedbackResponse.  # noqa: E501
        :rtype: FeedbackModelFeedbackResponse
        """
        return util.deserialize_model(dikt, cls)
