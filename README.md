# gosnmp-python

This library binds [our fork](https://github.com/ftpsolutions/gosnmp) of [gosnmp](https://github.com/gosnmp/gosnmp) for use in Python3
using [gopy](https://github.com/go-python/gopy).

## Versions

All versions 1.0.0 and up support Python 3 only! If you need Python 2 support, check out the following:

- [https://github.com/ftpsolutions/gosnmp-python/tree/v0.91](https://github.com/ftpsolutions/gosnmp-python/tree/v0.91)
- https://pypi.org/project/gosnmp-python/0.91/

## Concept

In the early days `gopy` was fairly limited in it's ability to track object allocation across the border of Go and Python.

As a result, our implementation is fairly naive- we only pass primitive types from Go to Python (nothing that comes by reference).

Session are managed entirely on the Go side and identified with an integer- here are a few function signatures to demonstrate:

- `NewRPCSessionV2c(hostname string, port int, community string, timeout, retries int) uint64`
- `RPCConnect(sessionID uint64) error`
- `RPCGet(sessionID uint64, oid string) (string, error)`
- `RPCClose(sessionID uint64) error`

The functions that return complex data do so in a special JSON-based format- at this point `gopy` does it's magic and those functions are
made available to Python.

We then have `RPCSession` abstraction on the Python side that pulls things together in a class for convenience (saving you need the to keep
track of the identifiers and handling deserialisation).

## Weird gotchas

We're building for Python3 and we use a `python-config` script for Python3 however we're using a `python.pc` file from Python2.

I dunno why this has to be this way, but it's the only way I can get the C Python API GIL lock/unlock calls to be available to
Go (`C.PyEval_SaveThread()` and `C.PyEval_RestoreThread(tState)`).

So if you're wondering why Python2 still comes into it here and there- that's why. Doesn't seem to cause any problems.

## Prerequisites

- MacOS
- CPython3.8+
- virtualenv
    - `pip install virtualenv`
- pkgconfig
    - `brew install pkg-config`
- Docker

## Install (from PyPI)

```
python -m pip install gosnmp-python
```

## Setup

```
virtualenv -p $(which python3) venv
source venv/bin/activate
./fix_venv.sh
pip install -r requirements-dev.txt
```

## Build

```
source venv/bin/activate
./native_build.sh
```

If you're deep in the grind and want to iterate faster, you can invoke:

```
./native_build.sh fast
```

This skips the setup (assuming you've already done that).
