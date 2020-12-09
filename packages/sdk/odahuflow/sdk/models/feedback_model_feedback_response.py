# coding: utf-8

from __future__ import absolute_import
from datetime import date, datetime  # noqa: F401

from typing import List, Dict  # noqa: F401

from odahuflow.sdk.models.base_model_ import Model
from odahuflow.sdk.models import util


class FeedbackModelFeedbackResponse(Model):
    """NOTE: This class is auto generated by the swagger code generator program.

    Do not edit the class manually.
    """

    def __init__(self):  # noqa: E501
        """FeedbackModelFeedbackResponse - a model defined in Swagger"""
        self.swagger_types = {}

        self.attribute_map = {}

    @classmethod
    def from_dict(cls, dikt) -> "FeedbackModelFeedbackResponse":
        """Returns the dict as a model

        :param dikt: A dict.
        :type: dict
        :return: The feedback.ModelFeedbackResponse of this FeedbackModelFeedbackResponse.  # noqa: E501
        :rtype: FeedbackModelFeedbackResponse
        """
        return util.deserialize_model(dikt, cls)
