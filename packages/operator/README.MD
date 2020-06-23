## Setup services

### API service

Execute the following commands to build and launch the API service:

```bash
cd packages/operator

dep ensure

make generate-all
make build-all

# The API service will use the current kube config.
./webserver
# If you specify the local backend then the API service will startup etcd and kube API locally.
# Download the the kubebuilder asset by [link](https://github.com/kubernetes-sigs/kubebuilder/releases/tag/v2.3.1)
# Specify the KUBEBUILDER_ASSETS environment variable that must point to the unzip asset dir.
./webserver --backend-type local
```
