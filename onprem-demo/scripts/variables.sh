#!/bin/bash
set -uo pipefail
IFS=$'\n\t'
export BACALHAU_API_HOST=127.0.0.1
export BACALHAU_API_PORT=20000
export DATAFOLDER=/tmp/bacalhau-onprem-demo/data
export LOGFOLDER="$DATAFOLDER/logs"
export LOGFILE="$LOGFOLDER/accesspoint.log"
export IMAGEFOLDER="$DATAFOLDER/images"
export HTTP_ENDPOINT=http://172.17.0.1:9600/publish
export SOURCE_LOGS_DOCKER_IMAGE=quay.io/bacalhau/onprem-demo-source-logs:latest
export SOURCE_IMAGES_DOCKER_IMAGE=quay.io/bacalhau/onprem-demo-source-images:latest
export SINK_INFERENCE_SERVER_DOCKER_IMAGE=quay.io/bacalhau/onprem-demo-sink-inference-server:latest
export PREDICTABLE_API_PORT=1
export BACALHAU_LOCAL_DIRECTORY_ALLOW_LIST=$DATAFOLDER
export SKIP_IMAGE_PULL=1
export BACALHAU_STREAMING_MODE=1
