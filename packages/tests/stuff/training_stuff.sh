#!/usr/bin/env bash
set -e
set -o pipefail

MODEL_NAMES=(simple-model fail counter feedback)
DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" >/dev/null 2>&1 && pwd)"
TRAINED_ARTIFACTS_DIR="${DIR}/trained_artifacts"
ODAHUFLOW_RESOURCES="${DIR}/odahuflow_resources"
TEST_DATA="${DIR}/data"
COMMAND=setup

# Test connection points to the valid gppi archive
TEST_VALID_GPPI_DIR_ID=test-valid-gppi-dir
# Test connection points to the odahu file inside valid gppi archive
TEST_VALID_GPPI_ODAHU_FILE_ID=test-valid-gppi-odahu-file
# Test connection points to the invalid gppi archive
TEST_INVALID_GPPI_DIR_ID=test-invalid-gppi-dir
# Test connection points to the odahu file inside invalid gppi archive
TEST_INVALID_GPPI_ODAHU_FILE_ID=test-invalid-gppi-odahu-file
# Wine data connection
TEST_WINE_CONN_ID=wine

# Test connection to custom output model folder
TEST_CUSTOM_OUTPUT_FOLDER=test-custom-output-folder

TEST_DATA_TI_ID=training-data-helper
# TODO: Remove after implementation of the issue https://github.com/odahuflow-platform/odahuflow/issues/1008
ODAHUFLOW_CONNECTION_DECRYPT_TOKEN=$(jq '.odahuflow_connection_decrypt_token' -r "${CLUSTER_PROFILE}")
EXAMPLES_VERSION=$(jq '.examples_version' -r "${CLUSTER_PROFILE}")

# Cleanups test model packaging from API server, cloud bucket and local filesystem.
# Arguments:
# $1 - Model packaging ID. It must be the same as the directory name where it locates.
function cleanup_pack_model() {
  local mp_id="${1}"

  # Removes the local zip file
  [ -f "${TRAINED_ARTIFACTS_DIR}/${mp_id}.zip" ] && rm "${TRAINED_ARTIFACTS_DIR}/${mp_id}.zip"

  # Removes the model packaging from API service
  odahuflowctl --verbose pack delete --id "${mp_id}" --ignore-not-found
}

# Creates a "trained" artifact from a directory and copy it to the cloud bucket.
# After all, function starts the packaging process.
# Arguments:
# $1 - Model packaging ID. It must be the same as the directory name where it locates.
function pack_model() {
  local mp_id="${1}"

  [ -f "${TRAINED_ARTIFACTS_DIR}/${mp_id}.zip" ] && rm "${TRAINED_ARTIFACTS_DIR}/${mp_id}.zip"

  # Creates trained zip artifact
  cd "${TRAINED_ARTIFACTS_DIR}/${mp_id}"
  zip -r "../${mp_id}.zip" -u ./
  cd -

  # Pushes trained zip artifact to the bucket
  case "${CLOUD_PROVIDER}" in
  aws)
    aws s3 cp "${TRAINED_ARTIFACTS_DIR}/${mp_id}.zip" "s3://${CLUSTER_NAME}-data-store/output/${mp_id}.zip"
    ;;
  azure)
    STORAGE_ACCOUNT=$(az storage account list -g "${CLUSTER_NAME}" --query "[?tags.cluster=='${CLUSTER_NAME}' && tags.purpose=='Odahuflow models storage'].[name]" -otsv)
    az storage blob upload --account-name "${STORAGE_ACCOUNT}" -c "${CLUSTER_NAME}-data-store" \
      -f "${TRAINED_ARTIFACTS_DIR}/${mp_id}.zip" -n "output/${mp_id}.zip"
    ;;
  gcp)
    gsutil cp "${TRAINED_ARTIFACTS_DIR}/${mp_id}.zip" "gs://${CLUSTER_NAME}-data-store/output/${mp_id}.zip"
    ;;
  *)
    echo "Unexpected CLOUD_PROVIDER: ${CLOUD_PROVIDER}"
    usage
    exit 1
    ;;
  esac

  rm "${TRAINED_ARTIFACTS_DIR}/${mp_id}.zip"

  # Repackages trained artifact to a docker image
  odahuflowctl --verbose pack delete --id "${mp_id}" --ignore-not-found
  odahuflowctl --verbose pack create --id "${mp_id}" --artifact-name "${mp_id}.zip" -f "${ODAHUFLOW_RESOURCES}/packaging.odahuflow.yaml"
}

# Waits for all background tasks.
# If one of a task fails, then function fails too.
function wait_all_background_task() {
  local fail_tasks=0

  for job in $(jobs -p); do
    echo "${job} waiting..."
    wait "${job}" || ((fail_tasks = fail_tasks + 1))
  done

  if [[ "$fail_tasks" -ne "0" ]]; then
    echo "Failed $fail_tasks background tasks"
    exit 1
  fi
}

# Create a test data OdahuFlow connection based on models-output connection.
# Arguments:
# $1 - OdahuFlow connection ID, which will be used for new connection
# $2 - OdahuFlow connection uri, which will be used for new connection
function create_test_data_connection() {
  local conn_id="${1}"
  local conn_uri="${2}"
  local conn_file="test-data-connection.yaml"

  # Replaced the uri with the test data directory and added the kind field
  # TODO: Remove the token after implementation of the issue https://github.com/odahuflow-platform/odahuflow/issues/1008
  odahuflowctl conn get --id models-output --decrypted "${ODAHUFLOW_CONNECTION_DECRYPT_TOKEN}" -o json |
    conn_uri="${conn_uri}" jq '.[0].spec.uri = env.conn_uri | .[] | .kind = "Connection"' \
      >"${conn_file}"

  odahuflowctl conn delete --id "${conn_id}" --ignore-not-found
  odahuflowctl conn create -f "${conn_file}" --id "${conn_id}"
  rm "${conn_file}"
}

# Main entrypoint for setup command.
# The function creates the model packaings and the toolchain integrations.
function setup() {
  for mp_id in "${MODEL_NAMES[@]}"; do
    pack_model "${mp_id}" &
  done

  # Create training-data-helper toolchain integration
  jq ".spec.defaultImage = \"${DOCKER_REGISTRY}/odahu-flow-robot-tests:${ODAHUFLOW_VERSION}\"" "${ODAHUFLOW_RESOURCES}/template.training_data_helper_ti.json" >"${ODAHUFLOW_RESOURCES}/training_data_helper_ti.json"
  odahuflowctl ti delete -f "${ODAHUFLOW_RESOURCES}/training_data_helper_ti.json" --ignore-not-found
  odahuflowctl ti create -f "${ODAHUFLOW_RESOURCES}/training_data_helper_ti.json"
  rm "${ODAHUFLOW_RESOURCES}/training_data_helper_ti.json"

  # Download training data for the wine model
  wget -O "${TEST_DATA}/wine-quality.csv" "https://raw.githubusercontent.com/odahu/odahu-examples/${EXAMPLES_VERSION}/mlflow/sklearn/wine/data/wine-quality.csv"

  # Pushes a test data to the bucket and create a file with the connection
  case "${CLOUD_PROVIDER}" in
  aws)
    remote_dir="s3://${CLUSTER_NAME}-data-store/test-data"

    aws s3 cp --recursive "${TEST_DATA}/" "${remote_dir}/data/"
    ;;
  azure)
    remote_dir="${CLUSTER_NAME}-data-store/test-data"
    STORAGE_ACCOUNT=$(az storage account list -g "${CLUSTER_NAME}" --query "[?tags.cluster=='${CLUSTER_NAME}' && tags.purpose=='Odahuflow models storage'].[name]" -otsv)

    az storage blob upload-batch --account-name "${STORAGE_ACCOUNT}" --source "${TEST_DATA}/" \
      --destination "${CLUSTER_NAME}-data-store" --destination-path "test-data/data"
    ;;
  gcp)
    remote_dir="gs://${CLUSTER_NAME}-data-store/test-data"

    gsutil cp -r "${TEST_DATA}/" "${remote_dir}/"
    ;;
  *)
    echo "Unexpected CLOUD_PROVIDER: ${CLOUD_PROVIDER}"
    usage
    exit 1
    ;;
  esac

  # Update test-data connections
  create_test_data_connection "${TEST_VALID_GPPI_ODAHU_FILE_ID}" "${remote_dir}/data/valid_gppi/odahuflow.project.yaml"
  create_test_data_connection "${TEST_VALID_GPPI_DIR_ID}" "${remote_dir}/data/valid_gppi/"
  create_test_data_connection "${TEST_INVALID_GPPI_ODAHU_FILE_ID}" "${remote_dir}/data/invalid_gppi/odahuflow.project.yaml"
  create_test_data_connection "${TEST_INVALID_GPPI_DIR_ID}" "${remote_dir}/data/invalid_gppi/"
  create_test_data_connection "${TEST_CUSTOM_OUTPUT_FOLDER}" "${remote_dir}/data/custom_output/"
  create_test_data_connection "${TEST_WINE_CONN_ID}" "${remote_dir}/data/wine-quality.csv"

  wait_all_background_task
}

# Main entrypoint for cleanup command.
# The function deletes the model packaings and the toolchain integrations.
function cleanup() {
  for mp_id in "${MODEL_NAMES[@]}"; do
    cleanup_pack_model "${mp_id}" &
  done

  odahuflowctl ti delete --id ${TEST_DATA_TI_ID} --ignore-not-found
  odahuflowctl conn delete --id ${TEST_VALID_GPPI_DIR_ID} --ignore-not-found
  odahuflowctl conn delete --id ${TEST_VALID_GPPI_ODAHU_FILE_ID} --ignore-not-found
  odahuflowctl conn delete --id ${TEST_INVALID_GPPI_DIR_ID} --ignore-not-found
  odahuflowctl conn delete --id ${TEST_INVALID_GPPI_ODAHU_FILE_ID} --ignore-not-found
}

# Prints the help message
function usage() {
  echo "Setup or cleanup training stuff for robot tests."
  echo "usage: training_stuff.sh [[setup|cleanup][--models][--help][--verbose]"
}

# The command line arguments parsing
while [ "${1}" != "" ]; do
  case "${1}" in
  setup)
    shift
    COMMAND=setup
    ;;
  cleanup)
    shift
    COMMAND=cleanup
    ;;
  --models)
    mapfile -t MODEL_NAMES <<<"${2}"
    shift 2
    ;;
  --verbose)
    set -x
    shift
    ;;
  --help)
    usage
    exit
    ;;
  *)
    usage
    exit 1
    ;;
  esac
done

# Main programm entrypoint
case "${COMMAND}" in
setup)
  setup
  ;;
cleanup)
  cleanup
  ;;
*)
  echo "Unxpected command: ${COMMAND}"
  usage
  exit 1
  ;;
esac
