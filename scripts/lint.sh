#!/bin/bash
set -e
set -o pipefail

PYLINT_FOLDER="target/pylint"

function pylint_cmd() {
    package_dir="${1}"
    output_name="${2}"
    additional_pylint_args="${3}"

    pylint --output-format=parseable \
           "${additional_pylint_args}" \
           --reports=no "packages/${package_dir}" 2>&1 | tee "${PYLINT_FOLDER}/odahu-flow-${output_name}.log" &
}

rm -rf "${PYLINT_FOLDER}"
mkdir -p "${PYLINT_FOLDER}"

# Ignoring models package because it contains autogenerated python code
pylint_cmd sdk/odahuflow sdk "--ignore=models"
pylint_cmd cli/odahuflow cli
pylint_cmd robot/odahuflow robot


rm -rf "${PYDOCSTYLE_FOLDER}"
mkdir -p "${PYDOCSTYLE_FOLDER}"

FAIL=0
# Wait all background linters
for job in $(jobs -p)
do
    echo "${job}"
    echo "waiting..."
    wait "${job}" || ((FAIL = FAIL + 1))
done

cat ${PYLINT_FOLDER}/*.log > "${PYLINT_FOLDER}/odahuflow.log"

echo "You can find the result of linting here: ${PYLINT_FOLDER}"

if [[ "$FAIL" -ne "0" ]];
then
    echo "Failed $FAIL linters"
    exit 1
fi
