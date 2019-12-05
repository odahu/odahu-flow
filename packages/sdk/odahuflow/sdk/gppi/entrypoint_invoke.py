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
This package must has no requirements except Python3.6
Package provide CLI that can be invoked in any gppi environment without installing extra deps.
The main purpose of package – interact with GPPI folder via CLI in the gppi environment.
This package not provide API to creating and installing Conda Envs by request. You need to create and install
all required GPPI requirements by yourself. (Or check `odahuflowctl gppi` command that can fulfill env requirements
before invoking GPPI)
"""

import argparse
import functools
import importlib
import json
import logging
import os
import pickle
import sys
from contextlib import contextmanager
from typing import Any, List, Union, Dict

_logger = logging.getLogger(__name__)

ODAHUFLOW_MODEL_LOCATION_ENV_VAR = 'MODEL_LOCATION'
MODEL_ENTRYPOINT_ENV = 'MODEL_ENTRYPOINT'

args = None


class CliError(Exception):
    pass


def build_error_response(message):
    raise CliError(message)


@functools.lru_cache()
def get_json_output_serializer():
    entrypoint = get_entrypoint()
    if hasattr(entrypoint, 'get_output_json_serializer'):
        return entrypoint.get_output_json_serializer()
    else:
        return None


def generate_input_props(input_schema: List[Dict[str, Union[Union[str, None, bool], Any]]]) -> Dict[str, Any]:
    """
    Generate input model properties in OpenAPI format
    :param input_schema: Input schema
    :return:
    """
    examples: List[Any] = []
    columns: List[str] = []
    for prop in input_schema:
        columns.append(prop['name'])
        examples.append(prop['example'])

    return {
        "columns": {
            "example": columns,
            "items": {
                "type": "string"
            },
            "type": "array"
        },
        "data": {
            "items": {
                "items": {
                    "type": "number"
                },
                "type": "array"
            },
            "type": "array",
            "example": [examples],
        }
    }


def generate_output_props(output_schema: List[Dict[str, Union[Union[str, None, bool], Any]]]) -> Dict[str, Any]:
    """
    Generate input model properties in OpenAPI format
    :param output_schema:
    :return:
    """
    examples: List[Any] = []
    columns: List[str] = []
    for prop in output_schema:
        columns.append(prop['name'])
        examples.append(prop['example'])

    return {
        "prediction": {
            "example": [examples],
            "items": {
                "type": "number"
            },
            "type": "array",
        },
        "columns": {
            "example": columns,
            "items": {
                "type": "string"
            },
            "type": "array"
        }
    }


def handle_prediction_on_matrix(parsed_data):
    matrix = parsed_data.get('data')
    columns = parsed_data.get('columns', None)

    if not matrix:
        return build_error_response('Matrix is not provided')

    entrypoint = get_entrypoint()

    try:
        prediction, columns = entrypoint.predict_on_matrix(matrix, provided_columns_names=columns)
    except Exception as predict_exception:
        return build_error_response(f'Exception during prediction: {predict_exception}')

    response = {
        'prediction': prediction,
        'columns': columns
    }

    return response


_entrypoint_module = None


def get_entrypoint():
    """
    Return model entrypoint module using required passed argument
    :return:
    """
    if _entrypoint_module is None:
        return importlib.import_module(args.entrypoint)
    else:
        return _entrypoint_module


def get_model_location() -> str:
    """
    Return model location path using environment variable
    :return:
    """
    return os.environ.get(ODAHUFLOW_MODEL_LOCATION_ENV_VAR, ".")


def get_model_input_sample() -> str:
    """
    Return input sample path
    :return:
    """
    return os.path.join(get_model_location(), 'head_input.pkl')


def get_model_output_sample() -> str:
    """
    Return input sample path
    :return:
    """
    return os.path.join(get_model_location(), 'head_output.pkl')


def load_pickle_artifact(name: str) -> Any:
    if os.path.exists(name):
        with open(name, 'rb') as f:
            return pickle.load(f)


def self_check():

    entrypoint = get_entrypoint()

    _logger.info('GPPI entrypoint module successfully imported')

    entrypoint.init()
    _logger.info('GPPI entrypoint.init() – successfully tested')

    entrypoint.info()
    _logger.info('GPPI entrypoint.info() – successfully tested')

    input_ = load_pickle_artifact(get_model_input_sample())
    if input_ is not None:
        entrypoint.predict_on_matrix(input_)
        _logger.info('GPPI entrypoint.predict_on_matrix(...) successfully tested')
    else:
        _logger.warning(f'GPPI {get_model_input_sample()} is not found in model directory. '
                        'Testing of a predicition API is impossible. '
                        'OpenAPI schema will not be generated for the model.')

    print('Self check is successful')


def self_check_caller(args):
    return self_check()


def predict(input_file: str, output_dir: str, output_file_name: str):

    prediction_mode = get_entrypoint().init()

    with open(input_file) as fp:
        parsed_data = json.load(fp)

    if prediction_mode == 'matrix':
        response = handle_prediction_on_matrix(parsed_data)
    else:
        return build_error_response(f'Unknown model\'s return type: {prediction_mode}')

    res_fp = os.path.join(output_dir, output_file_name)
    with open(res_fp, 'w') as fp:
        json.dump(response, fp, cls=get_json_output_serializer())

    print(f'Prediction is successful. Result file: {res_fp}')


def predict_caller(args):
    predict(args.input_file, args.output_dir, args.output_file_name)


def info():
    entrypoint = get_entrypoint()
    input_schema, output_schema = entrypoint.info()
    input_properties = generate_input_props(input_schema)
    output_properties = generate_output_props(output_schema)
    print('Input schema:')
    print(json.dumps(input_properties, indent=4))
    print('Output schema:')
    print(json.dumps(output_properties, indent=4))


def info_caller(args):
    info()


def _configure_logging(verbose: bool):
    if verbose:
        logging.basicConfig(level=logging.DEBUG)


def _configure_arg_parser() -> argparse.ArgumentParser:

    _parser = argparse.ArgumentParser(description="""
    Provide CLI to invoke GPPI entrypoint API.
    This module has only stdlib dependencies so it could be executed from any GPPI environment
    """)
    _parser.add_argument('-v', help='Verbosity logs', action='store_true')
    _parser.add_argument('--entrypoint', help='Name of entrypoint GPPI module')
    _parser.add_argument('--model', help="""
    Override $MODEL_LOCATION environment variable before entrypoint import.
    Clean overridden value after script execution (success or fail).
    """)
    _parser.set_defaults(func=lambda args_: _parser.print_help())

    subparsers = _parser.add_subparsers()

    parser_self_check = subparsers.add_parser('self_check', help='Self check GPPI correctness')
    parser_self_check.set_defaults(func=self_check_caller)

    parser_predict = subparsers.add_parser('predict', help='Make predictions using GPPI model')
    parser_predict.add_argument('input_file', help='JSON file with input data for prediction')
    parser_predict.add_argument('output_dir', help='Directory where JSON with predictions will be saved')
    parser_predict.add_argument('--output_file_name', help='JSON filename with predictions', default='results.json')
    parser_predict.set_defaults(func=predict_caller)

    parser_info = subparsers.add_parser('info', help='Show model input/output data schema')
    parser_info.set_defaults(func=info_caller)

    return _parser


@contextmanager
def model_location(model):
    """
    Set ODAHUFLOW_MODEL_LOCATION_ENV_VAR to `model` if `model` is not None
    Add ODAHUFLOW_MODEL_LOCATION_ENV_VAR to sys.path

    Clean state in context manager exit
    :param model:
    :return:
    """
    model_location_for_use = original_model_location = os.environ.get(ODAHUFLOW_MODEL_LOCATION_ENV_VAR, "")
    if original_model_location:
        _logger.debug(f'${ODAHUFLOW_MODEL_LOCATION_ENV_VAR} env var is set in a system '
                      f'(${ODAHUFLOW_MODEL_LOCATION_ENV_VAR}={original_model_location})')
    else:
        _logger.debug(f'${ODAHUFLOW_MODEL_LOCATION_ENV_VAR} env is not set in a system')

    if model:
        _logger.debug(f'--model option is passed. ${ODAHUFLOW_MODEL_LOCATION_ENV_VAR} will be replaced '
                      f'with {model}')
        model_location_for_use = os.environ[ODAHUFLOW_MODEL_LOCATION_ENV_VAR] = model

    if not model_location_for_use:
        raise RuntimeError(f'Either ${ODAHUFLOW_MODEL_LOCATION_ENV_VAR} env var or --model option MUST be specified')

    sys.path.append(model_location_for_use)
    _logger.debug(f'{model_location_for_use} is added to sys.path')

    try:
        yield
    finally:
        if model:
            os.environ[ODAHUFLOW_MODEL_LOCATION_ENV_VAR] = original_model_location
            _logger.debug(f'{ODAHUFLOW_MODEL_LOCATION_ENV_VAR} is set to original value={original_model_location}')
        if model_location_for_use:
            sys.path.remove(model_location_for_use)
            _logger.debug(f'{model_location_for_use} is removed from sys.path')


def main():

    global args

    parser = _configure_arg_parser()

    args = parser.parse_args()

    if not args.entrypoint:
        args.entrypoint = os.environ.get(MODEL_ENTRYPOINT_ENV)
    if not args.entrypoint:
        raise RuntimeError(f'Either ${MODEL_ENTRYPOINT_ENV} env var or --entrypoint option '
                           f'MUST be specified')

    _configure_logging(args.v)

    with model_location(args.model):
        try:
            args.func(args)
        except ImportError as exc_info:
            raise ImportError(
                'ImportError usually happens when you have not packed all required dependencies for you model.'
                'Please see your Training Toolchain documentation to get more info about packing your '
                'model script dependencies'
            ) from exc_info


if __name__ == '__main__':
    main()
