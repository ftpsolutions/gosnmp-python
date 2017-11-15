import unittest

from hamcrest import assert_that, equal_to
from mock import MagicMock, call, patch

from common import SNMPVariable
from common_test import _SNMP_VARIABLE
from session import Session, create_snmpv1_session, create_snmpv2c_session, \
    create_snmpv3_session


class SessionTest(unittest.TestCase):
    def setUp(self):
        self._subject = Session(
            session=MagicMock(),
        )

    @patch('gosnmp_python.session.handle_multi_result')
    def test_get(self, handle_multi_result):
        handle_multi_result.return_value = _SNMP_VARIABLE

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

    @patch('gosnmp_python.session.handle_multi_result')
    def test_get_next(self, handle_multi_result):
        handle_multi_result.return_value = _SNMP_VARIABLE

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


class ConstructorsTest(unittest.TestCase):
    @patch('gosnmp_python.session.Session')
    @patch('gosnmp_python.session._new_session_v1')
    def test_create_snmpv1_session(self, go_session_constructor, py_session_constructor):
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
                call(session=go_session_constructor())
            ])
        )

        assert_that(
            subject,
            equal_to(
                py_session_constructor()
            )
        )

    @patch('gosnmp_python.session.Session')
    @patch('gosnmp_python.session._new_session_v2c')
    def test_create_snmpv2c_session(self, go_session_constructor, py_session_constructor):
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
                call(session=go_session_constructor())
            ])
        )

        assert_that(
            subject,
            equal_to(
                py_session_constructor()
            )
        )

    @patch('gosnmp_python.session.Session')
    @patch('gosnmp_python.session._new_session_v3')
    def test_create_snmpv3_session(self, go_session_constructor, py_session_constructor):
        subject = create_snmpv3_session(
            hostname=u'some_hostname',
            security_username=u'some_username',
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
                call(session=go_session_constructor())
            ])
        )

        assert_that(
            subject,
            equal_to(
                py_session_constructor()
            )
        )
