#!/usr/bin/env bash
set -e
set -o pipefail

[[ $VERBOSE == true ]] && set -x

MODEL_NAMES=(simple-model fail counter feedback)
DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" >/dev/null 2>&1 && pwd)"
TRAINED_ARTIFACTS_DIR="${DIR}/trained_artifacts"
ODAHUFLOW_RESOURCES="${DIR}/odahuflow_resources"
TEST_DATA="${DIR}/data"
LOCAL_TEST_DATA="${DIR}/../e2e/robot/tests/local/resources/artifacts"
BATCH_TEST_DATA="${DIR}/../e2e/robot/tests/api/resources/batch"
COMMAND=setup

# array of image repos for local tests (in removal order)
IMAGE_REPO=(
  odahu/odahu-flow-mlflow-toolchain
  odahu/odahu-flow-packagers
  odahu/odahu-flow-docker-packager-base
  gcr.io/or2-msq-epmd-legn-t1iylu/odahu/odahu-flow-mlflow-toolchain
  gcr.io/or2-msq-epmd-legn-t1iylu/odahu/odahu-flow-packagers
  gcr.io/or2-msq-epmd-legn-t1iylu/odahu/odahu-flow-docker-packager-base
)

# Test connection points to the valid gppi archive
TEST_VALID_GPPI_DIR_ID=test-valid-gppi-dir
# Test connection points to the odahu file inside valid gppi archive
TEST_VALID_GPPI_ODAHU_FILE_ID=test-valid-gppi-odahu-file
# Wine data connection
TEST_WINE_CONN_ID=wine
# Test connection to custom output model folder
TEST_CUSTOM_OUTPUT_FOLDER=test-custom-output-folder

TEST_DATA_TI_ID=training-data-helper
EXAMPLES_VERSION=$(jq '.examples_version' -r "${CLUSTER_PROFILE}")
CLOUD_PROVIDER="$(jq '.cloud.type' -r "${CLUSTER_PROFILE}")"
BUCKET_NAME="$(jq '.data_bucket' -r "${CLUSTER_PROFILE}")"

GIT_REPO_DATA="https://raw.githubusercontent.com/odahu/odahu-examples/${EXAMPLES_VERSION}"
RCLONE_PROFILE_NAME="robot-tests"

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
  tar -cvzf "../${mp_id}.zip" "."
  cd -

  # Pushes trained zip artifact to the bucket
  copy_to_cluster_bucket "${TRAINED_ARTIFACTS_DIR}/${mp_id}.zip" "${BUCKET_NAME}/output/"

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

function configure_rclone() {
  [[ $VERBOSE == true ]] && set +x
  case "${CLOUD_PROVIDER}" in
  aws)
    local access_key_id secret_access_key region
    access_key_id="$(jq -r .cloud.aws.credentials.AWS_ACCESS_KEY_ID "${CLUSTER_PROFILE}")"
    secret_access_key="$(jq -r .cloud.aws.credentials.AWS_SECRET_ACCESS_KEY "${CLUSTER_PROFILE}")"
    region="$(jq -r .cloud.aws.region "${CLUSTER_PROFILE}")"

    rclone config create "${RCLONE_PROFILE_NAME}" "s3" \
      provider AWS \
      account bucket-owner-full-control \
      server_side_encryption AES256 \
      storage_class STANDARD \
      region "${region}" \
      access_key_id "${access_key_id}" \
      secret_access_key "${secret_access_key}" \
      1>/dev/null

    aws_account_id=$(aws sts get-caller-identity | jq .Account | xargs)
    aws ecr get-login-password --region ${region} | docker login --username AWS --password-stdin ${aws_account_id}.dkr.ecr.${region}.amazonaws.com
    ;;
  azure)
    local sas_url
    sas_url="$(odahuflowctl conn get --id models-output --decrypted -o json | jq -r '.[0].spec.keySecret' | base64 -d)"

    rclone config create "${RCLONE_PROFILE_NAME}" "azureblob" \
      sas_url "${sas_url}" \
      1>/dev/null

    RESOURCE_GROUP="$(jq '.cloud.azure.resource_group' -r "${CLUSTER_PROFILE}")"
    az acr login --name $(az acr list -g ${RESOURCE_GROUP} --query "[0].name" | xargs)
    ;;
  gcp)
    local service_account_credentials
    service_account_credentials="$(jq -r .cloud.gcp.credentials.GOOGLE_CREDENTIALS "${CLUSTER_PROFILE}")"
    rclone config create "${RCLONE_PROFILE_NAME}" "google cloud storage" \
      object_acl projectPrivate \
      bucket_acl projectPrivate \
      bucket_policy_only true \
      service_account_credentials "${service_account_credentials}" \
      1>/dev/null
    ;;
  *)
    echo "Unexpected CLOUD_PROVIDER: ${CLOUD_PROVIDER}"
    usage
    exit 1
    ;;
  esac
  [[ $VERBOSE == true ]] && set -x
}

# Copy local directory or file to a bucket
# $1 - local directory or file
# $2 - bucket directory or file
function copy_to_cluster_bucket() {
  local source="${1}"
  local target="${2}"

  rclone copy "${source}" "${RCLONE_PROFILE_NAME}:${target}"
}

# Create a test data OdahuFlow connection based on models-output connection.
# Arguments:
# $1 - OdahuFlow connection ID, which will be used for new connection
# $2 - OdahuFlow connection uri, which will be used for new connection
function create_test_data_connection() {
  case "${CLOUD_PROVIDER}" in
  aws)
    remote_dir="s3://${BUCKET_NAME}"
    ;;
  azure)
    remote_dir="${BUCKET_NAME}"
    ;;
  gcp)
    remote_dir="gs://${BUCKET_NAME}"
    ;;
  *)
    echo "Unexpected CLOUD_PROVIDER: ${CLOUD_PROVIDER}"
    usage
    exit 1
    ;;
  esac

  local conn_id="${1}"
  local conn_uri="${remote_dir}/${2}"
  local conn_file="test-data-connection.yaml"

  # Replaced the uri with the test data directory and added the kind field
  odahuflowctl conn get --id models-output --decrypted -o json |
    conn_uri="${conn_uri}" jq '.[0].spec.uri = env.conn_uri | .[] | .kind = "Connection"' \
      >"${conn_file}"

  odahuflowctl conn delete --id "${conn_id}" --ignore-not-found
  odahuflowctl conn create -f "${conn_file}" --id "${conn_id}"
  rm "${conn_file}"
}

# Upload test dags from odahu-examples repository to the cluster dags object storage bucket
function upload_test_dags() {
  git_url="https://github.com/odahu/odahu-examples.git"
  dag_dirs=("mlflow/sklearn/wine/airflow")
  tmp_odahu_example_dir=$(mktemp -d -t upload-test-dags-XXXXXXXXXX)

  git clone --branch "${EXAMPLES_VERSION}" "${git_url}" "${tmp_odahu_example_dir}"

  local dag_bucket dag_bucket_path
  dag_bucket="$(jq -r ".airflow.dag_bucket |= if . == null or . == \"\" then \"$BUCKET_NAME\" else . end | .airflow.dag_bucket" "${CLUSTER_PROFILE}")"
  dag_bucket_path="$(jq -r '.airflow.dag_bucket_path |= if . == null or . == "" then "/dags" else . end | .airflow.dag_bucket_path' "${CLUSTER_PROFILE}")"

  for dag_dir in "${dag_dirs[@]}"; do
    copy_to_cluster_bucket "${tmp_odahu_example_dir}/${dag_dir}/" "${dag_bucket}${dag_bucket_path}/${dag_dir}/"
  done

  rm -rf "${tmp_odahu_example_dir}"
}

function upload_test_algorithm() {
  git_url="https://github.com/odahu/odahu-examples.git"
  algorithm_code_path=("mlflow/sklearn/wine/MLproject")
  tmp_odahu_example_dir=$(mktemp -d -t upload-test-algorithm-XXXXXXXXXX)

  git clone --branch "${EXAMPLES_VERSION}" "${git_url}" "${tmp_odahu_example_dir}"

  copy_to_cluster_bucket "${tmp_odahu_example_dir}/${algorithm_code_path}" "${BUCKET_NAME}/test_algorithm/wine/"

  rm -rf "${tmp_odahu_example_dir}"
}

# Prepare for batch e2e test
function setup_batch_examples() {
  local git_url="https://github.com/odahu/odahu-examples.git"
  local dir="batch-inference"
  local tmp_odahu_example_dir=$(mktemp -d -t examples-XXXXXXXXXX)

  git clone --branch "${EXAMPLES_VERSION}" "${git_url}" "${tmp_odahu_example_dir}"

  # Build and predictor image
  docker build ${tmp_odahu_example_dir}/batch-inference/predictor -t ${DOCKER_REGISTRY}/odahu-flow-batch-predictor-test:${ODAHUFLOW_VERSION}
  docker push ${DOCKER_REGISTRY}/odahu-flow-batch-predictor-test:${ODAHUFLOW_VERSION}

  # Build and predictor image with embedded model
  docker build ${tmp_odahu_example_dir}/batch-inference/predictor_embedded -t ${DOCKER_REGISTRY}/odahu-flow-batch-predictor-test-embedded:${ODAHUFLOW_VERSION}
  docker push ${DOCKER_REGISTRY}/odahu-flow-batch-predictor-test-embedded:${ODAHUFLOW_VERSION}

  # Prepare test data by replacing image in spec of service and copying job manifest
  yq w ${tmp_odahu_example_dir}/batch-inference/manifests/inferenceservice.yaml \
    'spec.image' ${DOCKER_REGISTRY}/odahu-flow-batch-predictor-test:${ODAHUFLOW_VERSION} > "${DIR}/../e2e/robot/tests/api/resources/batch/inferenceservice.yaml"
  cp ${tmp_odahu_example_dir}/batch-inference/manifests/inferencejob.yaml "${DIR}/../e2e/robot/tests/api/resources/batch/inferencejob.yaml"

  yq w ${tmp_odahu_example_dir}/batch-inference/manifests/inferenceservice-packed.yaml \
    'spec.image' ${DOCKER_REGISTRY}/odahu-flow-batch-predictor-test:${ODAHUFLOW_VERSION} > "${DIR}/../e2e/robot/tests/api/resources/batch/inferenceservice-packed.yaml"
  cp ${tmp_odahu_example_dir}/batch-inference/manifests/inferencejob-packed.yaml "${DIR}/../e2e/robot/tests/api/resources/batch/inferencejob-packed.yaml"

  # embedded
  yq w ${tmp_odahu_example_dir}/batch-inference/manifests/inferenceservice-embedded.yaml \
    'spec.image' ${DOCKER_REGISTRY}/odahu-flow-batch-predictor-test-embedded:${ODAHUFLOW_VERSION} > "${DIR}/../e2e/robot/tests/api/resources/batch/inferenceservice-embedded.yaml"
  cp ${tmp_odahu_example_dir}/batch-inference/manifests/inferencejob-embedded.yaml "${DIR}/../e2e/robot/tests/api/resources/batch/inferencejob-embedded.yaml"

  cp -r ${tmp_odahu_example_dir}/batch-inference/output "${DIR}/../e2e/robot/tests/api/resources/batch/output/"
  # Upload model and input data to object storage
  copy_to_cluster_bucket ${tmp_odahu_example_dir}/batch-inference/input "${BUCKET_NAME}/test-data/batch_job_data/input"
  copy_to_cluster_bucket ${tmp_odahu_example_dir}/batch-inference/model "${BUCKET_NAME}/output/test-data/batch_job_data/model"
  copy_to_cluster_bucket ${tmp_odahu_example_dir}/batch-inference/model.tar.gz "${BUCKET_NAME}/output/test-data/batch_job_data/"
  # Clean tmp dir
  rm -rf "${tmp_odahu_example_dir}"
}

# updates tag for image in specifications for local tests
function change_image_tag() {
  local file_name=$1
  local json_path=$2
  local tag=$3

  image=$(jq -r "${json_path}" "${file_name}" | cut -d ':' -f 1)
  echo ${image}:${tag}
}

# Upload files for local training and packaging
function local_setup() {
  # download example files
  wget -O "${LOCAL_TEST_DATA}/../request.json" "${GIT_REPO_DATA}/mlflow/sklearn/wine/odahuflow/request.json"
  wget -O "${LOCAL_TEST_DATA}/MLproject" "${GIT_REPO_DATA}/mlflow/sklearn/wine/MLproject"
  wget -O "${LOCAL_TEST_DATA}/train.py" "${GIT_REPO_DATA}/mlflow/sklearn/wine/train.py"
  wget -O "${LOCAL_TEST_DATA}/conda.yaml" "${GIT_REPO_DATA}/mlflow/sklearn/wine/conda.yaml"
  wget -O "${LOCAL_TEST_DATA}/wine-quality.csv" "${GIT_REPO_DATA}/mlflow/sklearn/wine/data/wine-quality.csv"

  # configure Docker: https://cloud.google.com/container-registry/docs/advanced-authentication#gcloud-helper
  gcloud auth configure-docker

  if [ ! -x "$(command -v sponge)" ]; then
    printf "\nPlease install moreutils or sponge package to setup robot tests\n"
    exit 1
  fi

  CMDbase64="base64 --wrap=0"
  if [ -x "$(command -v brew)" ]; then
    CMDbase64=base64
  fi

  echo ${CMDbase64}

  # update specifications
  ## docker-pull target
  local docker_uri="$(jq -r .docker_repo "${CLUSTER_PROFILE}")"
  local docker_username="$(jq -r .docker_username "${CLUSTER_PROFILE}")"
  local docker_password="$(jq -r .docker_password "${CLUSTER_PROFILE}" | tr -d "\n" | $CMDbase64)"

  jq --arg uri "${docker_uri}" --arg username "${docker_username}" --arg password "${docker_password}" \
    '.spec.uri=$uri | .spec.username=$username | .spec.password=$password' "${LOCAL_TEST_DATA}/odahuflow/dir/docker-pull-target.json" | jq "." | sponge "${LOCAL_TEST_DATA}/odahuflow/dir/docker-pull-target.json"

  ## docker image tags
  local ti_version="$(jq -r .mlflow_toolchain_version "${CLUSTER_PROFILE}")"
  local pi_version="$(jq -r .packager_version "${CLUSTER_PROFILE}")"

  image=$(change_image_tag "${LOCAL_TEST_DATA}/odahuflow/dir/toolchain_integration.json" ".spec.defaultImage" "${ti_version}")
  jq --arg image "${image}" '.spec.defaultImage=$image' "${LOCAL_TEST_DATA}/odahuflow/dir/toolchain_integration.json" | sponge "${LOCAL_TEST_DATA}/odahuflow/dir/toolchain_integration.json"

  image=$(change_image_tag "${LOCAL_TEST_DATA}/odahuflow/file/training.json" ".[0].spec.defaultImage" "${ti_version}")
  jq --arg image "$image" '.[0].spec.defaultImage=$image' "${LOCAL_TEST_DATA}/odahuflow/file/training.json" | sponge "${LOCAL_TEST_DATA}/odahuflow/file/training.json"

  image=$(change_image_tag "${LOCAL_TEST_DATA}/odahuflow/file/training.json" ".[0].spec.defaultImage" "${ti_version}")
  jq --arg image "$image" '.[0].spec.defaultImage=$image' "${LOCAL_TEST_DATA}/odahuflow/file/training.default.artifact.template.json" | sponge "${LOCAL_TEST_DATA}/odahuflow/file/training.default.artifact.template.json"

  image=$(change_image_tag "${LOCAL_TEST_DATA}/odahuflow/dir/packaging_integration.json" ".spec.defaultImage" "${pi_version}")
  jq --arg image "$image" '.spec.defaultImage=$image' "${LOCAL_TEST_DATA}/odahuflow/dir/packaging_integration.json" | sponge "${LOCAL_TEST_DATA}/odahuflow/dir/packaging_integration.json"
  image=$(change_image_tag "${LOCAL_TEST_DATA}/odahuflow/dir/packaging_integration.json" ".spec.schema.arguments.properties[1].parameters[2].value" "${pi_version}")
  jq --arg image "$image" '.spec.schema.arguments.properties[1].parameters[2].value=$image' "${LOCAL_TEST_DATA}/odahuflow/dir/packaging_integration.json" | sponge "${LOCAL_TEST_DATA}/odahuflow/dir/packaging_integration.json"

  image=$(change_image_tag "${LOCAL_TEST_DATA}/odahuflow/file/packaging.json" ".[1].spec.defaultImage" "${pi_version}")
  jq --arg image "$image" '.[1].spec.defaultImage=$image' "${LOCAL_TEST_DATA}/odahuflow/file/packaging.json" | sponge "${LOCAL_TEST_DATA}/odahuflow/file/packaging.json"
  image=$(change_image_tag "${LOCAL_TEST_DATA}/odahuflow/file/packaging.json" ".[1].spec.schema.arguments.properties[1].parameters[2].value" "${pi_version}")
  jq --arg image "$image" '.[1].spec.schema.arguments.properties[1].parameters[2].value=$image' "${LOCAL_TEST_DATA}/odahuflow/file/packaging.json" | sponge "${LOCAL_TEST_DATA}/odahuflow/file/packaging.json"
}

# remove packager images locally
function local_cleanup() {
  echo "IMAGE_REPOS: ${IMAGE_REPO[*]}"
  for repo in ${IMAGE_REPO[@]}; do
    list_images=$(docker images -aq ${repo})
    for image in ${list_images}; do
      bash ${DIR}/docker-remove-image.sh ${image}
    done
  done
}

# Main entrypoint for setup command.
# The function creates the model packagings and the toolchain integrations.
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
  wget -O "${TEST_DATA}/wine-quality.csv" "${GIT_REPO_DATA}/mlflow/sklearn/wine/data/wine-quality.csv"

  # Pushes a test data to the bucket and create a file with the connection
  copy_to_cluster_bucket "${TEST_DATA}/" "${BUCKET_NAME}/test-data/data/"

  # Update test-data connections
  create_test_data_connection "${TEST_VALID_GPPI_ODAHU_FILE_ID}" "test-data/data/valid_gppi/odahuflow.project.yaml"
  create_test_data_connection "${TEST_VALID_GPPI_DIR_ID}" "test-data/data/valid_gppi/"
  create_test_data_connection "${TEST_CUSTOM_OUTPUT_FOLDER}" "test-data/data/custom_output/"
  create_test_data_connection "${TEST_WINE_CONN_ID}" "test-data/data/wine-quality.csv"

  upload_test_dags

  wait_all_background_task

  local_setup
  setup_batch_examples
  upload_test_algorithm
}

# Main entrypoint for cleanup command.
# The function deletes the model packagings and the toolchain integrations and the pulled locally packager image
function cleanup() {
  for mp_id in "${MODEL_NAMES[@]}"; do
    cleanup_pack_model "${mp_id}" &
  done

  odahuflowctl ti delete --id ${TEST_DATA_TI_ID} --ignore-not-found
  odahuflowctl conn delete --id ${TEST_VALID_GPPI_DIR_ID} --ignore-not-found
  odahuflowctl conn delete --id ${TEST_VALID_GPPI_ODAHU_FILE_ID} --ignore-not-found

  local_cleanup
}

# Prints the help message
function usage() {
  echo "Setup or cleanup training stuff for robot tests."
  echo "usage: ${0} [[setup|cleanup][--models][--help][--verbose]"
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

configure_rclone

# Main programm entrypoint
case "${COMMAND}" in
setup)
  setup
  ;;
cleanup)
  cleanup
  ;;
*)
  echo "Unexpected command: ${COMMAND}"
  usage
  exit 1
  ;;
esac
