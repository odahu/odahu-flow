import functools

import kubernetes
import kubernetes.client
import kubernetes.client.rest

ODAHU_DEPL_LABEL = "odahu.org/deploymentID"

version = "v1alpha1"
group = "odahuflow.odahu.org"

training_namespace = "odahu-flow-training"
training_plural = "modeltrainings"

packaging_namespace = "odahu-flow-packaging"
packaging_plural = "modelpackagings"

deployment_namespace = "odahu-flow-deployment"
deployment_plural = "modeldeployments"


def _api_exc_print(e: kubernetes.client.rest.ApiException) -> str:
    return f"ApiException: status={e.status}, reason={e.reason}, body: {e.body}, headers={e.headers}"


def print_report(func):
    @functools.wraps(func)
    def wrapper(self, name):
        report = func(self, name)
        print(report)
        return report

    return wrapper


class OdahuKubeReporter:
    """
    Print information about ODAHU CRD for training, pack, deployment
    and underlying k8s resources: deployment, pods, etc..
    Useful to log information from K8S conditions such as scheduler logs etc
    """

    def __init__(self):

        self._config = kubernetes.client.Configuration()

        kubernetes.config.load_kube_config(client_configuration=self._config)
        self._kube_client = kubernetes.client.ApiClient(self._config)

    @property
    def kube_client(self):
        return self._kube_client

    @print_report
    def report_training_pods(self, name: str) -> str:
        core_api = kubernetes.client.CoreV1Api(self.kube_client)
        crds = kubernetes.client.CustomObjectsApi(self.kube_client)

        report = ""

        try:
            model_training = crds.get_namespaced_custom_object(
                group, version, training_namespace, training_plural, name.lower()
            )
            report += f"Model training: {model_training}\n"
        except kubernetes.client.rest.ApiException as e:
            report += f"Model training: unable to retrieve. {_api_exc_print(e)}\n"
            return report

        mt_pod = model_training.get("status", {}).get("podName", "")
        if not mt_pod:
            report += (
                "Training pod: unable to fetch pod name from model training status\n"
            )
            return report

        try:
            pod = core_api.read_namespaced_pod(mt_pod, training_namespace)
            report += f"Pod: {pod}\n"
        except kubernetes.client.rest.ApiException as e:
            report += f"Pod: unable to retrieve. {_api_exc_print(e)}\n"
            return report

        return report

    @print_report
    def report_packaging_pods(self, name: str) -> str:

        core_api = kubernetes.client.CoreV1Api(self.kube_client)
        crds = kubernetes.client.CustomObjectsApi(self.kube_client)

        report = ""

        try:
            model_packaging = crds.get_namespaced_custom_object(
                group, version, packaging_namespace, packaging_plural, name.lower()
            )
            report += f"Model packaging: {model_packaging}\n"
        except kubernetes.client.rest.ApiException as e:
            report += f"Model packaging: unable to retrieve. {_api_exc_print(e)}\n"
            return report

        mt_pod = model_packaging.get("status", {}).get("podName", "")
        if not mt_pod:
            report += (
                "Packaging pod: unable to fetch pod name from model packaging status\n"
            )
            return report

        try:
            pod = core_api.read_namespaced_pod(mt_pod, packaging_namespace)
            report += f"Pod: {pod}\n"
        except kubernetes.client.rest.ApiException as e:
            report += f"Pod: unable to retrieve. {_api_exc_print(e)}\n"
            return report

        return report

    @print_report
    def report_model_deployment_pods(self, name: str) -> str:
        core_api = kubernetes.client.CoreV1Api(self.kube_client)
        crds = kubernetes.client.CustomObjectsApi(self.kube_client)

        report = ""

        try:
            model_deployment = crds.get_namespaced_custom_object(
                group, version, deployment_namespace, deployment_plural, name.lower()
            )
            report += f"Model deployment: {model_deployment}\n"
        except kubernetes.client.rest.ApiException as e:
            report += f"Model deployment: unable to retrieve. {_api_exc_print(e)}\n"
            return report

        pods = core_api.list_namespaced_pod(
            deployment_namespace, label_selector=f"{ODAHU_DEPL_LABEL}={name}"
        )

        report += f"Model deployment pods: {pods}"

        return report
