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

export PYTHONPATH=`pwd`/src/github.com/go-python/gopy/

echo "cleaning up output folder"
rm -frv gosnmp_python/*.pyc
rm -frv gosnmp_python/py2/*.pyc
rm -frv gosnmp_python/py2/gosnmp_python.so
rm -frv gosnmp_python/cffi/*.pyc
rm -frv gosnmp_python/cffi/_gosnmp_python.so
rm -frv gosnmp_python/cffi/gosnmp_python.py
echo ""

#echo "getting gosnmp"
#go get -v github.com/initialed85/gosnmp
#echo ""
#
#echo "building gosnmp"
#go build -a github.com/initialed85/gosnmp
#echo ""
#
#echo "getting gopy"
#go get -v github.com/go-python/gopy
#echo ""
#
#echo "building gopy"
#go build -a github.com/go-python/gopy
#echo ""
#
#echo "building gosnmp_python"
#go build -a gosnmp_python
#echo ""

echo "build gosnmp_python bindings for py2"
./gopy bind -lang="py2" -output="gosnmp_python/py2" gosnmp_python
echo ""

echo "build gosnmp_python bindings for cffi"
./gopy bind -lang="cffi" -output="gosnmp_python/cffi" gosnmp_python
echo ""

echo "fix broken cffi output (this is yuck)"
sed "s/py_kwd_011, \[\]int/py_kwd_011, list/g" gosnmp_python/cffi/gosnmp_python.py > temp.py
mv temp.py gosnmp_python/cffi/gosnmp_python.py
