#!/usr/bin/env bash
set -e
set -x
set -o pipefail

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

function upload_test_algorithm() {
  git_url="https://github.com/odahu/odahu-examples.git"
  algorithm_code_path=("mlflow/sklearn/wine/")
  tmp_odahu_example_dir=$(mktemp -d -t upload-test-algorithm-XXXXXXXXXX)

  git clone --branch "${EXAMPLES_VERSION}" "${git_url}" "${tmp_odahu_example_dir}"
  copy_to_cluster_bucket "${tmp_odahu_example_dir}/${algorithm_code_path}" "${BUCKET_NAME}/test_algorithm/wine/"
  rm -rf "${tmp_odahu_example_dir}"
}

upload_test_algorithm