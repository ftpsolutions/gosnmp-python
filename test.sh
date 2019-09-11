#!/bin/bash

set -e -o xtrace

# Allows us to run from a relative path
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
pushd ${DIR}

# permit overriding the container name (for avoiding name clashes)
if [[ -z "$CONTAINER_NAME" ]]; then
    CONTAINER_NAME=`basename $PWD-test`
fi

# Make sure we clean up
function finish {
    # Remove the container when we're done
    echo "Cleaning up...."
    # Go back to original dir
    popd
}
trap finish EXIT

IMAGE_TAG=gosnmp_python_test_build-${CONTAINER_NAME}

docker build --tag ${IMAGE_TAG} -f Dockerfile_test_build .

DOCKER_CMD=py.test
if [ "$#" -gt 0 ]; then
    echo "Using command from args"
    DOCKER_CMD=$@
fi

docker run --name ${CONTAINER_NAME} --rm -it ${IMAGE_TAG} ${DOCKER_CMD}
