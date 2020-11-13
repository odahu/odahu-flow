import json
from odahuflow.sdk import config

from odahuflow.sdk.clients.api_aggregated import parse_resources_file_with_one_item
from odahuflow.sdk.clients.configuration import ConfigurationClient
from odahuflow.sdk.clients.connection import ConnectionClient
from odahuflow.sdk.clients.deployment import ModelDeploymentClient
from odahuflow.sdk.clients.model import ModelClient
from odahuflow.sdk.clients.packaging import ModelPackagingClient
from odahuflow.sdk.clients.packaging_integration import PackagingIntegrationClient
from odahuflow.sdk.clients.route import ModelRouteClient
from odahuflow.sdk.clients.toolchain_integration import ToolchainIntegrationClient
from odahuflow.sdk.clients.training import ModelTrainingClient


class Login:

    @staticmethod
    def config_get(base_url: str = config.API_URL, **kwargs):
        return ConfigurationClient(base_url=base_url, **kwargs).get()


class Configuration:

    @staticmethod
    def config_get(**kwargs):
        return ConfigurationClient(**kwargs).get()


class Connection:

    @staticmethod
    def connection_get(**kwargs):
        return ConnectionClient(**kwargs).get_all()

    @staticmethod
    def connection_get_id(conn_id: str, **kwargs):
        return ConnectionClient(**kwargs).get(conn_id)

    @staticmethod
    def connection_get_id_decrypted(conn_id: str, **kwargs):
        return ConnectionClient(**kwargs).get_decrypted(conn_id)

    @staticmethod
    def connection_put(payload_file, **kwargs):
        api_object = parse_resources_file_with_one_item(payload_file).resource
        return ConnectionClient(**kwargs).edit(api_object)

    @staticmethod
    def connection_post(payload_file, **kwargs):
        api_object = parse_resources_file_with_one_item(payload_file).resource
        return ConnectionClient(**kwargs).create(api_object)

    @staticmethod
    def connection_delete(conn_id: str, **kwargs):
        return ConnectionClient(**kwargs).delete(conn_id)


class ModelDeployment:

    @staticmethod
    def deployment_get(**kwargs):
        return ModelDeploymentClient(**kwargs).get_all()

    @staticmethod
    def deployment_get_id(dep_id: str, **kwargs):
        return ModelDeploymentClient(**kwargs).get(dep_id)

    @staticmethod
    def deployment_put(payload_file, image=None, **kwargs):
        api_object = parse_resources_file_with_one_item(payload_file).resource

        if image:
            api_object.spec.image = image

        return ModelDeploymentClient().edit(api_object)

    @staticmethod
    def deployment_post(payload_file, image=None, **kwargs):
        api_object = parse_resources_file_with_one_item(payload_file).resource

        if image:
            api_object.spec.image = image

        return ModelDeploymentClient().create(api_object)

    @staticmethod
    def deployment_delete(dep_id: str, **kwargs):
        return ModelDeploymentClient(**kwargs).delete(dep_id)


class ModelPackaging:

    @staticmethod
    def packaging_get(**kwargs):
        return ModelPackagingClient(**kwargs).get_all()

    @staticmethod
    def packaging_get_id(pack_id: str, **kwargs):
        return ModelPackagingClient(**kwargs).get(pack_id)

    @staticmethod
    def packaging_put(payload_file, artifact_name=None, **kwargs):
        api_object = parse_resources_file_with_one_item(payload_file).resource

        if artifact_name:
            api_object.spec.artifact_name = artifact_name

        return ModelPackagingClient(**kwargs).edit(api_object)

    @staticmethod
    def packaging_post(payload_file, artifact_name=None, **kwargs):
        api_object = parse_resources_file_with_one_item(payload_file).resource

        if artifact_name:
            api_object.spec.artifact_name = artifact_name

        return ModelPackagingClient(**kwargs).create(api_object)

    @staticmethod
    def packaging_delete(pack_id: str, **kwargs):
        return ModelPackagingClient(**kwargs).delete(pack_id)

    @staticmethod
    def packaging_get_log(pack_id, **kwargs):
        log_generator = ModelPackagingClient(**kwargs).log(pack_id, follow=False)
        # logs_list will be list of log lines
        logs_list = list(log_generator)
        text = "\n".join(logs_list)
        return text


class ModelTraining:

    @staticmethod
    def training_get(**kwargs):
        return ModelTrainingClient(**kwargs).get_all()

    @staticmethod
    def training_get_id(train_id: str, **kwargs):
        return ModelTrainingClient(**kwargs).get(train_id)

    @staticmethod
    def training_put(payload_file, **kwargs):
        api_object = parse_resources_file_with_one_item(payload_file).resource
        return ModelTrainingClient(**kwargs).edit(api_object)

    @staticmethod
    def training_post(payload_file, **kwargs):
        api_object = parse_resources_file_with_one_item(payload_file).resource
        return ModelTrainingClient(**kwargs).create(api_object)

    @staticmethod
    def training_delete(train_id: str, **kwargs):
        return ModelTrainingClient(**kwargs).delete(train_id)

    @staticmethod
    def training_get_log(train_id, **kwargs):
        log_generator = ModelTrainingClient(**kwargs).log(train_id, follow=False)
        # logs_list will be list of log lines
        logs_list = list(log_generator)
        text = "\n".join(logs_list)
        return text


class ModelRoute:

    @staticmethod
    def route_get(**kwargs):
        return ModelRouteClient(**kwargs).get_all()

    @staticmethod
    def route_get_id(route_id: str, **kwargs):
        return ModelRouteClient(**kwargs).get(route_id)

    @staticmethod
    def route_put(payload_file, **kwargs):
        api_object = parse_resources_file_with_one_item(payload_file).resource
        return ModelRouteClient(**kwargs).edit(api_object)

    @staticmethod
    def route_post(payload_file, **kwargs):
        api_object = parse_resources_file_with_one_item(payload_file).resource
        return ModelRouteClient(**kwargs).create(api_object)

    @staticmethod
    def route_delete(route_id: str, **kwargs):
        return ModelRouteClient(**kwargs).delete(route_id)


class Packager:

    @staticmethod
    def packager_get(**kwargs):
        return PackagingIntegrationClient(**kwargs).get_all()

    @staticmethod
    def packager_get_id(pi_id: str, **kwargs):
        return PackagingIntegrationClient(**kwargs).get(pi_id)

    @staticmethod
    def packager_put(payload_file, **kwargs):
        api_object = parse_resources_file_with_one_item(payload_file).resource
        return PackagingIntegrationClient(**kwargs).edit(api_object)

    @staticmethod
    def packager_post(payload_file, **kwargs):
        api_object = parse_resources_file_with_one_item(payload_file).resource
        return PackagingIntegrationClient(**kwargs).create(api_object)

    @staticmethod
    def packager_delete(pi_id: str, **kwargs):
        return PackagingIntegrationClient(**kwargs).delete(pi_id)


class Toolchain:

    @staticmethod
    def toolchain_get(**kwargs):
        return ToolchainIntegrationClient(**kwargs).get_all()

    @staticmethod
    def toolchain_get_id(ti_id: str, **kwargs):
        return ToolchainIntegrationClient(**kwargs).get(ti_id)

    @staticmethod
    def toolchain_put(payload_file, **kwargs):
        api_object = parse_resources_file_with_one_item(payload_file).resource
        return ToolchainIntegrationClient(**kwargs).edit(api_object)

    @staticmethod
    def toolchain_post(payload_file, **kwargs):
        api_object = parse_resources_file_with_one_item(payload_file).resource
        return ToolchainIntegrationClient(**kwargs).create(api_object)

    @staticmethod
    def toolchain_delete(ti_id: str, **kwargs):
        return ToolchainIntegrationClient(**kwargs).delete(ti_id)


class Model:

    @staticmethod
    def model_get(**kwargs):
        return ModelClient(**kwargs).info()

    @staticmethod
    def model_post(json_input=None, **kwargs):
        return ModelClient(**kwargs).invoke(**json.loads(json_input))
