# coding: utf-8

from __future__ import absolute_import
from datetime import date, datetime  # noqa: F401

from typing import List, Dict  # noqa: F401

from odahuflow.sdk.models.base_model_ import Model
from odahuflow.sdk.models.node_pool import NodePool  # noqa: F401,E501
from odahuflow.sdk.models.resource_requirements import (
    ResourceRequirements,
)  # noqa: F401,E501
from odahuflow.sdk.models import util


class ModelTrainingConfig(Model):
    """NOTE: This class is auto generated by the swagger code generator program.

    Do not edit the class manually.
    """

    def __init__(
        self,
        default_resources: ResourceRequirements = None,
        enabled: bool = None,
        gpu_node_pools: List[NodePool] = None,
        gpu_tolerations: str = None,
        metric_url: str = None,
        model_trainer_image: str = None,
        namespace: str = None,
        node_pools: List[NodePool] = None,
        output_connection_id: str = None,
        service_account: str = None,
        timeout: str = None,
        tolerations: str = None,
        toolchain_integration_namespace: str = None,
        toolchain_integration_repository_type: str = None,
    ):  # noqa: E501
        """ModelTrainingConfig - a model defined in Swagger

        :param default_resources: The default_resources of this ModelTrainingConfig.  # noqa: E501
        :type default_resources: ResourceRequirements
        :param enabled: The enabled of this ModelTrainingConfig.  # noqa: E501
        :type enabled: bool
        :param gpu_node_pools: The gpu_node_pools of this ModelTrainingConfig.  # noqa: E501
        :type gpu_node_pools: List[NodePool]
        :param gpu_tolerations: The gpu_tolerations of this ModelTrainingConfig.  # noqa: E501
        :type gpu_tolerations: str
        :param metric_url: The metric_url of this ModelTrainingConfig.  # noqa: E501
        :type metric_url: str
        :param model_trainer_image: The model_trainer_image of this ModelTrainingConfig.  # noqa: E501
        :type model_trainer_image: str
        :param namespace: The namespace of this ModelTrainingConfig.  # noqa: E501
        :type namespace: str
        :param node_pools: The node_pools of this ModelTrainingConfig.  # noqa: E501
        :type node_pools: List[NodePool]
        :param output_connection_id: The output_connection_id of this ModelTrainingConfig.  # noqa: E501
        :type output_connection_id: str
        :param service_account: The service_account of this ModelTrainingConfig.  # noqa: E501
        :type service_account: str
        :param timeout: The timeout of this ModelTrainingConfig.  # noqa: E501
        :type timeout: str
        :param tolerations: The tolerations of this ModelTrainingConfig.  # noqa: E501
        :type tolerations: str
        :param toolchain_integration_namespace: The toolchain_integration_namespace of this ModelTrainingConfig.  # noqa: E501
        :type toolchain_integration_namespace: str
        :param toolchain_integration_repository_type: The toolchain_integration_repository_type of this ModelTrainingConfig.  # noqa: E501
        :type toolchain_integration_repository_type: str
        """
        self.swagger_types = {
            "default_resources": ResourceRequirements,
            "enabled": bool,
            "gpu_node_pools": List[NodePool],
            "gpu_tolerations": str,
            "metric_url": str,
            "model_trainer_image": str,
            "namespace": str,
            "node_pools": List[NodePool],
            "output_connection_id": str,
            "service_account": str,
            "timeout": str,
            "tolerations": str,
            "toolchain_integration_namespace": str,
            "toolchain_integration_repository_type": str,
        }

        self.attribute_map = {
            "default_resources": "defaultResources",
            "enabled": "enabled",
            "gpu_node_pools": "gpuNodePools",
            "gpu_tolerations": "gpuTolerations",
            "metric_url": "metricUrl",
            "model_trainer_image": "modelTrainerImage",
            "namespace": "namespace",
            "node_pools": "nodePools",
            "output_connection_id": "outputConnectionID",
            "service_account": "serviceAccount",
            "timeout": "timeout",
            "tolerations": "tolerations",
            "toolchain_integration_namespace": "toolchainIntegrationNamespace",
            "toolchain_integration_repository_type": "toolchainIntegrationRepositoryType",
        }

        self._default_resources = default_resources
        self._enabled = enabled
        self._gpu_node_pools = gpu_node_pools
        self._gpu_tolerations = gpu_tolerations
        self._metric_url = metric_url
        self._model_trainer_image = model_trainer_image
        self._namespace = namespace
        self._node_pools = node_pools
        self._output_connection_id = output_connection_id
        self._service_account = service_account
        self._timeout = timeout
        self._tolerations = tolerations
        self._toolchain_integration_namespace = toolchain_integration_namespace
        self._toolchain_integration_repository_type = (
            toolchain_integration_repository_type
        )

    @classmethod
    def from_dict(cls, dikt) -> "ModelTrainingConfig":
        """Returns the dict as a model

        :param dikt: A dict.
        :type: dict
        :return: The ModelTrainingConfig of this ModelTrainingConfig.  # noqa: E501
        :rtype: ModelTrainingConfig
        """
        return util.deserialize_model(dikt, cls)

    @property
    def default_resources(self) -> ResourceRequirements:
        """Gets the default_resources of this ModelTrainingConfig.

        Default resources for training pods  # noqa: E501

        :return: The default_resources of this ModelTrainingConfig.
        :rtype: ResourceRequirements
        """
        return self._default_resources

    @default_resources.setter
    def default_resources(self, default_resources: ResourceRequirements):
        """Sets the default_resources of this ModelTrainingConfig.

        Default resources for training pods  # noqa: E501

        :param default_resources: The default_resources of this ModelTrainingConfig.
        :type default_resources: ResourceRequirements
        """

        self._default_resources = default_resources

    @property
    def enabled(self) -> bool:
        """Gets the enabled of this ModelTrainingConfig.

        Enable deployment API/operator  # noqa: E501

        :return: The enabled of this ModelTrainingConfig.
        :rtype: bool
        """
        return self._enabled

    @enabled.setter
    def enabled(self, enabled: bool):
        """Sets the enabled of this ModelTrainingConfig.

        Enable deployment API/operator  # noqa: E501

        :param enabled: The enabled of this ModelTrainingConfig.
        :type enabled: bool
        """

        self._enabled = enabled

    @property
    def gpu_node_pools(self) -> List[NodePool]:
        """Gets the gpu_node_pools of this ModelTrainingConfig.

        Node pools to run GPU training tasks on  # noqa: E501

        :return: The gpu_node_pools of this ModelTrainingConfig.
        :rtype: List[NodePool]
        """
        return self._gpu_node_pools

    @gpu_node_pools.setter
    def gpu_node_pools(self, gpu_node_pools: List[NodePool]):
        """Sets the gpu_node_pools of this ModelTrainingConfig.

        Node pools to run GPU training tasks on  # noqa: E501

        :param gpu_node_pools: The gpu_node_pools of this ModelTrainingConfig.
        :type gpu_node_pools: List[NodePool]
        """

        self._gpu_node_pools = gpu_node_pools

    @property
    def gpu_tolerations(self) -> str:
        """Gets the gpu_tolerations of this ModelTrainingConfig.

        Kubernetes tolerations for GPU model trainings pods  # noqa: E501

        :return: The gpu_tolerations of this ModelTrainingConfig.
        :rtype: str
        """
        return self._gpu_tolerations

    @gpu_tolerations.setter
    def gpu_tolerations(self, gpu_tolerations: str):
        """Sets the gpu_tolerations of this ModelTrainingConfig.

        Kubernetes tolerations for GPU model trainings pods  # noqa: E501

        :param gpu_tolerations: The gpu_tolerations of this ModelTrainingConfig.
        :type gpu_tolerations: str
        """

        self._gpu_tolerations = gpu_tolerations

    @property
    def metric_url(self) -> str:
        """Gets the metric_url of this ModelTrainingConfig.


        :return: The metric_url of this ModelTrainingConfig.
        :rtype: str
        """
        return self._metric_url

    @metric_url.setter
    def metric_url(self, metric_url: str):
        """Sets the metric_url of this ModelTrainingConfig.


        :param metric_url: The metric_url of this ModelTrainingConfig.
        :type metric_url: str
        """

        self._metric_url = metric_url

    @property
    def model_trainer_image(self) -> str:
        """Gets the model_trainer_image of this ModelTrainingConfig.


        :return: The model_trainer_image of this ModelTrainingConfig.
        :rtype: str
        """
        return self._model_trainer_image

    @model_trainer_image.setter
    def model_trainer_image(self, model_trainer_image: str):
        """Sets the model_trainer_image of this ModelTrainingConfig.


        :param model_trainer_image: The model_trainer_image of this ModelTrainingConfig.
        :type model_trainer_image: str
        """

        self._model_trainer_image = model_trainer_image

    @property
    def namespace(self) -> str:
        """Gets the namespace of this ModelTrainingConfig.

        Kubernetes namespace, where model trainings will be deployed  # noqa: E501

        :return: The namespace of this ModelTrainingConfig.
        :rtype: str
        """
        return self._namespace

    @namespace.setter
    def namespace(self, namespace: str):
        """Sets the namespace of this ModelTrainingConfig.

        Kubernetes namespace, where model trainings will be deployed  # noqa: E501

        :param namespace: The namespace of this ModelTrainingConfig.
        :type namespace: str
        """

        self._namespace = namespace

    @property
    def node_pools(self) -> List[NodePool]:
        """Gets the node_pools of this ModelTrainingConfig.

        Node pools to run training tasks on  # noqa: E501

        :return: The node_pools of this ModelTrainingConfig.
        :rtype: List[NodePool]
        """
        return self._node_pools

    @node_pools.setter
    def node_pools(self, node_pools: List[NodePool]):
        """Sets the node_pools of this ModelTrainingConfig.

        Node pools to run training tasks on  # noqa: E501

        :param node_pools: The node_pools of this ModelTrainingConfig.
        :type node_pools: List[NodePool]
        """

        self._node_pools = node_pools

    @property
    def output_connection_id(self) -> str:
        """Gets the output_connection_id of this ModelTrainingConfig.


        :return: The output_connection_id of this ModelTrainingConfig.
        :rtype: str
        """
        return self._output_connection_id

    @output_connection_id.setter
    def output_connection_id(self, output_connection_id: str):
        """Sets the output_connection_id of this ModelTrainingConfig.


        :param output_connection_id: The output_connection_id of this ModelTrainingConfig.
        :type output_connection_id: str
        """

        self._output_connection_id = output_connection_id

    @property
    def service_account(self) -> str:
        """Gets the service_account of this ModelTrainingConfig.


        :return: The service_account of this ModelTrainingConfig.
        :rtype: str
        """
        return self._service_account

    @service_account.setter
    def service_account(self, service_account: str):
        """Sets the service_account of this ModelTrainingConfig.


        :param service_account: The service_account of this ModelTrainingConfig.
        :type service_account: str
        """

        self._service_account = service_account

    @property
    def timeout(self) -> str:
        """Gets the timeout of this ModelTrainingConfig.

        Timeout for full training process  # noqa: E501

        :return: The timeout of this ModelTrainingConfig.
        :rtype: str
        """
        return self._timeout

    @timeout.setter
    def timeout(self, timeout: str):
        """Sets the timeout of this ModelTrainingConfig.

        Timeout for full training process  # noqa: E501

        :param timeout: The timeout of this ModelTrainingConfig.
        :type timeout: str
        """

        self._timeout = timeout

    @property
    def tolerations(self) -> str:
        """Gets the tolerations of this ModelTrainingConfig.

        Kubernetes tolerations for model trainings pods  # noqa: E501

        :return: The tolerations of this ModelTrainingConfig.
        :rtype: str
        """
        return self._tolerations

    @tolerations.setter
    def tolerations(self, tolerations: str):
        """Sets the tolerations of this ModelTrainingConfig.

        Kubernetes tolerations for model trainings pods  # noqa: E501

        :param tolerations: The tolerations of this ModelTrainingConfig.
        :type tolerations: str
        """

        self._tolerations = tolerations

    @property
    def toolchain_integration_namespace(self) -> str:
        """Gets the toolchain_integration_namespace of this ModelTrainingConfig.


        :return: The toolchain_integration_namespace of this ModelTrainingConfig.
        :rtype: str
        """
        return self._toolchain_integration_namespace

    @toolchain_integration_namespace.setter
    def toolchain_integration_namespace(self, toolchain_integration_namespace: str):
        """Sets the toolchain_integration_namespace of this ModelTrainingConfig.


        :param toolchain_integration_namespace: The toolchain_integration_namespace of this ModelTrainingConfig.
        :type toolchain_integration_namespace: str
        """

        self._toolchain_integration_namespace = toolchain_integration_namespace

    @property
    def toolchain_integration_repository_type(self) -> str:
        """Gets the toolchain_integration_repository_type of this ModelTrainingConfig.

        Storage backend for toolchain integrations. Available options:   * kubernetes   * postgres  # noqa: E501

        :return: The toolchain_integration_repository_type of this ModelTrainingConfig.
        :rtype: str
        """
        return self._toolchain_integration_repository_type

    @toolchain_integration_repository_type.setter
    def toolchain_integration_repository_type(
        self, toolchain_integration_repository_type: str
    ):
        """Sets the toolchain_integration_repository_type of this ModelTrainingConfig.

        Storage backend for toolchain integrations. Available options:   * kubernetes   * postgres  # noqa: E501

        :param toolchain_integration_repository_type: The toolchain_integration_repository_type of this ModelTrainingConfig.
        :type toolchain_integration_repository_type: str
        """

        self._toolchain_integration_repository_type = (
            toolchain_integration_repository_type
        )
