#!/bin/bash

which go 2>/dev/null 1>/dev/null
if [[ $? -ne 0 ]]; then
    echo "error: failed to find go binary- do you have Go 1.9 installed?"
    exit 1
fi

GOVERSION=`go version`
if [[ $GOVERSION != *"go1.9"* ]]; then
    echo "error: Go version is not 1.9 (was $GOVERSION)"
    exit 1
fi

export GOPATH=`pwd`

echo "getting gopy"
go get -v github.com/go-python/gopy

export PYTHONPATH=`pwd`/src/github.com/go-python/gopy/

echo "building gopy"
go build -a -x github.com/go-python/gopy

echo "building gosnmp_python"
go build -a -x gosnmp_python

echo "build gosnmp_python bindings"
gopy bind -lang="py2" -output="gosnmp" gosnmp_python
