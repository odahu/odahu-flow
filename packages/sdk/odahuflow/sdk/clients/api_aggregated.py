#
#    Copyright 2019 EPAM Systems
#
#    Licensed under the Apache License, Version 2.0 (the "License");
#    you may not use this file except in compliance with the License.
#    You may obtain a copy of the License at
#
#        http://www.apache.org/licenses/LICENSE-2.0
#
#    Unless required by applicable law or agreed to in writing, software
#    distributed under the License is distributed on an "AS IS" BASIS,
#    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
#    See the License for the specific language governing permissions and
#    limitations under the License.
#
"""
Aggregated API client (can apply multiple resources)
"""
import asyncio
import json
import logging
import os
import typing
from os import listdir
from os.path import isdir, join, isfile

import yaml
from odahuflow.sdk.clients.api import RemoteAPIClient, WrongHttpStatusCode, AsyncRemoteAPIClient
from odahuflow.sdk.clients.connection import ConnectionClient, AsyncConnectionClient
from odahuflow.sdk.clients.deployment import ModelDeploymentClient, ModelDeployment, AsyncModelDeploymentClient
from odahuflow.sdk.clients.packaging import ModelPackagingClient, AsyncModelPackagingClient
from odahuflow.sdk.clients.packaging_integration import PackagingIntegrationClient, AsyncPackagingIntegrationClient
from odahuflow.sdk.clients.route import ModelRoute, ModelRouteClient, AsyncModelRouteClient
from odahuflow.sdk.clients.toolchain_integration import ToolchainIntegrationClient, AsyncToolchainIntegrationClient
from odahuflow.sdk.clients.training import ModelTrainingClient, ModelTraining, AsyncModelTrainingClient
from odahuflow.sdk.models import Connection, ToolchainIntegration, ModelPackaging, PackagingIntegration

TARGET_CLASSES = {
    'ModelTraining': ModelTraining,
    'ToolchainIntegration': ToolchainIntegration,
    'ModelDeployment': ModelDeployment,
    'ModelRoute': ModelRoute,
    'Connection': Connection,
    'ModelPackaging': ModelPackaging,
    'PackagingIntegration': PackagingIntegration,
}

LOGGER = logging.getLogger(__name__)


class InvalidResourceType(Exception):
    """
    Invalid resource type (unsupported) exception
    """

    pass


class OdahuflowCloudResourceUpdatePair(typing.NamedTuple):
    """
    Information about resources to update
    """

    resource_id: str
    resource: object


class OdahuflowCloudResourcesUpdateList(typing.NamedTuple):
    """
    Bulk update request (multiple resources)
    """

    changes: typing.Tuple[OdahuflowCloudResourceUpdatePair] = tuple()


class ApplyResult(typing.NamedTuple):
    """
    Result of bulk applying
    """

    created: typing.Tuple[OdahuflowCloudResourceUpdatePair] = tuple()
    removed: typing.Tuple[OdahuflowCloudResourceUpdatePair] = tuple()
    changed: typing.Tuple[OdahuflowCloudResourceUpdatePair] = tuple()
    errors: typing.Tuple[Exception] = tuple()


# pylint: disable=R0911
def build_client(resource: OdahuflowCloudResourceUpdatePair, api_client: RemoteAPIClient) -> typing.Optional[object]:
    """
    Build client for particular resource (e.g. it builds ModelTrainingClient for ModelTraining resource)

    :param resource: target resource
    :param api_client: base API client to extract connection options from
    :return: remote client or None
    """
    if isinstance(resource.resource, ModelTraining):
        return ModelTrainingClient.construct_from_other(api_client)
    elif isinstance(resource.resource, ModelDeployment):
        return ModelDeploymentClient.construct_from_other(api_client)
    elif isinstance(resource.resource, Connection):
        return ConnectionClient.construct_from_other(api_client)
    elif isinstance(resource.resource, ToolchainIntegration):
        return ToolchainIntegrationClient.construct_from_other(api_client)
    elif isinstance(resource.resource, ModelRoute):
        return ModelRouteClient.construct_from_other(api_client)
    elif isinstance(resource.resource, ModelPackaging):
        return ModelPackagingClient.construct_from_other(api_client)
    elif isinstance(resource.resource, PackagingIntegration):
        return PackagingIntegrationClient.construct_from_other(api_client)
    else:
        raise InvalidResourceType('{!r} is invalid resource '.format(resource.resource))


# pylint: disable=R0911
def build_async_client(resource: OdahuflowCloudResourceUpdatePair,
                       async_api_client: AsyncRemoteAPIClient
                       ) -> typing.Optional[object]:
    """
    Build client for particular resource (e.g. it builds ModelTrainingClient for ModelTraining resource)

    :param resource: target resource
    :param async_api_client: base async API client to extract connection options from
    :return: remote client or None
    """
    if isinstance(resource.resource, ModelTraining):
        return AsyncModelTrainingClient.construct_from_other(async_api_client)
    elif isinstance(resource.resource, ModelDeployment):
        return AsyncModelDeploymentClient.construct_from_other(async_api_client)
    elif isinstance(resource.resource, Connection):
        return AsyncConnectionClient.construct_from_other(async_api_client)
    elif isinstance(resource.resource, ToolchainIntegration):
        return AsyncToolchainIntegrationClient.construct_from_other(async_api_client)
    elif isinstance(resource.resource, ModelRoute):
        return AsyncModelRouteClient.construct_from_other(async_api_client)
    elif isinstance(resource.resource, ModelPackaging):
        return AsyncModelPackagingClient.construct_from_other(async_api_client)
    elif isinstance(resource.resource, PackagingIntegration):
        return AsyncPackagingIntegrationClient.construct_from_other(async_api_client)
    else:
        raise InvalidResourceType('{!r} is invalid resource '.format(resource.resource))


def build_resource(declaration: dict) -> OdahuflowCloudResourceUpdatePair:
    """
    Build resource from it's declaration

    :param declaration: declaration of resource
    :return: built resource
    """
    resource_type = declaration.get('kind')
    if resource_type is None:
        raise Exception(f'Kind of object {declaration} must be not null')

    if not isinstance(resource_type, str):
        raise Exception(f'Kind of object {declaration} should be string')

    if resource_type not in TARGET_CLASSES:
        raise Exception(f'Unknown kind of object: \'{resource_type}\'')

    resource = TARGET_CLASSES[resource_type].from_dict(declaration)

    return OdahuflowCloudResourceUpdatePair(
        resource_id=resource.id,
        resource=resource
    )


def parse_stream(data_stream: typing.TextIO) -> OdahuflowCloudResourcesUpdateList:
    """
    Parse YAML/JSON TextIO for Odahuflow resources

    :param data_stream: text stream with yaml/json data
    :raises Exception: if parsing of file is impossible
    :raises ValueError: if invalid Odahuflow resource detected
    :return: parsed resources
    """
    try:
        # Technically YAML is a superset of JSON.
        # So we have to try parse as json first.
        items = json.load(data_stream)

        LOGGER.debug('Successfully parsed as JSON format')
    except json.JSONDecodeError:
        try:
            # Return to the beginning of the stream
            data_stream.seek(0)

            items = tuple(yaml.load_all(data_stream, Loader=yaml.SafeLoader))

            LOGGER.debug('Successfully parsed as YAML format')
        except yaml.YAMLError as e:
            raise Exception('not valid JSON or YAML') from e

    if not isinstance(items, (list, tuple)):
        items = [items]

    result: typing.List[OdahuflowCloudResourceUpdatePair] = []

    for item in items:
        if item is None:
            continue

        if not isinstance(item, dict):
            raise ValueError(f'Invalid Odahuflow resource in file: {item}')

        result.append(build_resource(item))

    return OdahuflowCloudResourcesUpdateList(
        changes=tuple(result)
    )


def parse_resources_dir(path: str) -> typing.List[OdahuflowCloudResourceUpdatePair]:
    """
    Parse dir with YAML/JSON files for Odahuflow resources

    :param path: path to file (local)
    :raises Exception: if parsing of file is impossible
    :raises ValueError: if invalid Odahuflow resource detected
    :return: parsed resources
    """
    if not os.path.exists(path):
        raise FileNotFoundError(f'Resource directory {path} not found')

    if not isdir(path):
        raise FileNotFoundError(f'{path} is not a directory')

    resource_files = [join(path, f) for f in listdir(path) if
                      isfile(join(path, f))]

    entities: typing.List[OdahuflowCloudResourceUpdatePair] = []
    for file in resource_files:
        try:
            LOGGER.debug(f'Parsing the {file} file')
            entities.extend(parse_resources_file(file).changes)
        except Exception:
            LOGGER.exception('parse error')

    return entities


def parse_resources_file(path: str) -> OdahuflowCloudResourcesUpdateList:
    """
    Parse file (YAML/JSON) for Odahuflow resources

    :param path: path to file (local)
    :raises Exception: if parsing of file is impossible
    :raises ValueError: if invalid Odahuflow resource detected
    :return: parsed resources
    """
    if not os.path.exists(path):
        raise FileNotFoundError(f'Resource file \'{path}\' not found')

    with open(path, 'r') as data_stream:
        return parse_stream(data_stream)


def parse_resources_file_with_one_item(path: str) -> OdahuflowCloudResourceUpdatePair:
    """
    Parse file (YAML/JSON) for Odahuflow resource. Raise exception if it is more then one resource

    :param path: path to file (local)
    :raises Exception: if parsing of file is impossible
    :raises Exception: if more then one resource found
    :raises ValueError: if invalid Odahuflow resource detected
    :return: parsed resource
    """
    resources = parse_resources_file(path)
    if len(resources.changes) != 1:
        raise Exception('{!r} should contain 1 item, but {!r} found'.format(path, len(resources)))
    return resources.changes[0]


async def async_apply(updates: OdahuflowCloudResourcesUpdateList,
                      async_api_client: AsyncRemoteAPIClient,
                      is_removal: bool) -> ApplyResult:
    """
    Apply changes on Odahuflow cloud

    :param updates: changes to apply
    :param async_api_client: client to extract connection properties from
    :param is_removal: is it removal?
    :return: result of applying
    """
    created = []
    removed = []
    changed = []
    errors = []

    # Operate over all resources
    for idx, change in enumerate(updates.changes):
        resource_str_identifier = f'#{idx + 1}. {change.resource_id}' if change.resource_id else f'#{idx + 1}'

        LOGGER.debug('Processing resource %r', resource_str_identifier)
        # Build and check client
        try:
            client = build_async_client(change, async_api_client)
        except Exception as general_exception:
            errors.append(Exception(f'Can not get build client for {resource_str_identifier}: {general_exception}'))
            continue

        # Check is resource exist or not
        try:
            await client.get(change.resource_id)
            resource_exist = True
        except WrongHttpStatusCode as http_exception:
            if http_exception.status_code == 404:
                resource_exist = False
            else:
                errors.append(Exception(f'Can not get status of resource '
                                        f'{resource_str_identifier}: {http_exception}'))
                continue
        except Exception as general_exception:
            errors.append(Exception(f'Can not get status of resource '
                                    f'{resource_str_identifier}: {general_exception}'))
            continue

        # Change resource (update/create/delete)
        try:
            # If not removal (creation / update)
            if not is_removal:
                if resource_exist:
                    LOGGER.info('Editing of #%d %s (name: %s)', idx + 1, change.resource, change.resource_id)
                    await client.edit(change.resource)
                    changed.append(change)
                else:
                    LOGGER.info('Creating of #%d %s (name: %s)', idx + 1, change.resource, change.resource_id)
                    await client.create(change.resource)
                    created.append(change)
            # If removal
            else:
                # Only if resource exists on a cluster
                if resource_exist:
                    LOGGER.info('Removing of #%d %s (name: %s)', idx + 1, change.resource, change.resource_id)
                    await client.delete(change.resource_id)
                    removed.append(change)
        except Exception as general_exception:
            errors.append(Exception(f'Can not update resource {resource_str_identifier}: {general_exception}'))
            continue

    return ApplyResult(tuple(created), tuple(removed), tuple(changed), tuple(errors))


def apply(updates: OdahuflowCloudResourcesUpdateList,
          api_client: typing.Union[AsyncRemoteAPIClient, RemoteAPIClient],
          is_removal: bool) -> ApplyResult:
    """
    Apply changes on Odahuflow cloud (wrapper for async_apply). Used for not async client (For backward compatibility)

    :param updates: changes to apply
    :param api_client: client to extract connection properties from
    :param is_removal: is it removal?
    :return:  result of applying
    """
    loop = asyncio.get_event_loop()

    feature = async_apply(updates, api_client, is_removal)

    result = loop.run_until_complete(feature)

    return result
