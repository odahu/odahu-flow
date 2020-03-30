#!/usr/bin/env bash
set -e

function ReadArguments() {
  if [[ $# == 0 ]]; then
    echo "ERROR: Options not specified! Use -h for help!"
    exit 1
  fi

  while [[ $# -gt 0 ]]; do
    case "$1" in
    -h | --help)
      echo "test_composer.sh - Launch and wait finish of the dags"
      echo -e "Usage: ./test_composer.sh [OPTIONS]\n\noptions:"
      echo -e "--dags\t\tDags for testing, for example: --dags 'airflow-wine,airflow-tensorflow'"
      echo -e "-v  --verbose\t\tverbose mode for debug purposes"
      echo -e "-h  --help\t\tshow brief help"
      exit 0
      ;;
    --dags)
      export TEST_DAG_IDS_RAW="$2"
      shift 2
      ;;
    -v | --verbose)
      export VERBOSE=true
      shift
      ;;
    *)
      echo "ERROR: Unknown option: $1. Use -h for help."
      exit 1
      ;;
    esac
  done

  # Check mandatory parameters
  if [[ ! TEST_DAG_IDS_RAW ]]; then
    echo "ERROR: dags argument must be specified. Use -h for help!"
    exit 1
  else
    IFS=',' read -r -a TEST_DAG_IDS <<< "${TEST_DAG_IDS_RAW}"
    export TEST_DAG_IDS
  fi

  if [[ $VERBOSE == true ]]; then
    set -x
  fi
}

function wait_dags_finish() {
  for i in ${!TEST_DAG_RUN_IDS[@]}; do
    dag_run_id="${TEST_DAG_RUN_IDS[${i}]}"
    dag_id="${TEST_DAG_IDS[${i}]}"

    echo "Wait for the finishing of ${dag_id} and ${dag_run_id} its run"

    while [ true ]; do
      # Extract a dag state from the following output table.
      #------------------------------------------------------------------------------------------------------------------------
      #DAG RUNS
      #------------------------------------------------------------------------------------------------------------------------
      #id  | run_id               | state      | execution_date       | state_date           |
      #152 | manual__2019-12-25T14:53:03+00:00 | success    | 2019-12-25T14:53:03+00:00 | 2019-12-25T14:53:03.236701+00:00 |
      #152 | manual__2019-12-25T13:52:01+00:00 | success    | 2019-12-25T14:53:03+00:00 | 2019-12-25T14:53:03.236701+00:00 |
      state=$(kubectl exec "$POD" -n airflow -it -- airflow list_dag_runs "${dag_id}" | grep -- "${dag_run_id}" | awk '{print $5}')

      case "${state}" in
      "success")
        echo "DAG run ${dag_run_id} finished"
        break
        ;;
      "running")
        echo "DAG run ${dag_run_id} is running. Slepping 30 sec..."
        sleep 30
        ;;
      "failed")
        echo "DAG run ${dag_run_id} failed"
        exit 1
        ;;
      *)
        echo "${state} is unknown state of the ${dag_run_id} DAG"
        exit 1
        ;;
      esac
    done
  done
}

export TEST_DAG_RUN_IDS=()
ReadArguments "$@"
export POD=$(kubectl get pods -l app=airflow -n airflow -o custom-columns=:metadata.name --no-headers | head -n 1)

# Run all test dags
for dag_id in ${TEST_DAG_IDS[@]}; do
  dag_run_id="${dag_id}-ci-$(date +%s)"
  TEST_DAG_RUN_IDS+=("${dag_run_id}")

  echo "Run the ${dag_run_id} of ${dag_id} dag"
  kubectl exec "$POD" -n airflow -it -- airflow trigger_dag -r "${dag_run_id}" "${dag_id}"
done

wait_dags_finish