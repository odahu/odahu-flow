# coding: utf-8

from __future__ import absolute_import
from datetime import date, datetime  # noqa: F401

from typing import List, Dict  # noqa: F401

from odahuflow.sdk.models.base_model_ import Model
from odahuflow.sdk.models.model_route import ModelRoute  # noqa: F401,E501
from odahuflow.sdk.models import util


class RouteEvent(Model):
    """NOTE: This class is auto generated by the swagger code generator program.

    Do not edit the class manually.
    """

    def __init__(self, _datetime: str=None, entity_id: str=None, payload: ModelRoute=None, type: str=None):  # noqa: E501
        """RouteEvent - a model defined in Swagger

        :param _datetime: The _datetime of this RouteEvent.  # noqa: E501
        :type _datetime: str
        :param entity_id: The entity_id of this RouteEvent.  # noqa: E501
        :type entity_id: str
        :param payload: The payload of this RouteEvent.  # noqa: E501
        :type payload: ModelRoute
        :param type: The type of this RouteEvent.  # noqa: E501
        :type type: str
        """
        self.swagger_types = {
            '_datetime': str,
            'entity_id': str,
            'payload': ModelRoute,
            'type': str
        }

        self.attribute_map = {
            '_datetime': 'datetime',
            'entity_id': 'entityID',
            'payload': 'payload',
            'type': 'type'
        }

        self.__datetime = _datetime
        self._entity_id = entity_id
        self._payload = payload
        self._type = type

    @classmethod
    def from_dict(cls, dikt) -> 'RouteEvent':
        """Returns the dict as a model

        :param dikt: A dict.
        :type: dict
        :return: The RouteEvent of this RouteEvent.  # noqa: E501
        :rtype: RouteEvent
        """
        return util.deserialize_model(dikt, cls)

    @property
    def _datetime(self) -> str:
        """Gets the _datetime of this RouteEvent.

        When event is raised  # noqa: E501

        :return: The _datetime of this RouteEvent.
        :rtype: str
        """
        return self.__datetime

    @_datetime.setter
    def _datetime(self, _datetime: str):
        """Sets the _datetime of this RouteEvent.

        When event is raised  # noqa: E501

        :param _datetime: The _datetime of this RouteEvent.
        :type _datetime: str
        """

        self.__datetime = _datetime

    @property
    def entity_id(self) -> str:
        """Gets the entity_id of this RouteEvent.

        EntityID contains ID of ModelRoute for ModelRouteDeleted and ModelRouteDeletionMarkIsSet event types Does not make sense in case of ModelRouteUpdate, ModelRouteCreate, ModelRouteStatusUpdated events  # noqa: E501

        :return: The entity_id of this RouteEvent.
        :rtype: str
        """
        return self._entity_id

    @entity_id.setter
    def entity_id(self, entity_id: str):
        """Sets the entity_id of this RouteEvent.

        EntityID contains ID of ModelRoute for ModelRouteDeleted and ModelRouteDeletionMarkIsSet event types Does not make sense in case of ModelRouteUpdate, ModelRouteCreate, ModelRouteStatusUpdated events  # noqa: E501

        :param entity_id: The entity_id of this RouteEvent.
        :type entity_id: str
        """

        self._entity_id = entity_id

    @property
    def payload(self) -> ModelRoute:
        """Gets the payload of this RouteEvent.

        Payload contains ModelRoute for ModelRouteUpdate, ModelRouteCreate, ModelRouteStatusUpdated  events. Does not make sense in case of ModelRouteDelete, ModelRouteDeletionMarkIsSet events  # noqa: E501

        :return: The payload of this RouteEvent.
        :rtype: ModelRoute
        """
        return self._payload

    @payload.setter
    def payload(self, payload: ModelRoute):
        """Sets the payload of this RouteEvent.

        Payload contains ModelRoute for ModelRouteUpdate, ModelRouteCreate, ModelRouteStatusUpdated  events. Does not make sense in case of ModelRouteDelete, ModelRouteDeletionMarkIsSet events  # noqa: E501

        :param payload: The payload of this RouteEvent.
        :type payload: ModelRoute
        """

        self._payload = payload

    @property
    def type(self) -> str:
        """Gets the type of this RouteEvent.

        Possible values: ModelRouteCreate, ModelRouteUpdate, ModelRouteDeleted, ModelRouteDeletionMarkIsSet, ModelRouteStatusUpdated  # noqa: E501

        :return: The type of this RouteEvent.
        :rtype: str
        """
        return self._type

    @type.setter
    def type(self, type: str):
        """Sets the type of this RouteEvent.

        Possible values: ModelRouteCreate, ModelRouteUpdate, ModelRouteDeleted, ModelRouteDeletionMarkIsSet, ModelRouteStatusUpdated  # noqa: E501

        :param type: The type of this RouteEvent.
        :type type: str
        """

        self._type = type
