DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" >/dev/null 2>&1 && pwd)"
LOCAL_TEST_DATA="${DIR}/../resources/artifacts"

# array of image repos for local tests (in removal order)
IMAGE_REPO=(
  gcr.io/or2-msq-epmd-legn-t1iylu/gke-dev04/wine-local-1
  gcr.io/or2-msq-epmd-legn-t1iylu/gke-dev04/wine-artifact-hardcoded-1
  wine-local-1
  wine-artifact-hardcoded-1
  odahu/odahu-flow-mlflow-toolchain
  odahu/odahu-flow-packagers
  odahu/odahu-flow-docker-packager-base
  gcr.io/or2-msq-epmd-legn-t1iylu/odahu/odahu-flow-mlflow-toolchain
  gcr.io/or2-msq-epmd-legn-t1iylu/odahu/odahu-flow-packagers
  gcr.io/or2-msq-epmd-legn-t1iylu/odahu/odahu-flow-docker-packager-base
)

EXAMPLES_VERSION=$(jq '.examples.examples_version' -r "${CLUSTER_PROFILE}")
GIT_REPO_DATA="https://raw.githubusercontent.com/odahu/odahu-examples/${EXAMPLES_VERSION}"

# updates tag for image in specifications for local tests
# Arguments:
# $1 - file name
# $2 - jsonpath to the docker image to be updated
# $3 - tag to replace with
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
  gcloud auth configure-docker --quiet

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
  jq --arg image "${image}" '.[0].spec.defaultImage=$image' "${LOCAL_TEST_DATA}/odahuflow/file/training.json" | sponge "${LOCAL_TEST_DATA}/odahuflow/file/training.json"

  image=$(change_image_tag "${LOCAL_TEST_DATA}/odahuflow/file/training.json" ".[0].spec.defaultImage" "${ti_version}")
  jq --arg image "${image}" '.[0].spec.defaultImage=$image' "${LOCAL_TEST_DATA}/odahuflow/file/training.default.artifact.template.json" | sponge "${LOCAL_TEST_DATA}/odahuflow/file/training.default.artifact.template.json"

  image=$(change_image_tag "${LOCAL_TEST_DATA}/odahuflow/dir/packaging_integration.json" ".spec.defaultImage" "${pi_version}")
  jq --arg image "${image}" '.spec.defaultImage=$image' "${LOCAL_TEST_DATA}/odahuflow/dir/packaging_integration.json" | sponge "${LOCAL_TEST_DATA}/odahuflow/dir/packaging_integration.json"
  image=$(change_image_tag "${LOCAL_TEST_DATA}/odahuflow/dir/packaging_integration.json" ".spec.schema.arguments.properties[1].parameters[2].value" "${pi_version}")
  jq --arg image "${image}" '.spec.schema.arguments.properties[1].parameters[2].value=$image' "${LOCAL_TEST_DATA}/odahuflow/dir/packaging_integration.json" | sponge "${LOCAL_TEST_DATA}/odahuflow/dir/packaging_integration.json"

  image=$(change_image_tag "${LOCAL_TEST_DATA}/odahuflow/file/packaging.json" ".[1].spec.defaultImage" "${pi_version}")
  jq --arg image "${image}" '.[1].spec.defaultImage=$image' "${LOCAL_TEST_DATA}/odahuflow/file/packaging.json" | sponge "${LOCAL_TEST_DATA}/odahuflow/file/packaging.json"
  image=$(change_image_tag "${LOCAL_TEST_DATA}/odahuflow/file/packaging.json" ".[1].spec.schema.arguments.properties[1].parameters[2].value" "${pi_version}")
  jq --arg image "${image}" '.[1].spec.schema.arguments.properties[1].parameters[2].value=$image' "${LOCAL_TEST_DATA}/odahuflow/file/packaging.json" | sponge "${LOCAL_TEST_DATA}/odahuflow/file/packaging.json"
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
  local_setup
  ;;
cleanup)
  local_cleanup
  ;;
*)
  echo "Unexpected command: ${COMMAND}"
  usage
  exit 1
  ;;
esac