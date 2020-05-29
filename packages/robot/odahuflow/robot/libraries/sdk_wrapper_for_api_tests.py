from odahuflow.robot.cloud import object_storage
from odahuflow.sdk.clients.api_aggregated import parse_resources_file_with_one_item

from odahuflow.sdk.clients.configuration import ConfigurationClient
from odahuflow.sdk.clients.connection import ConnectionClient
from odahuflow.sdk.clients.deployment import ModelDeploymentClient
from odahuflow.sdk.clients.training import ModelTrainingClient


class SDKWrapper:

    def __init__(self, file):
        api_object = parse_resources_file_with_one_item(file).resource


class Connection(SDKWrapper, ConnectionClient):

    def create(self, api_object):
        ConnectionClient.create(self, api_object)
