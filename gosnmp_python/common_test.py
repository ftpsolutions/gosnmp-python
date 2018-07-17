from __future__ import (absolute_import, division, print_function,
                        unicode_literals)

import unittest
from builtins import *
from collections import namedtuple

from future import standard_library
from hamcrest import assert_that, calling, equal_to, raises

from .common import (SNMPVariable, UnknownSNMPTypeError, handle_multi_result,
                     handle_multi_result_json)

standard_library.install_aliases()



# this is what comes across the border from Go for .Get and .GetNext
MultiResult = namedtuple('MultiResult', [
    'OID',
    'Type',
    'IsNull',
    'IsUnknown',
    'IsNoSuchInstance',
    'IsNoSuchObject',
    'IsEndOfMibView',
    'BoolValue',
    'IntValue',
    'FloatValue',
    'ByteArrayValue',
    'StringValue',
])

# this is what comes across the border from Go for .GetJSON and .GetNextJSON
_MULTI_RESULT_JSON_STRING = '{"OID":".1.2.3.4","Type":"string","IsNull":false,"IsUnknown":false,"IsNoSuchInstance":false,"IsNoSuchObject":false,"IsEndOfMibView":false,"BoolValue":false,"IntValue":0,"FloatValue":0,"ByteArrayValue":[],"StringValue":"some string"}'

_MULTI_RESULT_NOSUCHINSTANCE = MultiResult(
    OID='.1.2.3.4',
    Type='noSuchInstance',
    IsNull=False,
    IsUnknown=False,
    IsNoSuchInstance=True,
    IsNoSuchObject=False,
    IsEndOfMibView=False,
    BoolValue=False,
    IntValue=0,
    FloatValue=0.0,
    ByteArrayValue=[],
    StringValue='',
)

_MULTI_RESULT_NOSUCHOBJECT = MultiResult(
    OID='.1.2.3.4',
    Type='noSuchObject',
    IsNull=False,
    IsUnknown=False,
    IsNoSuchInstance=False,
    IsNoSuchObject=True,
    IsEndOfMibView=False,
    BoolValue=False,
    IntValue=0,
    FloatValue=0.0,
    ByteArrayValue=[],
    StringValue='',
)

_MULTI_RESULT_ENDOFMIBVIEW = MultiResult(
    OID='.1.2.3.4',
    Type='endOfMibView',
    IsNull=False,
    IsUnknown=False,
    IsNoSuchInstance=False,
    IsNoSuchObject=False,
    IsEndOfMibView=True,
    BoolValue=False,
    IntValue=0,
    FloatValue=0.0,
    ByteArrayValue=[],
    StringValue='',
)

_MULTI_RESULT_BOOL = MultiResult(
    OID='.1.2.3.4',
    Type='bool',
    IsNull=False,
    IsUnknown=False,
    IsNoSuchInstance=False,
    IsNoSuchObject=False,
    IsEndOfMibView=False,
    BoolValue=True,
    IntValue=0,
    FloatValue=0.0,
    ByteArrayValue=[],
    StringValue='',
)

_MULTI_RESULT_INT = MultiResult(
    OID='.1.2.3.4',
    Type='int',
    IsNull=False,
    IsUnknown=False,
    IsNoSuchInstance=False,
    IsNoSuchObject=False,
    IsEndOfMibView=False,
    BoolValue=False,
    IntValue=1337,
    FloatValue=0.0,
    ByteArrayValue=[],
    StringValue='',
)

_MULTI_RESULT_FLOAT = MultiResult(
    OID='.1.2.3.4',
    Type='float',
    IsNull=False,
    IsUnknown=False,
    IsNoSuchInstance=False,
    IsNoSuchObject=False,
    IsEndOfMibView=False,
    BoolValue=False,
    IntValue=0,
    FloatValue=1.337,
    ByteArrayValue=[],
    StringValue='',
)

_MULTI_RESULT_STRING = MultiResult(
    OID='.1.2.3.4',
    Type='string',
    IsNull=False,
    IsUnknown=False,
    IsNoSuchInstance=False,
    IsNoSuchObject=False,
    IsEndOfMibView=False,
    BoolValue=False,
    IntValue=0,
    FloatValue=0.0,
    ByteArrayValue=[],
    StringValue='some string',
)

_MULTI_RESULT_BYTEARRAY = MultiResult(
    OID='.1.2.3.4',
    Type='bytearray',
    IsNull=False,
    IsUnknown=False,
    IsNoSuchInstance=False,
    IsNoSuchObject=False,
    IsEndOfMibView=False,
    BoolValue=False,
    IntValue=0,
    FloatValue=0.0,
    ByteArrayValue=[0x41, 0x42, 0x43, 0x44, 0x45, 0x46],
    StringValue='',
)

_MULTI_RESULT_GARBAGE = MultiResult(
    OID='.1.2.3.4',
    Type='ham sandwich',
    IsNull=True,
    IsUnknown=False,
    IsNoSuchInstance=False,
    IsNoSuchObject=False,
    IsEndOfMibView=False,
    BoolValue=False,
    IntValue=0,
    FloatValue=0.0,
    ByteArrayValue=[],
    StringValue='',
)

_SNMP_VARIABLE = SNMPVariable(
    oid='.1.2.3',
    oid_index=4,
    snmp_type='string',
    value='some value',
)


class CommonTest(unittest.TestCase):
    def test_handle_multi_result(self):
        assert_that(
            handle_multi_result(
                _MULTI_RESULT_NOSUCHINSTANCE
            ),
            equal_to(
                SNMPVariable(oid='.1.2.3', oid_index=4, snmp_type='noSuchInstance', value=None)
            )
        )

        assert_that(
            handle_multi_result(
                _MULTI_RESULT_NOSUCHOBJECT
            ),
            equal_to(
                SNMPVariable(oid='.1.2.3', oid_index=4, snmp_type='noSuchObject', value=None)
            )
        )

        assert_that(
            handle_multi_result(
                _MULTI_RESULT_ENDOFMIBVIEW
            ),
            equal_to(
                SNMPVariable(oid='.1.2.3', oid_index=4, snmp_type='endOfMibView', value=None)
            )
        )

        assert_that(
            handle_multi_result(
                _MULTI_RESULT_BOOL
            ),
            equal_to(
                SNMPVariable(oid='.1.2.3', oid_index=4, snmp_type='bool', value=True)
            )
        )

        assert_that(
            handle_multi_result(
                _MULTI_RESULT_INT
            ),
            equal_to(
                SNMPVariable(oid='.1.2.3', oid_index=4, snmp_type='int', value=1337)
            )
        )

        assert_that(
            handle_multi_result(
                _MULTI_RESULT_FLOAT
            ),
            equal_to(
                SNMPVariable(oid='.1.2.3', oid_index=4, snmp_type='float', value=1.337)
            )
        )

        assert_that(
            handle_multi_result(
                _MULTI_RESULT_BYTEARRAY
            ),
            equal_to(
                SNMPVariable(oid='.1.2.3', oid_index=4, snmp_type='bytearray', value='ABCDEF')
            )
        )

        assert_that(
            handle_multi_result(
                _MULTI_RESULT_STRING
            ),
            equal_to(
                SNMPVariable(oid='.1.2.3', oid_index=4, snmp_type='string', value='some string')
            )
        )

        assert_that(
            calling(handle_multi_result).with_args(
                _MULTI_RESULT_GARBAGE
            ),
            raises(UnknownSNMPTypeError)
        )

    def test_handle_multi_result_json(self):
        assert_that(
            handle_multi_result_json(_MULTI_RESULT_JSON_STRING),
            equal_to(_MULTI_RESULT_STRING)
        )
