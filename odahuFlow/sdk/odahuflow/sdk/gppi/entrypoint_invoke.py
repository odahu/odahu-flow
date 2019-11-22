import argparse
import importlib
import logging
import os
import pickle
import sys
from typing import Any

MODEL_LOCATION = os.environ.get('MODEL_LOCATION', '.')
sys.path.append(MODEL_LOCATION)

# Optional. Examples of input and output pandas DataFrames
MODEL_INPUT_SAMPLE_FILE = os.path.join(MODEL_LOCATION, 'head_input.pkl')
MODEL_OUTPUT_SAMPLE_FILE = os.path.join(MODEL_LOCATION, 'head_output.pkl')

_logger = logging.getLogger(__name__)


def get_entrypoint():
    return importlib.import_module(args.entrypoint)


def load_pickle_artifact(name: str) -> Any:
    if os.path.exists(name):
        with open(name, 'rb') as f:
            return pickle.load(f)


def main():

    entrypoint = get_entrypoint()

    _logger.info('GPPI entrypoint module successfully imported')

    entrypoint.init()
    _logger.info('GPPI entrypoint.init() – successfully tested')

    entrypoint.info()
    _logger.info('GPPI entrypoint.info() – successfully tested')

    input_ = load_pickle_artifact(MODEL_INPUT_SAMPLE_FILE)
    if input_ is not None:
        entrypoint.predict_on_matrix(input_)
        _logger.info('GPPI entrypoint.predict_on_matrix(...) successfully tested')
    else:
        _logger.warning(f'GPPI {MODEL_INPUT_SAMPLE_FILE} is not found in model directory. '
                        'Testing of a predicition API is impossible. '
                        'OpenAPI schema will not be generated for the model.')


if __name__ == '__main__':

    logging.basicConfig(level=logging.INFO)

    parser = argparse.ArgumentParser()
    parser.add_argument('entrypoint', help='Name of entrypoint GPPI module inside $MODEL_LOCATION path')
    args = parser.parse_args()

    try:
        main()
    except ImportError as e:
        raise ImportError(
            'ImportError usually happens when you have not packed all required dependencies for you model.'
            'Please see your Training Toolchain documentation to get more info about packing your '
            'model script dependencies'
        ) from e
