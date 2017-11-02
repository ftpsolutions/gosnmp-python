## gosnmp-python

The purpose of this module is to provide a Python interface to the Golang
[gosnmp](https://github.com/soniah/gosnmp) module.

It was made very easy with the help of the Golang
[gosnmp](https://github.com/soniah/gosnmp) module.

#### Limitations

* I haven't written any tests for the Go side because (excuses)
* It doesn't work with Python 2.7.10 delivered with Mac OS (use brew to install Python 2.7.13 or something)
* Python command needs to be prefixd with GODEBUG=cgocheck=0 (or have that in the environment)  

#### How do I make use of this?

Right now I'm still working on how to put it all together as a Python module, so here are the raw steps.

#### Prerequisites

* Go 1.9
* Python 2.7.13 or newer (Python 2.7.10 delivered with Mac OS doesn't seem to work)
* pip
* virtualenvwrapper

#### Setup (for dev)

* ```mkvirtualenvwrapper -p (/path/to/python) gosnmp-python``` 
* ```pip install -r requirements-dev.txt```
* ```./build.sh```
* ```py.test -v```

#### What's worth knowing if I want to further the development?

* gopy doesn't like Go interfaces; so make sure you don't have any public (exported) interfaces
    * this includes a struct with a public property that may eventually lead to an interface
    * e.g. Session.snmp is private (because that object leads to gosnmp which has interfaces)

#### Example Python usage

To create an SNMPv2 session in Python do the following:

```
from gosnmp import create_snmpv2c_session

session = create_snmpv2c_session(
    hostname='1.2.3.4',
    community='public',
)

print session.get('.1.3.6.1.2.1.1.5.0')
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
