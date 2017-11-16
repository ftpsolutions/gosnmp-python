## gosnmp-python

The purpose of this module is to provide a Python interface to the Golang
[gosnmp](https://github.com/soniah/gosnmp) module.

It was made very easy with the help of the Golang
[gopy](https://github.com/go-python/gopy) module.

#### Limitations

* I haven't written any tests for the Go side because (excuses)
* Doesn't seem to work with Python any more (after the CFFI change)
* Python command needs to be prefixed with GODEBUG=cgocheck=0 (or have that in the environment)
* I've not implemented walk (as I didn't need it for my use-case, I just use get_next with Python)
* PyPy seems to have memory problems (particularly when closing/opening lots of sessions in multiple threads)
    * I've put lots of mutexes in that apply only when PyPy is running- still seems flaky

#### How do I make use of this?

Right now I'm still working on how to put it all together as a Python module, so here are the raw steps.

#### Prerequisites

* Go 1.9
* Python 2.7.14+, or PyPy 5.9+ (PyPy is a bit unstable at the moment)
* pip
* virtualenvwrapper
* pkgconfig/pkg-config

#### Setup (for dev)

* ```mkvirtualenvwrapper -p (/path/to/pypy) gosnmp-python``` 
* ```pip install -r requirements-dev.txt```
* ```./build.sh```
* ```py.test -v```

#### What's worth knowing if I want to further the development?

* gopy doesn't like Go interfaces; so make sure you don't have any public (exported) interfaces
    * this includes a struct with a public property that may eventually lead to an interface
    * e.g. Session.snmp is private (because that object leads to gosnmp which has interfaces)
* I've left the GIL handling in potentially blocking call (for performance on the Python side);
  it doesn't seem to be used by PyPy and has odd behaviour with Python depending on the version,
  you may want to remove it altogether or try and make it work properly (if you use Python).

#### Example Go Session usage (complex Session struct, calls return MultiResult struct)

```
package main

import (
    "fmt"
    "gosnmp_python"
)

func main() {

    session := gosnmp_python.NewSessionv2c(
        "1.2.3.4"
        161
        "public"
        5
        1
    )
    
    err := session.Connect()
    if err != nil {
        panic(err)
    }
    
    multiResult, err := session.Get(".1.3.6.1.2.1.1.5.0")
    if err != nil {
        panic(err)
    }
    
    session.Close()
    if err != nil {
        panic(err)
    }

}
```

#### Example Go RPCSession usage (simple session ID, calls return JSON)

```
package main

import (
    "fmt"
    "gosnmp_python"
)

func main() {

    session := gosnmp_python.NewRPCSessionv2c(
        "1.2.3.4"
        161
        "public"
        5
        1
    )
    
    err := session.Connect()
    if err != nil {
        panic(err)
    }
    
    multiResult, err := session.Get(".1.3.6.1.2.1.1.5.0")
    if err != nil {
        panic(err)
    }
    
    session.Close()
    if err != nil {
        panic(err)
    }

}
```

#### Example Python usage (uses RPCSession underneath because of memory leaks between Go and Python with structs)

To create an SNMPv2 session in Python do the following:

```
from gosnmp import create_snmpv2c_session

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
    index=0, 
    snmp_type='string', 
    value='FTP_Switch1.ftpsolutions.com.au'
)
```
 
Some of this may feel a bit like [easysnmp](https://github.com/kamakazikamikaze/easysnmp); that's intentional,
I was originally using that but I think its got some underlying thread-safety issues on the C side (particularly
to do with SNMPv3).

No offence to the guys that contribute to that project- it's served me very well.
