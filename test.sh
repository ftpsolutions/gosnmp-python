#!/bin/bash

set -e -o xtrace

# Allows us to run from a relative path
DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
pushd ${DIR}

# permit overriding the container name (for avoiding name clashes)
if [[ -z "$CONTAINER_NAME" ]]; then
  CONTAINER_NAME=$(basename $PWD-test)
fi

# Make sure we clean up
function finish() {
  echo "Cleaning up...."

  # Remove the container when we're done
  docker rm -f ${CONTAINER_NAME} || true

  # Go back to original dir
  popd
}
trap finish EXIT

IMAGE_TAG=gosnmp_python_test_build-${CONTAINER_NAME}

if [ -z "${SKIP_BUILD}" ]; then
  docker build --tag ${IMAGE_TAG} -f Dockerfile_test_build .
fi

DOCKER_CMD=py.test
if [ "$#" -gt 0 ]; then
  echo "Using command from args"
  DOCKER_CMD=$@
fi

# Define MOUNT_WORKSPACE to mount this workspace inside the docker container
WORKSPACE_VOLUME=""
if [ ! -z "${MOUNT_WORKSPACE}" ]; then
  WORKSPACE_VOLUME="-v $(pwd):/workspace"
fi

# run the tests / shell
docker run --name ${CONTAINER_NAME} --rm -d ${WORKSPACE_VOLUME} ${IMAGE_TAG} tail -F /dev/null
docker exec -it ${CONTAINER_NAME} ${DOCKER_CMD}

# extract the package for deployment
NAME=$(cat setup.py | grep 'name=' | xargs | cut -d '=' -f 2 | tr -d ',')
VERSION=$(cat setup.py | grep 'version=' | xargs | cut -d '=' -f 2 | tr -d ',')
PACKAGE="${NAME}-${VERSION}"
docker cp ${CONTAINER_NAME}:/workspace/dist/${PACKAGE}.tar.gz .
echo "Deployable package follows:"
ls -al ${PACKAGE}.tar.gz
echo ""
