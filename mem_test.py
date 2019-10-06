"""
Manual mem test.
1. Make a virtual env and install all deps
2. run `build.sh`
3. run `python mem_test.py`
4. Watch memory usage and make sure it doesn't get out of control
"""
from gosnmp_python import create_snmpv2c_session
from gosnmp_python.common import GoRuntimeError
import traceback
import time
from guppy import hpy

session = create_snmpv2c_session(hostname='127.0.0.1', community='dummy',
                                 timeout=1, port=1234)
session.connect()

print(hpy().heap())

err_count = 0
while True:
    try:
        session.get('.1.3.6.1.2.1.1.1.0')
        session.get(None)
    except KeyboardInterrupt:
        break
    except (RuntimeError, GoRuntimeError):
        err_count += 1

print(hpy().heap())

print('done')


