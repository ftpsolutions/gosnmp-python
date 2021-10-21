#!/usr/bin/env bash

set -e

if [[ "${VIRTUAL_ENV}" == "" ]]; then
  echo "error: a virtualenv must be activated"
  exit 1
fi

export TEST_COMMUNITY
export TEST_DATABASE
export TEST_HOSTNAME
export TEST_PASSWORD
export TEST_PORT
export TEST_RETRIES
export TEST_TIMEOUT
export TEST_USERNAME

FOCUS="test/smoke_test.py"
if [[ "${1}" != "" ]]; then
  FOCUS="${*}"
fi

GODEBUG=cgocheck=0 python -m pytest --junit-xml "/srv/test_results/junit_results.xml" -vv -s "${FOCUS}"
