#!/usr/bin/env bash

set -e -x

if [[ "${VIRTUAL_ENV}" == "" ]]; then
  echo "error: VIRTUAL_ENV not set"
  exit 1
fi

DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# this is needed so that Go builds can find the Python headers etc
cp -frv "${DIR}/python-config" "${VIRTUAL_ENV}/bin/python-config"

# this is needed so that our Go builds can find the common tooling when invoked from within a Virtualenv
cp -frv "${DIR}/native_build.sh" "${VIRTUAL_ENV}/bin/native_build.sh"
