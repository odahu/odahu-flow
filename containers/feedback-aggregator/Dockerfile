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

FROM golang:1.11.6-alpine3.9 as builder

ENV FEEDBACK_DIR="/go/src/github.com/odahu/odahu-flow/odahuFlow/feedback-aggregator"
WORKDIR "${FEEDBACK_DIR}"

RUN apk add -u ca-certificates git gcc musl-dev make && \
    wget https://github.com/golang/dep/releases/download/v0.5.1/dep-linux-amd64 -O /usr/local/bin/dep && \
    chmod +x /usr/local/bin/dep

COPY odahuFlow/feedback-aggregator/Gopkg.toml odahuFlow/feedback-aggregator/Gopkg.lock ./
RUN dep ensure -vendor-only

COPY odahuFlow/feedback-aggregator/ ./

RUN GOOS=linux GOARCH=amd64 make build-all

#########################################################
#########################################################
#########################################################

FROM alpine:3.9 as aggregator

ENV ODAHUFLOW_DIR="/opt/odahuFlow"

WORKDIR "${ODAHUFLOW_DIR}"
COPY --from=builder /go/src/github.com/odahu/odahu-flow/odahuFlow/feedback-aggregator/aggregator ./
CMD ["./aggregator"]

#########################################################
#########################################################
#########################################################

FROM alpine:3.9 as collector

ENV ODAHUFLOW_DIR="/opt/odahuFlow"

WORKDIR "${ODAHUFLOW_DIR}"
COPY --from=builder /go/src/github.com/odahu/odahu-flow/odahuFlow/feedback-aggregator/collector ./
CMD ["./collector"]