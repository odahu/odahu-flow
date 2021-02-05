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

FROM golang:1.14-alpine as builder

ENV FEEDBACK_DIR="/go/src/github.com/odahu/odahu-flow/packages/feedback"
WORKDIR "${FEEDBACK_DIR}"

RUN apk add -u ca-certificates git gcc musl-dev make

COPY packages/commons/ ../commons
COPY packages/feedback/ ./

RUN GOOS=linux GOARCH=amd64 make build-all

#########################################################
#########################################################
#########################################################

FROM alpine:3.12 as collector

ENV ODAHUFLOW_DIR="/opt/odahu-flow"

WORKDIR "${ODAHUFLOW_DIR}"
COPY --from=builder /go/src/github.com/odahu/odahu-flow/packages/feedback/collector ./
CMD ["./collector"]

#########################################################
#########################################################
#########################################################

FROM alpine:3.12 as rq-catcher

ENV ODAHUFLOW_DIR="/opt/odahu-flow"

WORKDIR "${ODAHUFLOW_DIR}"
COPY --from=builder /go/src/github.com/odahu/odahu-flow/packages/feedback/rq-catcher ./
CMD ["./rq-catcher"]