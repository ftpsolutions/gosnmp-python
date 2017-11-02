import unittest
from collections import namedtuple

from hamcrest import assert_that, equal_to, calling, raises
from mock import MagicMock, call

from session import Session, SNMPVariable, UnknownSNMPTypeError

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
    'ByteArray',
    'StringValue',
])

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
    ByteArray='[]byte{}',
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
    ByteArray='[]byte{}',
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
    ByteArray='[]byte{}',
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
    ByteArray='[]byte{}',
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
    ByteArray='[]byte{}',
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
    ByteArray='[]byte{}',
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
    ByteArray='[]byte{}',
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
    ByteArray='[]byte{0x41, 0x42, 0x43, 0x44, 0x45, 0x46}',
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
    ByteArray='[]byte{}',
    StringValue='',
)

_SNMP_VARIABLE = SNMPVariable(
    oid='.1.2.3',
    oid_index=4,
    snmp_type='string',
    value='some value',
)


class SessionTest(unittest.TestCase):
    def setUp(self):
        self._subject = Session(
            session=MagicMock(),
        )

    def test_handle_multi_result(self):
        assert_that(
            self._subject._handle_multi_result(
                _MULTI_RESULT_NOSUCHINSTANCE
            ),
            equal_to(
                SNMPVariable(oid='.1.2.3', oid_index=4, snmp_type='noSuchInstance', value=None)
            )
        )

        assert_that(
            self._subject._handle_multi_result(
                _MULTI_RESULT_NOSUCHOBJECT
            ),
            equal_to(
                SNMPVariable(oid='.1.2.3', oid_index=4, snmp_type='noSuchObject', value=None)
            )
        )

        assert_that(
            self._subject._handle_multi_result(
                _MULTI_RESULT_ENDOFMIBVIEW
            ),
            equal_to(
                SNMPVariable(oid='.1.2.3', oid_index=4, snmp_type='endOfMibView', value=None)
            )
        )

        assert_that(
            self._subject._handle_multi_result(
                _MULTI_RESULT_BOOL
            ),
            equal_to(
                SNMPVariable(oid='.1.2.3', oid_index=4, snmp_type='bool', value=True)
            )
        )

        assert_that(
            self._subject._handle_multi_result(
                _MULTI_RESULT_INT
            ),
            equal_to(
                SNMPVariable(oid='.1.2.3', oid_index=4, snmp_type='int', value=1337)
            )
        )

        assert_that(
            self._subject._handle_multi_result(
                _MULTI_RESULT_FLOAT
            ),
            equal_to(
                SNMPVariable(oid='.1.2.3', oid_index=4, snmp_type='float', value=1.337)
            )
        )

        assert_that(
            self._subject._handle_multi_result(
                _MULTI_RESULT_BYTEARRAY
            ),
            equal_to(
                SNMPVariable(oid='.1.2.3', oid_index=4, snmp_type='bytearray', value='ABCDEF')
            )
        )

        assert_that(
            self._subject._handle_multi_result(
                _MULTI_RESULT_STRING
            ),
            equal_to(
                SNMPVariable(oid='.1.2.3', oid_index=4, snmp_type='string', value='some string')
            )
        )

        assert_that(
            calling(self._subject._handle_multi_result).with_args(
                _MULTI_RESULT_GARBAGE
            ),
            raises(UnknownSNMPTypeError)
        )

    def test_get(self):
        self._subject._handle_multi_result = MagicMock()
        self._subject._handle_multi_result.return_value = _SNMP_VARIABLE

        assert_that(
            self._subject.get('1.2.3.4'),
            equal_to(
                SNMPVariable(oid='.1.2.3', oid_index=4, snmp_type='string', value='some value')
            )
        )

        assert_that(
            self._subject._session.mock_calls,
            equal_to([
                call.Get('1.2.3.4')
            ])
        )

    def test_get_next(self):
        self._subject._handle_multi_result = MagicMock()
        self._subject._handle_multi_result.return_value = _SNMP_VARIABLE

        assert_that(
            self._subject.get_next('1.2.3.3'),
            equal_to(
                SNMPVariable(oid='.1.2.3', oid_index=4, snmp_type='string', value='some value')
            )
        )

        assert_that(
            self._subject._session.mock_calls,
            equal_to([
                call.GetNext('1.2.3.3')
            ])
        )
