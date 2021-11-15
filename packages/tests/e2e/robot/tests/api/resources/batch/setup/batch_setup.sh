#!/usr/bin/env bash
set -e
set -x
set -o pipefail

DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" >/dev/null 2>&1 && pwd)"

BATCH_TEST_DATA="${DIR}/.."

EXAMPLES_VERSION=$(jq '.examples.examples_version' -r "${CLUSTER_PROFILE}")
BUCKET_NAME="$(jq '.data_bucket' -r "${CLUSTER_PROFILE}")"
RCLONE_PROFILE_NAME="robot-tests"

# Copy local directory or file to a bucket
# $1 - local directory or file
# $2 - bucket directory or file
function copy_to_cluster_bucket() {
  local source="${1}"
  local target="${2}"

  rclone copy "${source}" "${RCLONE_PROFILE_NAME}:${target}"
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
    'spec.image' ${DOCKER_REGISTRY}/odahu-flow-batch-predictor-test:${ODAHUFLOW_VERSION} > "${BATCH_TEST_DATA}/inferenceservice.yaml"
  cp ${tmp_odahu_example_dir}/batch-inference/manifests/inferencejob.yaml "${BATCH_TEST_DATA}/inferencejob.yaml"

  yq w ${tmp_odahu_example_dir}/batch-inference/manifests/inferenceservice-packed.yaml \
    'spec.image' ${DOCKER_REGISTRY}/odahu-flow-batch-predictor-test:${ODAHUFLOW_VERSION} > "${BATCH_TEST_DATA}/inferenceservice-packed.yaml"
  cp ${tmp_odahu_example_dir}/batch-inference/manifests/inferencejob-packed.yaml "${BATCH_TEST_DATA}/inferencejob-packed.yaml"

  # embedded
  yq w ${tmp_odahu_example_dir}/batch-inference/manifests/inferenceservice-embedded.yaml \
    'spec.image' ${DOCKER_REGISTRY}/odahu-flow-batch-predictor-test-embedded:${ODAHUFLOW_VERSION} > "${BATCH_TEST_DATA}/inferenceservice-embedded.yaml"
  cp ${tmp_odahu_example_dir}/batch-inference/manifests/inferencejob-embedded.yaml "${BATCH_TEST_DATA}/inferencejob-embedded.yaml"

  cp -r ${tmp_odahu_example_dir}/batch-inference/output "${BATCH_TEST_DATA}/output/"
  # Upload model and input data to object storage
  copy_to_cluster_bucket ${tmp_odahu_example_dir}/batch-inference/input "${BUCKET_NAME}/test-data/batch_job_data/input"
  copy_to_cluster_bucket ${tmp_odahu_example_dir}/batch-inference/model "${BUCKET_NAME}/output/test-data/batch_job_data/model"
  copy_to_cluster_bucket ${tmp_odahu_example_dir}/batch-inference/model.tar.gz "${BUCKET_NAME}/output/test-data/batch_job_data/"
  # Clean tmp dir
  rm -rf "${tmp_odahu_example_dir}"
}

# The command line arguments parsing
while [ "${1}" != "" ]; do
  case "${1}" in
  --docker-registry)
    DOCKER_REGISTRY=${2}
    shift 2
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

setup_batch_examples