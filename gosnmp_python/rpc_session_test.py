from __future__ import (absolute_import, division, print_function,
                        unicode_literals)

import gc
import unittest

from builtins import *
from future import standard_library
from hamcrest import assert_that, equal_to, greater_than_or_equal_to
from mock import call, patch

from .common import SNMPVariable
from .common_test import _SNMP_VARIABLE
from .rpc_session import (RPCSession, create_snmpv1_session,
                          create_snmpv2c_session, create_snmpv3_session)

standard_library.install_aliases()


class SessionTest(unittest.TestCase):
    def setUp(self):
        self._subject = RPCSession(
            session_id=0,
        )

    @patch('gosnmp_python.rpc_session.RPCConnect')
    def test_connect(self, rpc_call):
        rpc_call.return_value = None

        assert_that(
            self._subject.connect(),
            equal_to(None)
        )

        assert_that(
            rpc_call.mock_calls,
            equal_to([
                call(0)
            ])
        )

    @patch('gosnmp_python.rpc_session.RPCGet')
    @patch('gosnmp_python.rpc_session.handle_multi_result')
    @patch('gosnmp_python.rpc_session.handle_multi_result_json')
    def test_get(self, handle_multi_result_json, handle_multi_result, rpc_call):
        handle_multi_result.return_value = _SNMP_VARIABLE

        assert_that(
            self._subject.get('1.2.3.4'),
            equal_to(
                SNMPVariable(oid='.1.2.3', oid_index=4, snmp_type='string', value='some value')
            )
        )

        assert_that(
            rpc_call.mock_calls,
            equal_to([
                call.RPCGet(0, '1.2.3.4')
            ])
        )

    @patch('gosnmp_python.rpc_session.RPCGetNext')
    @patch('gosnmp_python.rpc_session.handle_multi_result')
    @patch('gosnmp_python.rpc_session.handle_multi_result_json')
    def test_get_next(self, handle_multi_result_json, handle_multi_result, rpc_call):
        handle_multi_result.return_value = _SNMP_VARIABLE

        assert_that(
            self._subject.get_next('1.2.3.3'),
            equal_to(
                SNMPVariable(oid='.1.2.3', oid_index=4, snmp_type='string', value='some value')
            )
        )

        assert_that(
            rpc_call.mock_calls,
            equal_to([
                call.RPCGetNext(0, '1.2.3.3')
            ])
        )

    @patch('gosnmp_python.rpc_session.RPCSetInteger')
    @patch('gosnmp_python.rpc_session.handle_multi_result')
    @patch('gosnmp_python.rpc_session.handle_multi_result_json')
    def test_set_integer(self, handle_multi_result_json, handle_multi_result, rpc_call):
        handle_multi_result.return_value = _SNMP_VARIABLE

        assert_that(
            self._subject.set('1.2.3.4', 1),
            equal_to(
                SNMPVariable(oid='.1.2.3', oid_index=4, snmp_type='string', value='some value')
            )
        )

        assert_that(
            rpc_call.mock_calls,
            equal_to([
                call.RPCSetInteger(0, '1.2.3.4', 1)
            ])
        )

    @patch('gosnmp_python.rpc_session.RPCSetString')
    @patch('gosnmp_python.rpc_session.handle_multi_result')
    @patch('gosnmp_python.rpc_session.handle_multi_result_json')
    def test_set_string(self, handle_multi_result_json, handle_multi_result, rpc_call):
        handle_multi_result.return_value = _SNMP_VARIABLE

        assert_that(
            self._subject.set('1.2.3.4', 'string'),
            equal_to(
                SNMPVariable(oid='.1.2.3', oid_index=4, snmp_type='string', value='some value')
            )
        )

        assert_that(
            rpc_call.mock_calls,
            equal_to([
                call.RPCSetString(0, '1.2.3.4', 'string')
            ])
        )

    @patch('gosnmp_python.rpc_session.RPCSetIPAddress')
    @patch('gosnmp_python.rpc_session.handle_multi_result')
    @patch('gosnmp_python.rpc_session.handle_multi_result_json')
    def test_set_ip_address(self, handle_multi_result_json, handle_multi_result, rpc_call):
        handle_multi_result.return_value = _SNMP_VARIABLE

        assert_that(
            self._subject.set('1.2.3.4', '1.2.3.3'),
            equal_to(
                SNMPVariable(oid='.1.2.3', oid_index=4, snmp_type='string', value='some value')
            )
        )

        assert_that(
            rpc_call.mock_calls,
            equal_to([
                call(0, '1.2.3.4', '1.2.3.3')
            ])
        )

    @patch('gosnmp_python.rpc_session.RPCClose')
    def test_close(self, rpc_call):
        rpc_call.return_value = None

        assert_that(
            self._subject.close(),
            equal_to(None)
        )

        assert_that(
            rpc_call.mock_calls,
            equal_to([
                call(0)
            ])
        )

    @patch('gosnmp_python.rpc_session.RPCClose')
    def test_del(self, rpc_call):
        rpc_call.return_value = None

        del (self._subject)

        gc.collect()

        assert_that(
            len(rpc_call.mock_calls),
            greater_than_or_equal_to(1)
        )


class ConstructorsTest(unittest.TestCase):
    @patch('gosnmp_python.rpc_session.RPCSession')
    @patch('gosnmp_python.rpc_session._new_rpc_session_v1')
    def test_create_snmpv1_session(self, go_session_constructor, py_session_constructor):
        go_session_constructor.return_value = -1

        subject = create_snmpv1_session(
            hostname=u'some_hostname',
            community=u'some_community',
            port='161',
            timeout='5',
            retries='1',
        )

        assert_that(
            go_session_constructor.mock_calls,
            equal_to([
                call('some_hostname', 161, 'some_community', 5, 1)
            ])
        )

        assert_that(
            py_session_constructor.mock_calls,
            equal_to([
                call(
                    community=u'some_community',
                    hostname=u'some_hostname',
                    port='161',
                    retries='1',
                    session_id=-1,
                    timeout='5')
            ])
        )

        assert_that(
            subject,
            equal_to(
                py_session_constructor()
            )
        )

    @patch('gosnmp_python.rpc_session.RPCSession')
    @patch('gosnmp_python.rpc_session._new_rpc_session_v2c')
    def test_create_snmpv2c_session(self, go_session_constructor, py_session_constructor):
        go_session_constructor.return_value = -1

        subject = create_snmpv2c_session(
            hostname=u'some_hostname',
            community=u'some_community',
            port='161',
            timeout='5',
            retries='1',
        )

        assert_that(
            go_session_constructor.mock_calls,
            equal_to([
                call('some_hostname', 161, 'some_community', 5, 1)
            ])
        )

        assert_that(
            py_session_constructor.mock_calls,
            equal_to([
                call(
                    community=u'some_community',
                    hostname=u'some_hostname',
                    port='161',
                    retries='1',
                    session_id=-1,
                    timeout='5')
            ])
        )

        assert_that(
            subject,
            equal_to(
                py_session_constructor()
            )
        )

    @patch('gosnmp_python.rpc_session.RPCSession')
    @patch('gosnmp_python.rpc_session._new_rpc_session_v3')
    def test_create_snmpv3_session(self, go_session_constructor, py_session_constructor):
        go_session_constructor.return_value = -1

        subject = create_snmpv3_session(
            hostname=u'some_hostname',
            security_username=u'some_username',
            context_name=u'some_context_name',
            security_level=u'authPriv',
            auth_password=u'some_password',
            auth_protocol=u'SHA',
            privacy_password=u'other_password',
            privacy_protocol=u'AES',
            port='161',
            timeout='5',
            retries='1',
        )

        assert_that(
            go_session_constructor.mock_calls,
            equal_to([
                call(
                    'some_hostname',
                    161,
                    'some_context_name',
                    'some_username',
                    'other_password',
                    'some_password',
                    'authPriv',
                    'SHA',
                    'AES',
                    5,
                    1
                )
            ])
        )

        assert_that(
            py_session_constructor.mock_calls,
            equal_to([
                call(
                    auth_password=u'some_password',
                    auth_protocol=u'SHA',
                    context_name=u'some_context_name',
                    hostname=u'some_hostname',
                    port='161',
                    privacy_password=u'other_password',
                    privacy_protocol=u'AES',
                    retries='1',
                    security_level=u'authPriv',
                    security_username=u'some_username',
                    session_id=-1,
                    timeout='5')
            ])
        )

        assert_that(
            subject,
            equal_to(
                py_session_constructor()
            )
        )
