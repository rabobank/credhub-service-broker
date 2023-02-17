#!/bin/bash

OUTPUT_DIR=$PWD/dist
mkdir -p ${OUTPUT_DIR}

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ${OUTPUT_DIR}/credhub-service-broker -ldflags "-X github.com/rabobank/credhub-service-broker/conf.VERSION=${VERSION} -X github.com/rabobank/credhub-service-broker/conf.COMMIT=${COMMIT}" .
