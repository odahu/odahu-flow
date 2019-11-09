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

ENV OPERATOR_DIR="/go/src/github.com/odahu/odahu-flow/odahuFlow/operator" \
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

COPY odahuFlow/operator/Gopkg.toml odahuFlow/operator/Gopkg.lock ./
RUN dep ensure -v -vendor-only

COPY odahuFlow/operator/ ./

RUN GOOS=linux GOARCH=amd64 make build-all

#########################################################
#########################################################
#########################################################

FROM ubuntu:18.04 as operator

ENV ODAHUFLOW_DIR="/opt/odahuFlow"
RUN  apt-get -yq update && \
     apt-get -yqq install ca-certificates

COPY --from=builder /go/src/github.com/odahu/odahu-flow/odahuFlow/operator/operator "${ODAHUFLOW_DIR}/"
WORKDIR "${ODAHUFLOW_DIR}"
CMD ["./operator"]

#########################################################
#########################################################
#########################################################

FROM ubuntu:18.04 as edi

ENV ODAHUFLOW_DIR="/opt/odahuFlow" \
    TEMPLATE_FOLDER="/opt/odahuFlow/templates" \
    GIN_MODE="release"
RUN  apt-get -yq update && \
     apt-get -yqq install openssh-client ca-certificates

COPY --from=builder /go/src/github.com/odahu/odahu-flow/odahuFlow/operator/edi "${ODAHUFLOW_DIR}/"
COPY odahuFlow/operator/templates "${TEMPLATE_FOLDER}/"

WORKDIR "${ODAHUFLOW_DIR}"
CMD ["./edi"]

#########################################################
#########################################################
#########################################################

FROM ubuntu:18.04 as service-catalog

ENV ODAHUFLOW_DIR="/opt/odahuFlow" \
    TEMPLATE_FOLDER="/opt/odahuFlow/templates" \
    GIN_MODE="release"

RUN  apt-get -yq update && \
     apt-get -yqq install ca-certificates

COPY --from=builder /go/src/github.com/odahu/odahu-flow/odahuFlow/operator/service-catalog "${ODAHUFLOW_DIR}/"
COPY odahuFlow/operator/templates "${TEMPLATE_FOLDER}/"

WORKDIR "${ODAHUFLOW_DIR}"
CMD ["./edi"]

#########################################################
#########################################################
#########################################################

FROM ubuntu:18.04 as model-trainer

ENV DEBIAN_FRONTEND=noninteractive \
    LC_ALL=en_US.UTF-8 LANG=en_US.UTF-8 LANGUAGE=en_US.UTF-8 \
    WORK_DIR="/opt/odahuFlow"

WORKDIR "${WORK_DIR}/"

RUN  apt-get -yq update && \
     apt-get -yqq install ca-certificates

COPY --from=builder /go/src/github.com/odahu/odahu-flow/odahuFlow/operator/trainer "${WORK_DIR}/"

CMD ["./trainer"]

#########################################################
#########################################################
#########################################################

FROM ubuntu:18.04 as model-packager

ENV DEBIAN_FRONTEND=noninteractive \
    LC_ALL=en_US.UTF-8 LANG=en_US.UTF-8 LANGUAGE=en_US.UTF-8 \
    WORK_DIR="/opt/odahuFlow"

RUN  apt-get -yq update && \
     apt-get -yqq install ca-certificates

WORKDIR "${WORK_DIR}/"

COPY --from=builder /go/src/github.com/odahu/odahu-flow/odahuFlow/operator/packager "${WORK_DIR}/"

CMD ["./packager"]
