#!/bin/bash

OUTPUT_DIR=$PWD/dist
mkdir -p ${OUTPUT_DIR}

COMMIT_HASH=$(git rev-parse --short=8 HEAD 2>/dev/null)
BUILD_TIME=$(date +%FT%T%z)
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ${OUTPUT_DIR}/credhub-service-broker -ldflags "-X github.com/rabobank/credhub-service-broker/conf/CommitHash=${COMMIT_HASH} -X github.com/rabobank/credhub-service-broker/conf.BuildTime=${BUILD_TIME}" .
