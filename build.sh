#!/bin/bash

set -e -o xtrace

which go 2>/dev/null 1>/dev/null
if [[ $? -ne 0 ]]; then
    echo "error: failed to find go binary- do you have Go 1.13 installed?"
    exit 1
fi

GOVERSION=`go version`
if [[ $GOVERSION != *"go1.13"* ]]; then
    echo "error: Go version is not 1.13 (was $GOVERSION)"
    exit 1
fi

export PYTHONPATH=`pwd`/src/github.com/go-python/gopy/
# Use the default go binary path - the way to do it with newer versions of golang!
PATH=${PATH}:~/go/bin

echo "cleaning up output folder"
rm -frv gosnmp_python/*.pyc
rm -frv gosnmp_python/py2/*
echo ""

if [[ "$1" == "clean" ]]; then
    exit 0
fi

if [[ "$1" != "fast" ]]; then
    echo "building gosnmp"
    go build -x -a -mod readonly github.com/ftpsolutions/gosnmp
    echo ""

    echo "building gopy"
    go build -x -a github.com/go-python/gopy
    echo ""

    echo "installing gopy"
    go install -i -mod readonly github.com/go-python/gopy
    echo ""

    # Use a specific version!
    echo "getting goimports"
    go get golang.org/x/tools/cmd/goimports@v0.0.0-20190910044552-dd2b5c81c578
fi

# Using a special version of pybindgen to fix some memory leaks specific to our use case
# https://github.com/ftpsolutions/pybindgen
echo "installing pybindgen - required for gopy"
pip install --trusted-host imdemo.ftpsolutions.com.au \
    --extra-index-url http://imdemo.ftpsolutions.com.au:9011/ \
    pybindgen==0.20.0.post2+gcab0b4a

echo "build gosnmp_python bindings for py2"
./gopy build -output="gosnmp_python/py2" -symbols=true -vm=$(which python) gosnmp_python/gosnmp_python_go
echo ""

# Yep - this is highly questionable
# This requires an entry in LD_LIBRARY_PATH to work
SHARED_OBJ_DIR=/usr/local/lib/gopy/
echo "copying shared objects to ${SHARED_OBJ_DIR}"
mkdir -p ${SHARED_OBJ_DIR}
cp gosnmp_python/py2/gosnmp_python_go_go.so ${SHARED_OBJ_DIR}

# gopy doesn't seem to support Python3 as yet
# echo "build gosnmp_python bindings for py3"
# ./gopy bind -lang="py3" -output="gosnmp_python/py3" -symbols=true -work=false gosnmp_python
# echo ""

# No more support for cffi
# echo "build gosnmp_python bindings for cffi"
# ./gopy bind -api="cffi" -output="gosnmp_python/cffi" -symbols=true -work=false gosnmp_python
# echo ""

echo "cleaning up"
find . | grep -E "(__pycache__|\.pyc|\.pyo$)" | xargs rm -rf
echo ""