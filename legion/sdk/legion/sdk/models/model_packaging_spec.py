# coding: utf-8

from __future__ import absolute_import
from datetime import date, datetime  # noqa: F401

from typing import List, Dict  # noqa: F401

from legion.sdk.models.base_model_ import Model
from legion.sdk.models.target import Target  # noqa: F401,E501
from legion.sdk.models import util


class ModelPackagingSpec(Model):
    """NOTE: This class is auto generated by the swagger code generator program.

    Do not edit the class manually.
    """

    def __init__(self, arguments: object=None, artifact_name: str=None, image: str=None, integration_name: str=None, targets: List[Target]=None):  # noqa: E501
        """ModelPackagingSpec - a model defined in Swagger

        :param arguments: The arguments of this ModelPackagingSpec.  # noqa: E501
        :type arguments: object
        :param artifact_name: The artifact_name of this ModelPackagingSpec.  # noqa: E501
        :type artifact_name: str
        :param image: The image of this ModelPackagingSpec.  # noqa: E501
        :type image: str
        :param integration_name: The integration_name of this ModelPackagingSpec.  # noqa: E501
        :type integration_name: str
        :param targets: The targets of this ModelPackagingSpec.  # noqa: E501
        :type targets: List[Target]
        """
        self.swagger_types = {
            'arguments': object,
            'artifact_name': str,
            'image': str,
            'integration_name': str,
            'targets': List[Target]
        }

        self.attribute_map = {
            'arguments': 'arguments',
            'artifact_name': 'artifactName',
            'image': 'image',
            'integration_name': 'integrationName',
            'targets': 'targets'
        }

        self._arguments = arguments
        self._artifact_name = artifact_name
        self._image = image
        self._integration_name = integration_name
        self._targets = targets

    @classmethod
    def from_dict(cls, dikt) -> 'ModelPackagingSpec':
        """Returns the dict as a model

        :param dikt: A dict.
        :type: dict
        :return: The ModelPackagingSpec of this ModelPackagingSpec.  # noqa: E501
        :rtype: ModelPackagingSpec
        """
        return util.deserialize_model(dikt, cls)

    @property
    def arguments(self) -> object:
        """Gets the arguments of this ModelPackagingSpec.

        List of arguments. This parameter depends on the specific packaging integration  # noqa: E501

        :return: The arguments of this ModelPackagingSpec.
        :rtype: object
        """
        return self._arguments

    @arguments.setter
    def arguments(self, arguments: object):
        """Sets the arguments of this ModelPackagingSpec.

        List of arguments. This parameter depends on the specific packaging integration  # noqa: E501

        :param arguments: The arguments of this ModelPackagingSpec.
        :type arguments: object
        """

        self._arguments = arguments

    @property
    def artifact_name(self) -> str:
        """Gets the artifact_name of this ModelPackagingSpec.

        Training output artifact name  # noqa: E501

        :return: The artifact_name of this ModelPackagingSpec.
        :rtype: str
        """
        return self._artifact_name

    @artifact_name.setter
    def artifact_name(self, artifact_name: str):
        """Sets the artifact_name of this ModelPackagingSpec.

        Training output artifact name  # noqa: E501

        :param artifact_name: The artifact_name of this ModelPackagingSpec.
        :type artifact_name: str
        """

        self._artifact_name = artifact_name

    @property
    def image(self) -> str:
        """Gets the image of this ModelPackagingSpec.

        Image name. Packaging integration image will be used if this parameters is missed  # noqa: E501

        :return: The image of this ModelPackagingSpec.
        :rtype: str
        """
        return self._image

    @image.setter
    def image(self, image: str):
        """Sets the image of this ModelPackagingSpec.

        Image name. Packaging integration image will be used if this parameters is missed  # noqa: E501

        :param image: The image of this ModelPackagingSpec.
        :type image: str
        """

        self._image = image

    @property
    def integration_name(self) -> str:
        """Gets the integration_name of this ModelPackagingSpec.

        Packaging integration ID  # noqa: E501

        :return: The integration_name of this ModelPackagingSpec.
        :rtype: str
        """
        return self._integration_name

    @integration_name.setter
    def integration_name(self, integration_name: str):
        """Sets the integration_name of this ModelPackagingSpec.

        Packaging integration ID  # noqa: E501

        :param integration_name: The integration_name of this ModelPackagingSpec.
        :type integration_name: str
        """

        self._integration_name = integration_name

    @property
    def targets(self) -> List[Target]:
        """Gets the targets of this ModelPackagingSpec.

        List of targets. This parameter depends on the specific packaging integration  # noqa: E501

        :return: The targets of this ModelPackagingSpec.
        :rtype: List[Target]
        """
        return self._targets

    @targets.setter
    def targets(self, targets: List[Target]):
        """Sets the targets of this ModelPackagingSpec.

        List of targets. This parameter depends on the specific packaging integration  # noqa: E501

        :param targets: The targets of this ModelPackagingSpec.
        :type targets: List[Target]
        """

        self._targets = targets
