#
#      Copyright 2021 EPAM Systems
#
#      Licensed under the Apache License, Version 2.0 (the "License");
#      you may not use this file except in compliance with the License.
#      You may obtain a copy of the License at
#
#          http://www.apache.org/licenses/LICENSE-2.0
#
#      Unless required by applicable law or agreed to in writing, software
#      distributed under the License is distributed on an "AS IS" BASIS,
#      WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
#      See the License for the specific language governing permissions and
#      limitations under the License.


import os
import sys
import pathlib
import json
import logging

logging.basicConfig(level=logging.DEBUG)


class BadConfigurationError(Exception):
    def __init__(self, text):
        self.text = text


def predict():
    """
    Predict function load model that is found at ODAHU_MODEL_PATH env variable.
    Model is multiplier – integer on which every input tensor value will be multiplied
    For example: [1,2,10] => [2,4,20]
    Input must be valid kfserving InputRequest object with extra requirements:
    Only first tensor is processed
    Tensor must be single-dimensional ([2,18,50,11] – correct, [[2,18], [50, 11]] – NOT correct)
    :return:
    """
    model_path = pathlib.Path(os.environ.get("ODAHU_MODEL_PATH"))
    data_path = pathlib.Path(os.environ.get("ODAHU_INPUT_PATH"))
    output_path = pathlib.Path(os.environ.get("ODAHU_OUTPUT_PATH"))

    full_model_path = model_path / "multiplier.txt"
    with open(full_model_path, "r") as f:
        try:
            multiplier = int(f.read())
        except FileNotFoundError:
            raise BadConfigurationError(f"Unable to find model at {full_model_path}")
        except ValueError:
            raise BadConfigurationError(f"multiplier.txt file must be integer")

    try:
        files = os.listdir(data_path)
    except FileNotFoundError:
        raise BadConfigurationError(f"unable to find input data at {data_path}")
    except NotADirectoryError:
        raise BadConfigurationError(f"input data must be directory: {data_path}")

    i = 0
    for f_path in files:
        if f_path.endswith(".json"):
            logging.info(f"processing {f_path}")
            with open(data_path / f_path, "r") as f:
                data = json.load(f)

            try:
                tensor = data["inputs"][0]
                tensor_name = tensor["name"]
                tensor_shape = tensor["shape"]
                tensor_data = tensor["data"]
            except KeyError:
                raise BadConfigurationError("input json files must follow kfserving prediction specification (v2)")

            pathlib.Path(output_path).mkdir(exist_ok=True, parents=True)
            full_o = output_path / f"response{i}.json"
            with open(full_o, "w") as f:
                json.dump({
                    "id": "custom",
                    "outputs": [
                        {
                            "name": tensor_name,
                            "datatype": "INT32",
                            "shape": tensor_shape,
                            "data": [v*multiplier for v in tensor_data]
                        }
                    ]
                }, f)
                logging.info(f"output: {full_o}")
                i += 1


if __name__ == '__main__':

    try:
        predict()
    except BadConfigurationError as e:
        print(f"Script input or configuration is not a correct: {e.text}")
        sys.exit(1)

