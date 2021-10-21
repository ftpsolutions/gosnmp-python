#!/usr/bin/env bash

set -e -x

DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
pushd "${DIR}"

if [[ -f "/tmp/native_build.sh" ]]; then
  /tmp/native_build.sh gosnmp "${1}"
elif which native_build.sh; then
  native_build.sh gosnmp "${1}"
else
  .common/go-binding-utils/native_build.sh gosnmp "${1}"
fi
