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

ENV OPERATOR_DIR="/app" \
    PATH="$PATH:/go/bin:/usr/lib/go-1.14/bin" \
    GOPATH="/go"

WORKDIR "${OPERATOR_DIR}"

RUN apt-get update -qq && \
    apt-get install -y software-properties-common && \
    add-apt-repository -y ppa:longsleep/golang-backports && \
    apt-get update -qq && \
    apt-get install -y git gcc make golang-1.14-go

COPY packages/commons ../commons
COPY packages/operator ./

RUN GOOS=linux GOARCH=amd64 make build-all

#########################################################
#########################################################
#########################################################

FROM ubuntu:18.04 as operator

ENV ODAHUFLOW_DIR="/opt/odahu-flow"
RUN apt-get -yq update && \
    apt-get -yqq install ca-certificates && \
    apt-get clean && rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/operator "${ODAHUFLOW_DIR}/"
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

COPY --from=builder /app/webserver "${ODAHUFLOW_DIR}/"

WORKDIR "${ODAHUFLOW_DIR}"
CMD ["./webserver"]


#########################################################
#########################################################
#########################################################

FROM ubuntu:18.04 as tools

ENV ODAHUFLOW_DIR="/opt/odahu-flow"

RUN apt-get -yq update && \
    apt-get -yqq install openssh-client ca-certificates && \
    apt-get clean && rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/odahu-tools "${ODAHUFLOW_DIR}/"

WORKDIR "${ODAHUFLOW_DIR}"
CMD ["./odahu-tools"]


#########################################################
#########################################################
#########################################################

FROM ubuntu:18.04 as controller

ENV ODAHUFLOW_DIR="/opt/odahu-flow" \
    GIN_MODE="release"
RUN apt-get -yq update && \
    apt-get -yqq install openssh-client ca-certificates && \
    apt-get clean && rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/controller "${ODAHUFLOW_DIR}/"

WORKDIR "${ODAHUFLOW_DIR}"
CMD ["./controller"]

#########################################################
#########################################################
#########################################################


FROM ubuntu:18.04 as service-catalog

ENV ODAHUFLOW_DIR="/opt/odahu-flow" \
    GIN_MODE="release"

RUN apt-get -yq update && \
    apt-get -yqq install ca-certificates && \
    apt-get clean && rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/service-catalog "${ODAHUFLOW_DIR}/"

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

COPY --from=builder /app/trainer "${WORK_DIR}/"

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

COPY --from=builder /app/packager "${WORK_DIR}/"

CMD ["./packager"]
