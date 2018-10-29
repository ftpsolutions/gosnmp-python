## gosnmp-python

The purpose of this module is to provide a Python interface to the Golang [gosnmp](https://github.com/soniah/gosnmp) module.

It was made very easy with the help of the Golang [gopy](https://github.com/go-python/gopy) module.

#### Limitations

* Python command needs to be prefixed with GODEBUG=cgocheck=0 (or have that in the environment)
* I've not implemented walk (as I didn't need it for my use-case, I just use get_next with Python)
* Seems to have some odd memory problems with PyPy (via CFFI); lots of locks and stuff to try and alleviate that

#### How do I make use of this?

Right now I'm still working on how to put it all together as a Python module, so here are the raw steps.

#### Prerequisites

* Go 1.9 or Go 1.10
* Python 2.7+, Python 3.5+, PyPy 5.10+ or PyPy3 5.10+
* pip
* virtualenvwrapper
* pkgconfig/pkg-config

#### Installation (for prod)
* ```python setup.py install``` 

#### Making a python wheel install file (for distribution)
* ```python setup.py bdist_wheel``` 

#### Setup (for dev)

* ```mkvirtualenvwrapper -p (/path/to/pypy) gosnmp-python``` 
* ```pip install -r requirements-dev.txt```
* ```./build.sh```
* ```GODEBUG=cgocheck=0 py.test -v```

#### What's worth knowing if I want to further the development?

* gopy doesn't like Go interfaces; so make sure you don't have any public (exported) interfaces
    * this includes a struct with a public property that may eventually lead to an interface

#### Example Go RPCSession usage (simple session ID, calls return JSON)

There's no real reason why you'd want to do this (just use gosnmp on it's own)- it's more for reference:

```
package main

import (
	"gosnmp_python"
	"fmt"
)

func main() {

	sessionID := gosnmp_python.NewRPCSessionV2c(
		"1.2.3.4",
		161,
		"public",
		5,
		1,
	)

	err := gosnmp_python.RPCConnect(sessionID)
	if err != nil {
		panic(err)
	}

	jsonResult, err := gosnmp_python.RPCGet(sessionID, ".1.3.6.1.2.1.1.5.0")
	if err != nil {
		panic(err)
	}

	fmt.Println(jsonResult)

	err = gosnmp_python.RPCClose(sessionID)
	if err != nil {
		panic(err)
	}

}
```

#### Example Python usage (uses RPCSession underneath because of memory leaks between Go and Python with structs)

To create an SNMPv2 session in Python do the following:

```
from gosnmp_python import create_snmpv2c_session

session = create_snmpv2c_session(
    hostname='1.2.3.4',
    community='public',
)

session.connect()

print session.get('.1.3.6.1.2.1.1.5.0')

session.close()
```

Which returns an object like this:

```
SNMPVariable(
    oid='.1.3.6.1.2.1.1.5', 
    oid_index=0, 
    snmp_type=u'bytearray', 
    value='hostname.domain.com.au'
)
```
 
Some of this may feel a bit like [easysnmp](https://github.com/kamakazikamikaze/easysnmp); that's intentional, I was originally using that but I think its got some underlying thread-safety issues on the C side (particularly to do with SNMPv3).

No offence to the guys that contribute to that project- it's served me very well.
