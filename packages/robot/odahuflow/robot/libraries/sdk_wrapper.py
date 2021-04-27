# for class Model
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
from odahuflow.sdk.clients.batch_service import BatchInferenceServiceClient
from odahuflow.sdk.clients.batch_job import BatchInferenceJobClient
from odahuflow.sdk.clients.api import EntityAlreadyExists


class Login:

    @staticmethod
    def reload_config():
        config._INI_FILE_TRIED_TO_BE_LOADED = False
        config.reinitialize_variables()


class Configuration:

    @staticmethod
    def config_get(**kwargs):
        return ConfigurationClient(**kwargs).get()


class Connection:

    @staticmethod
    def connection_get():
        return ConnectionClient().get_all()

    @staticmethod
    def connection_get_id(conn_id: str):
        return ConnectionClient().get(conn_id)

    @staticmethod
    def connection_get_id_decrypted(conn_id: str):
        return ConnectionClient().get_decrypted(conn_id)

    @staticmethod
    def connection_put(payload_file):
        api_object = parse_resources_file_with_one_item(payload_file).resource
        return ConnectionClient().edit(api_object)

    @staticmethod
    def connection_post(payload_file):
        api_object = parse_resources_file_with_one_item(payload_file).resource
        return ConnectionClient().create(api_object)

    @staticmethod
    def connection_delete(conn_id: str):
        return ConnectionClient().delete(conn_id)


class ModelDeployment:

    @staticmethod
    def deployment_get():
        return ModelDeploymentClient().get_all()

    @staticmethod
    def deployment_get_id(dep_id: str):
        return ModelDeploymentClient().get(dep_id)

    @staticmethod
    def deployment_put(payload_file, image=None):
        api_object = parse_resources_file_with_one_item(payload_file).resource

        if image:
            api_object.spec.image = image

        return ModelDeploymentClient().edit(api_object)

    @staticmethod
    def deployment_post(payload_file, image=None):
        api_object = parse_resources_file_with_one_item(payload_file).resource

        if image:
            api_object.spec.image = image

        return ModelDeploymentClient().create(api_object)

    @staticmethod
    def deployment_delete(dep_id: str):
        return ModelDeploymentClient().delete(dep_id)

    @staticmethod
    def deployment_get_default_route(dep_id: str):
        return ModelDeploymentClient().get_default_route(dep_id)


class ModelPackaging:

    @staticmethod
    def packaging_get():
        return ModelPackagingClient().get_all()

    @staticmethod
    def packaging_get_id(pack_id: str):
        return ModelPackagingClient().get(pack_id)

    @staticmethod
    def packaging_put(payload_file, artifact_name=None):
        api_object = parse_resources_file_with_one_item(payload_file).resource

        if artifact_name:
            api_object.spec.artifact_name = artifact_name

        return ModelPackagingClient().edit(api_object)

    @staticmethod
    def packaging_post(payload_file, artifact_name=None):
        api_object = parse_resources_file_with_one_item(payload_file).resource

        if artifact_name:
            api_object.spec.artifact_name = artifact_name

        return ModelPackagingClient().create(api_object)

    @staticmethod
    def packaging_delete(pack_id: str):
        return ModelPackagingClient().delete(pack_id)

    @staticmethod
    def packaging_get_log(pack_id):
        log_generator = ModelPackagingClient().log(pack_id, follow=False)
        # logs_list will be list of log lines
        logs_list = list(log_generator)
        text = "\n".join(logs_list)
        return text


class ModelTraining:

    @staticmethod
    def training_get():
        return ModelTrainingClient().get_all()

    @staticmethod
    def training_get_id(train_id: str):
        return ModelTrainingClient().get(train_id)

    @staticmethod
    def training_put(payload_file):
        api_object = parse_resources_file_with_one_item(payload_file).resource
        return ModelTrainingClient().edit(api_object)

    @staticmethod
    def training_post(payload_file):
        api_object = parse_resources_file_with_one_item(payload_file).resource
        return ModelTrainingClient().create(api_object)

    @staticmethod
    def training_delete(train_id: str):
        return ModelTrainingClient().delete(train_id)

    @staticmethod
    def training_get_log(train_id):
        log_generator = ModelTrainingClient().log(train_id, follow=False)
        # logs_list will be list of log lines
        logs_list = list(log_generator)
        text = "\n".join(logs_list)
        return text


class ModelRoute:

    @staticmethod
    def route_get():
        return ModelRouteClient().get_all()

    @staticmethod
    def route_get_id(route_id: str):
        return ModelRouteClient().get(route_id)

    @staticmethod
    def route_put(payload_file):
        api_object = parse_resources_file_with_one_item(payload_file).resource
        return ModelRouteClient().edit(api_object)

    @staticmethod
    def route_post(payload_file):
        api_object = parse_resources_file_with_one_item(payload_file).resource
        return ModelRouteClient().create(api_object)

    @staticmethod
    def route_delete(route_id: str):
        return ModelRouteClient().delete(route_id)


class Packager:

    @staticmethod
    def packager_get():
        return PackagingIntegrationClient().get_all()

    @staticmethod
    def packager_get_id(pi_id: str):
        return PackagingIntegrationClient().get(pi_id)

    @staticmethod
    def packager_put(payload_file):
        api_object = parse_resources_file_with_one_item(payload_file).resource
        return PackagingIntegrationClient().edit(api_object)

    @staticmethod
    def packager_post(payload_file):
        api_object = parse_resources_file_with_one_item(payload_file).resource
        return PackagingIntegrationClient().create(api_object)

    @staticmethod
    def packager_delete(pi_id: str):
        return PackagingIntegrationClient().delete(pi_id)


class Toolchain:

    @staticmethod
    def toolchain_get():
        return ToolchainIntegrationClient().get_all()

    @staticmethod
    def toolchain_get_id(ti_id: str):
        return ToolchainIntegrationClient().get(ti_id)

    @staticmethod
    def toolchain_put(payload_file):
        api_object = parse_resources_file_with_one_item(payload_file).resource
        return ToolchainIntegrationClient().edit(api_object)

    @staticmethod
    def toolchain_post(payload_file):
        api_object = parse_resources_file_with_one_item(payload_file).resource
        return ToolchainIntegrationClient().create(api_object)

    @staticmethod
    def toolchain_delete(ti_id: str):
        return ToolchainIntegrationClient().delete(ti_id)


class Model:

    @staticmethod
    def model_get(base_url, model_route=None, model_deployment=None, url_prefix=None, **kwargs):
        return ModelClient(
            base_url,
            model_route=model_route,
            model_deployment=model_deployment,
            url_prefix=url_prefix,
            token=config.API_TOKEN
        ).info()

    @staticmethod
    def model_post(base_url, model_route=None, model_deployment=None, url_prefix=None, json_input=None, **kwargs):
        return ModelClient(
            base_url,
            model_route=model_route,
            model_deployment=model_deployment,
            url_prefix=url_prefix,
            token=config.API_TOKEN
        ).invoke(**json.loads(json_input))


class InferenceService:

    @staticmethod
    def service_post(payload_file):
        api_object = parse_resources_file_with_one_item(payload_file).resource
        try:
            BatchInferenceServiceClient().create(api_object)
        except EntityAlreadyExists:
            pass


class InferenceJob:

    @staticmethod
    def job_post(payload_file):
        api_object = parse_resources_file_with_one_item(payload_file).resource
        en = BatchInferenceJobClient().create(api_object)
        return en.id

    @staticmethod
    def job_get_id(id_: str):
        return BatchInferenceJobClient().get(id_)
