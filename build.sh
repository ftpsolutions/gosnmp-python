#!/bin/bash

which go 2>/dev/null 1>/dev/null
if [[ $? -ne 0 ]]; then
    echo "error: failed to find go binary- do you have Go 1.9 or Go 1.10 installed?"
    exit 1
fi

GOVERSION=`go version`
if [[ $GOVERSION != *"go1.9"* ]] && [[ $GOVERSION != *"go1.10"* ]]; then
    echo "error: Go version is not 1.9 or 1.10 (was $GOVERSION)"
    exit 1
fi

export GOPATH=`pwd`

export PYTHONPATH=`pwd`/src/github.com/go-python/gopy/

echo "cleaning up output folder"
rm -frv gosnmp_python/*.pyc
rm -frv gosnmp_python/py2/*.pyc
rm -frv gosnmp_python/py2/gosnmp_python.so
rm -frv gosnmp_python/cffi/*.pyc
rm -frv gosnmp_python/cffi/_gosnmp_python.so
rm -frv gosnmp_python/cffi/gosnmp_python.py
echo ""

if [[ "$1" != "fast" ]]; then
    echo "getting assert"
    go get -v -a github.com/stretchr/testify/assert
    echo ""

    echo "getting gosnmp"
    go get -v -a github.com/soniah/gosnmp
    echo ""

    echo "building gosnmp"
    go build -a -x github.com/soniah/gosnmp
    echo ""

    echo "getting gopy"
    go get -v -a github.com/go-python/gopy
    echo ""

    if [[ $GOVERSION == *"go1.10"* ]]; then
        echo "fix errant pkg-config call in gopy (because we're running Go1.10)"
        sed 's^//#cgo pkg-config: %\[2\]s --cflags --libs^//#cgo pkg-config: %\[2\]s^g' src/github.com/go-python/gopy/bind/gengo.go > temp.go
        mv temp.go src/github.com/go-python/gopy/bind/gengo.go
    fi

    echo "building gopy"
    go build -x -a github.com/go-python/gopy
    echo ""

    echo "building gosnmp_python"
    go build -x -a gosnmp_python
    echo ""
fi

echo "build gosnmp_python bindings for py2"
./gopy bind -lang="py2" -output="gosnmp_python/py2" -symbols=true -work=false gosnmp_python
echo ""

echo "build gosnmp_python bindings for cffi"
./gopy bind -lang="cffi" -output="gosnmp_python/cffi" -symbols=true -work=false gosnmp_python
echo ""

echo "fix broken cffi output (this is yuck)"
sed 's/py_kwd_011, \[\]int/py_kwd_011, list/g' gosnmp_python/cffi/gosnmp_python.py > temp.py
mv temp.py gosnmp_python/cffi/gosnmp_python.py
