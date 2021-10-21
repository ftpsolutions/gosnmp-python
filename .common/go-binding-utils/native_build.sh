#!/usr/bin/env bash

set -e -x

if [[ "${1}" == "" ]]; then
  echo "error: first argument must be name of lib (e.g. gomssql)"
  exit 1
fi

NAME_OF_LIB="${1}"

if [[ "${VIRTUAL_ENV}" == "" ]]; then
  echo "error: a virtualenv must be activated"
  exit 1
fi

# TODO
# GO_VERSION="$(go version)"
# if [[ ${GO_VERSION} != *"go1.13"* ]]; then
#   echo "error: Go version is not 1.13 (was ${GO_VERSION})"
#   exit 1
# fi

PKG_CONFIG_PATH="$(pwd)"
export PKG_CONFIG_PATH

OUTPUT_PATH="${NAME_OF_LIB}_python/built"

if [[ "${2}" != "fast" ]]; then
  echo "installing python deps..."
  pip install --upgrade -r requirements-dev.txt
  echo ""

  echo "installing goimports..."
  go get golang.org/x/tools/cmd/goimports
  echo ""

  echo "installing gopy..."
  go get github.com/go-python/gopy@v0.3.4
  echo ""
fi

echo "cleaning up output folder..."
rm -frv "${OUTPUT_PATH:?}/*" || true
mkdir -p "${OUTPUT_PATH}" || true
touch "${OUTPUT_PATH}/__init__.py" || true
echo ""

echo "building ${NAME_OF_LIB}-python..."
export PATH=${PATH}:~/go/bin/
export CFLAGS
export C_INCLUDE_PATH
gopy build -output="${OUTPUT_PATH}" -symbols=true -vm="$(command -v python)" "./${NAME_OF_LIB}_python_go"
echo ""

echo "hacking in some sed fixes..."
sed -i'.bak' "s/import _${NAME_OF_LIB}_python_go/from ${NAME_OF_LIB}_python.built import _${NAME_OF_LIB}_python_go/g" "${NAME_OF_LIB}_python/built/${NAME_OF_LIB}_python_go.py"
sed -i'.bak' "s/import go/from ${NAME_OF_LIB}_python.built import go/g" "${NAME_OF_LIB}_python/built/${NAME_OF_LIB}_python_go.py"
sed -i'.bak' "s/import _${NAME_OF_LIB}_python_go/from ${NAME_OF_LIB}_python.built import _${NAME_OF_LIB}_python_go/g" "${NAME_OF_LIB}_python/built/go.py"
echo ""
