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

ENV CLI_RCLONE_VERSION=v1.51.0 \
    KUBECTL_VERSION=v1.14.10

RUN wget -qO /usr/local/bin/kubectl "https://storage.googleapis.com/kubernetes-release/release/${KUBECTL_VERSION}/bin/linux/amd64/kubectl" && \
    wget -qO /tmp/rclone.zip "https://downloads.rclone.org/${CLI_RCLONE_VERSION}/rclone-${CLI_RCLONE_VERSION}-linux-amd64.zip" && \
    mkdir /tmp/rclone && unzip /tmp/rclone.zip -d /tmp/rclone && \
    cp "/tmp/rclone/rclone-${CLI_RCLONE_VERSION}-linux-amd64/rclone" /usr/local/bin/rclone && \
    chmod a+x /usr/local/bin/rclone /usr/local/bin/kubectl

RUN pip install --disable-pip-version-check --upgrade pip==18.1 pipenv==2018.10.13 awscli

ARG MINICONDA_URL=https://repo.anaconda.com/miniconda/Miniconda3-4.7.12-Linux-x86_64.sh
ENV LANG=C.UTF-8 LC_ALL=C.UTF-8 \
    PATH="${PATH}:/opt/conda/bin" \
    ODAHU_CONDA_ENV_NAME=odahu_model

# Install conda
RUN wget --quiet ${MINICONDA_URL} -O ~/miniconda.sh && \
    /bin/bash ~/miniconda.sh -b -p /opt/conda && \
    rm ~/miniconda.sh && \
    conda clean -tipsy && \
    ln -s /opt/conda/etc/profile.d/conda.sh /etc/profile.d/conda.sh && \
    conda create --name ${ODAHU_CONDA_ENV_NAME} python --no-default-packages

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
