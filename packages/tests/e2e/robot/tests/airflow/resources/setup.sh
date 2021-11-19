#!/usr/bin/env bash
set -e
set -x
set -o pipefail

EXAMPLES_VERSION=$(jq '.examples.examples_version' -r "${CLUSTER_PROFILE}")
BUCKET_NAME="$(jq '.data_bucket' -r "${CLUSTER_PROFILE}")"
# the same as RCLONE_PROFILE_NAME variable in "packages/tests/stuff/training_stuff.sh" file
RCLONE_PROFILE_NAME="robot-tests"

# Copy local directory or file to a bucket
# $1 - local directory or file
# $2 - bucket directory or file
function copy_to_cluster_bucket() {
  local source="${1}"
  local target="${2}"

  rclone copy "${source}" "${RCLONE_PROFILE_NAME}:${target}"
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

upload_test_dags
