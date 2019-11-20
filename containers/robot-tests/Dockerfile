#
#    Copyright 2019 EPAM Systems
#
#    Licensed under the Apache License, Version 2.0 (the "License");
#    you may not use this file except in compliance with the License.
#    You may obtain a copy of the License at
#
#       http://www.apache.org/licenses/LICENSE-2.0
#
#    Unless required by applicable law or agreed to in writing, software
#    distributed under the License is distributed on an "AS IS" BASIS,
#    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
#    See the License for the specific language governing permissions and
#    limitations under the License.
#

FROM python:3.6.8

ENV KUBECTL_VERSION=v1.13.8 \
    HELM_VERSION=v2.14.3

# Install python package dependencies and cloud CLI tools
RUN apt-get update && apt-get install -y --no-install-recommends apt-transport-https software-properties-common \
    make build-essential zip libssl-dev libffi-dev zlib1g-dev libjpeg-dev git jq=1.5+dfsg-1.3 && \
    DISTR="stretch" && \
    curl -s https://packages.cloud.google.com/apt/doc/apt-key.gpg | apt-key add - && \
    echo "deb http://packages.cloud.google.com/apt cloud-sdk-$DISTR main" | tee /etc/apt/sources.list.d/google-cloud-sdk.list && \
    curl -s https://packages.microsoft.com/keys/microsoft.asc | apt-key add - && \
    echo "deb [arch=amd64] https://packages.microsoft.com/repos/azure-cli/ $DISTR main" | tee /etc/apt/sources.list.d/azure-cli.list && \
    apt-get -qqy update && \
    apt-get install -y --no-install-recommends azure-cli google-cloud-sdk &&\
    apt-get clean all && rm -rf /var/lib/apt/lists/*

RUN pip install --disable-pip-version-check --upgrade pip==18.1 pipenv==2018.10.13 awscli

# Install kubectl
RUN curl -fsSLO https://storage.googleapis.com/kubernetes-release/release/${KUBECTL_VERSION}/bin/linux/amd64/kubectl && \
    mv kubectl /usr/local/bin/kubectl && \
    chmod a+x /usr/local/bin/kubectl

# Install Helm
RUN curl -fsSLO https://kubernetes-helm.storage.googleapis.com/helm-${HELM_VERSION}-linux-amd64.tar.gz && \
    mkdir -p /tmp/helm && mv helm-${HELM_VERSION}-linux-amd64.tar.gz /tmp/helm && \
    tar xzf /tmp/helm/helm-${HELM_VERSION}-linux-amd64.tar.gz -C /tmp/helm && \
    mv /tmp/helm/linux-amd64/helm /usr/local/bin/helm && rm -rf /tmp/helm

# Install python dependecies
COPY packages/sdk/Pipfile packages/sdk/Pipfile.lock /opt/odahu-flow/packages/sdk/
RUN  cd /opt/odahu-flow/packages/sdk && pipenv install --system --three --dev
COPY packages/cli/Pipfile packages/cli/Pipfile.lock /opt/odahu-flow/packages/cli/
RUN  cd /opt/odahu-flow/packages/cli && pipenv install --system --three --dev
COPY packages/robot/Pipfile packages/robot/Pipfile.lock /opt/odahu-flow/packages/robot/
RUN  cd /opt/odahu-flow/packages/robot && pipenv install --system --three --dev

COPY scripts /opt/odahu-flow/scripts
RUN chmod -R a+x /opt/odahu-flow/scripts/*
COPY Makefile /opt/odahu-flow/Makefile
COPY packages /opt/odahu-flow/packages

RUN cd /opt/odahu-flow/ && make BUILD_PARAMS="--no-deps" install-all
