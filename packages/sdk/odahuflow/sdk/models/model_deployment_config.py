# coding: utf-8

from __future__ import absolute_import
from datetime import date, datetime  # noqa: F401

from typing import List, Dict  # noqa: F401

from odahuflow.sdk.models.base_model_ import Model
from odahuflow.sdk.models.edge_config import EdgeConfig  # noqa: F401,E501
from odahuflow.sdk.models.model_deployment_istio_config import ModelDeploymentIstioConfig  # noqa: F401,E501
from odahuflow.sdk.models.model_deployment_security_config import ModelDeploymentSecurityConfig  # noqa: F401,E501
from odahuflow.sdk.models import util


class ModelDeploymentConfig(Model):
    """NOTE: This class is auto generated by the swagger code generator program.

    Do not edit the class manually.
    """

    def __init__(self, default_docker_pull_conn_name: str=None, edge: EdgeConfig=None, enabled: bool=None, istio: ModelDeploymentIstioConfig=None, namespace: str=None, node_selector: Dict[str, str]=None, security: ModelDeploymentSecurityConfig=None, toleration: Dict[str, str]=None):  # noqa: E501
        """ModelDeploymentConfig - a model defined in Swagger

        :param default_docker_pull_conn_name: The default_docker_pull_conn_name of this ModelDeploymentConfig.  # noqa: E501
        :type default_docker_pull_conn_name: str
        :param edge: The edge of this ModelDeploymentConfig.  # noqa: E501
        :type edge: EdgeConfig
        :param enabled: The enabled of this ModelDeploymentConfig.  # noqa: E501
        :type enabled: bool
        :param istio: The istio of this ModelDeploymentConfig.  # noqa: E501
        :type istio: ModelDeploymentIstioConfig
        :param namespace: The namespace of this ModelDeploymentConfig.  # noqa: E501
        :type namespace: str
        :param node_selector: The node_selector of this ModelDeploymentConfig.  # noqa: E501
        :type node_selector: Dict[str, str]
        :param security: The security of this ModelDeploymentConfig.  # noqa: E501
        :type security: ModelDeploymentSecurityConfig
        :param toleration: The toleration of this ModelDeploymentConfig.  # noqa: E501
        :type toleration: Dict[str, str]
        """
        self.swagger_types = {
            'default_docker_pull_conn_name': str,
            'edge': EdgeConfig,
            'enabled': bool,
            'istio': ModelDeploymentIstioConfig,
            'namespace': str,
            'node_selector': Dict[str, str],
            'security': ModelDeploymentSecurityConfig,
            'toleration': Dict[str, str]
        }

        self.attribute_map = {
            'default_docker_pull_conn_name': 'defaultDockerPullConnName',
            'edge': 'edge',
            'enabled': 'enabled',
            'istio': 'istio',
            'namespace': 'namespace',
            'node_selector': 'nodeSelector',
            'security': 'security',
            'toleration': 'toleration'
        }

        self._default_docker_pull_conn_name = default_docker_pull_conn_name
        self._edge = edge
        self._enabled = enabled
        self._istio = istio
        self._namespace = namespace
        self._node_selector = node_selector
        self._security = security
        self._toleration = toleration

    @classmethod
    def from_dict(cls, dikt) -> 'ModelDeploymentConfig':
        """Returns the dict as a model

        :param dikt: A dict.
        :type: dict
        :return: The ModelDeploymentConfig of this ModelDeploymentConfig.  # noqa: E501
        :rtype: ModelDeploymentConfig
        """
        return util.deserialize_model(dikt, cls)

    @property
    def default_docker_pull_conn_name(self) -> str:
        """Gets the default_docker_pull_conn_name of this ModelDeploymentConfig.

        Default connection ID which will be used if a user doesn't specify it in a model deployment  # noqa: E501

        :return: The default_docker_pull_conn_name of this ModelDeploymentConfig.
        :rtype: str
        """
        return self._default_docker_pull_conn_name

    @default_docker_pull_conn_name.setter
    def default_docker_pull_conn_name(self, default_docker_pull_conn_name: str):
        """Sets the default_docker_pull_conn_name of this ModelDeploymentConfig.

        Default connection ID which will be used if a user doesn't specify it in a model deployment  # noqa: E501

        :param default_docker_pull_conn_name: The default_docker_pull_conn_name of this ModelDeploymentConfig.
        :type default_docker_pull_conn_name: str
        """

        self._default_docker_pull_conn_name = default_docker_pull_conn_name

    @property
    def edge(self) -> EdgeConfig:
        """Gets the edge of this ModelDeploymentConfig.


        :return: The edge of this ModelDeploymentConfig.
        :rtype: EdgeConfig
        """
        return self._edge

    @edge.setter
    def edge(self, edge: EdgeConfig):
        """Sets the edge of this ModelDeploymentConfig.


        :param edge: The edge of this ModelDeploymentConfig.
        :type edge: EdgeConfig
        """

        self._edge = edge

    @property
    def enabled(self) -> bool:
        """Gets the enabled of this ModelDeploymentConfig.

        Enable deployment API/operator  # noqa: E501

        :return: The enabled of this ModelDeploymentConfig.
        :rtype: bool
        """
        return self._enabled

    @enabled.setter
    def enabled(self, enabled: bool):
        """Sets the enabled of this ModelDeploymentConfig.

        Enable deployment API/operator  # noqa: E501

        :param enabled: The enabled of this ModelDeploymentConfig.
        :type enabled: bool
        """

        self._enabled = enabled

    @property
    def istio(self) -> ModelDeploymentIstioConfig:
        """Gets the istio of this ModelDeploymentConfig.


        :return: The istio of this ModelDeploymentConfig.
        :rtype: ModelDeploymentIstioConfig
        """
        return self._istio

    @istio.setter
    def istio(self, istio: ModelDeploymentIstioConfig):
        """Sets the istio of this ModelDeploymentConfig.


        :param istio: The istio of this ModelDeploymentConfig.
        :type istio: ModelDeploymentIstioConfig
        """

        self._istio = istio

    @property
    def namespace(self) -> str:
        """Gets the namespace of this ModelDeploymentConfig.

        Kubernetes namespace, where model deployments will be deployed  # noqa: E501

        :return: The namespace of this ModelDeploymentConfig.
        :rtype: str
        """
        return self._namespace

    @namespace.setter
    def namespace(self, namespace: str):
        """Sets the namespace of this ModelDeploymentConfig.

        Kubernetes namespace, where model deployments will be deployed  # noqa: E501

        :param namespace: The namespace of this ModelDeploymentConfig.
        :type namespace: str
        """

        self._namespace = namespace

    @property
    def node_selector(self) -> Dict[str, str]:
        """Gets the node_selector of this ModelDeploymentConfig.

        Kubernetes node selector for model deployments  # noqa: E501

        :return: The node_selector of this ModelDeploymentConfig.
        :rtype: Dict[str, str]
        """
        return self._node_selector

    @node_selector.setter
    def node_selector(self, node_selector: Dict[str, str]):
        """Sets the node_selector of this ModelDeploymentConfig.

        Kubernetes node selector for model deployments  # noqa: E501

        :param node_selector: The node_selector of this ModelDeploymentConfig.
        :type node_selector: Dict[str, str]
        """

        self._node_selector = node_selector

    @property
    def security(self) -> ModelDeploymentSecurityConfig:
        """Gets the security of this ModelDeploymentConfig.


        :return: The security of this ModelDeploymentConfig.
        :rtype: ModelDeploymentSecurityConfig
        """
        return self._security

    @security.setter
    def security(self, security: ModelDeploymentSecurityConfig):
        """Sets the security of this ModelDeploymentConfig.


        :param security: The security of this ModelDeploymentConfig.
        :type security: ModelDeploymentSecurityConfig
        """

        self._security = security

    @property
    def toleration(self) -> Dict[str, str]:
        """Gets the toleration of this ModelDeploymentConfig.

        Kubernetes tolerations for model deployments  # noqa: E501

        :return: The toleration of this ModelDeploymentConfig.
        :rtype: Dict[str, str]
        """
        return self._toleration

    @toleration.setter
    def toleration(self, toleration: Dict[str, str]):
        """Sets the toleration of this ModelDeploymentConfig.

        Kubernetes tolerations for model deployments  # noqa: E501

        :param toleration: The toleration of this ModelDeploymentConfig.
        :type toleration: Dict[str, str]
        """

        self._toleration = toleration