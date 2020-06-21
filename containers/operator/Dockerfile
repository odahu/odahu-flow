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

FROM ubuntu:18.04 as builder

ENV OPERATOR_DIR="/go/src/github.com/odahu/odahu-flow/packages/operator" \
    PATH="$PATH:/go/bin:/usr/lib/go-1.12/bin" \
    GOPATH="/go"

WORKDIR "${OPERATOR_DIR}"

RUN apt-get update -qq && \
    apt-get install -y software-properties-common && \
    add-apt-repository -y ppa:longsleep/golang-backports && \
    apt-get update -qq && \
    apt-get install -y git gcc make golang-1.12-go wget && \
    wget -q https://github.com/golang/dep/releases/download/v0.5.1/dep-linux-amd64 -O /usr/local/bin/dep && \
    chmod +x /usr/local/bin/dep

COPY packages/operator/Gopkg.toml packages/operator/Gopkg.lock ./
RUN dep ensure -v -vendor-only

COPY packages/operator/ ./

RUN GOOS=linux GOARCH=amd64 make build-all

#########################################################
#########################################################
#########################################################

FROM ubuntu:18.04 as operator

ENV ODAHUFLOW_DIR="/opt/odahu-flow"
RUN apt-get -yq update && \
    apt-get -yqq install ca-certificates && \
    apt-get clean && rm -rf /var/lib/apt/lists/*

COPY --from=builder /go/src/github.com/odahu/odahu-flow/packages/operator/operator "${ODAHUFLOW_DIR}/"
WORKDIR "${ODAHUFLOW_DIR}"
CMD ["./operator"]

#########################################################
#########################################################
#########################################################

FROM ubuntu:18.04 as api

ENV ODAHUFLOW_DIR="/opt/odahu-flow" \
    GIN_MODE="release"
RUN apt-get -yq update && \
    apt-get -yqq install openssh-client ca-certificates && \
    apt-get clean && rm -rf /var/lib/apt/lists/*

COPY --from=builder /go/src/github.com/odahu/odahu-flow/packages/operator/api "${ODAHUFLOW_DIR}/"

WORKDIR "${ODAHUFLOW_DIR}"
CMD ["./webserver"]

#########################################################
#########################################################
#########################################################

FROM ubuntu:18.04 as odahu-migrate

ENV ODAHUFLOW_DIR="/opt/odahu-flow" \
    GIN_MODE="release"
RUN apt-get -yq update && \
    apt-get -yqq install openssh-client ca-certificates && \
    apt-get clean && rm -rf /var/lib/apt/lists/*

COPY --from=builder /go/src/github.com/odahu/odahu-flow/packages/operator/odahu-migrate "${ODAHUFLOW_DIR}/"

WORKDIR "${ODAHUFLOW_DIR}"
CMD ["./odahu-migrate"]

#########################################################
#########################################################
#########################################################

FROM ubuntu:18.04 as service-catalog

ENV ODAHUFLOW_DIR="/opt/odahu-flow" \
    GIN_MODE="release"

RUN apt-get -yq update && \
    apt-get -yqq install ca-certificates && \
    apt-get clean && rm -rf /var/lib/apt/lists/*

COPY --from=builder /go/src/github.com/odahu/odahu-flow/packages/operator/service-catalog "${ODAHUFLOW_DIR}/"

WORKDIR "${ODAHUFLOW_DIR}"
CMD ["./service-catalog"]

#########################################################
#########################################################
#########################################################

FROM ubuntu:18.04 as model-trainer

ENV DEBIAN_FRONTEND=noninteractive \
    LC_ALL=en_US.UTF-8 LANG=en_US.UTF-8 LANGUAGE=en_US.UTF-8 \
    WORK_DIR="/opt/odahu-flow"

WORKDIR "${WORK_DIR}/"

RUN apt-get -yq update && \
    apt-get -yqq install ca-certificates pigz && \
    apt-get clean && rm -rf /var/lib/apt/lists/*

COPY --from=builder /go/src/github.com/odahu/odahu-flow/packages/operator/trainer "${WORK_DIR}/"

CMD ["./trainer"]

#########################################################
#########################################################
#########################################################

FROM ubuntu:18.04 as model-packager

ENV DEBIAN_FRONTEND=noninteractive \
    LC_ALL=en_US.UTF-8 LANG=en_US.UTF-8 LANGUAGE=en_US.UTF-8 \
    WORK_DIR="/opt/odahu-flow"

RUN apt-get -yq update && \
    apt-get -yqq install ca-certificates pigz && \
    apt-get clean && rm -rf /var/lib/apt/lists/*

WORKDIR "${WORK_DIR}/"

COPY --from=builder /go/src/github.com/odahu/odahu-flow/packages/operator/packager "${WORK_DIR}/"

CMD ["./packager"]
