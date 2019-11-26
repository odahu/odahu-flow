#!/usr/bin/env bash
set -e

PROJECTS="packages/sdk packages/cli packages/robot"
ROOT_DIR="$(pwd)"

for project in ${PROJECTS}
do
    echo "Update dependencies in ${ROOT_DIR}/${project}"
    cd "${ROOT_DIR}/${project}"

    pipenv update

    ${ROOT_DIR}/scripts/convert_pipenv_to_requirements.py $(pwd)
done
