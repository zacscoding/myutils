#!/usr/bin/env bash

SCRIPT_PATH=$(cd "$(dirname $0)" && pwd)

rm -rf ${SCRIPT_PATH}/bin/myutils
go build -o ${SCRIPT_PATH}/bin/myutils ${SCRIPT_PATH}/cmd/myutils/
