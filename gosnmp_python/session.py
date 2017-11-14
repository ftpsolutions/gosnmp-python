from collections import namedtuple

from gosnmp_python import NewSessionV1, NewSessionV2c, NewSessionV3

SNMPVariable = namedtuple('SNMPVariable', ['oid', 'oid_index', 'snmp_type', 'value'])

_new_session_v1 = lambda *args: NewSessionV1(*args)

_new_session_v2c = lambda *args: NewSessionV2c(*args)

_new_session_v3 = lambda *args: NewSessionV3(*args)


class UnknownSNMPTypeError(Exception):
    pass


class GoRuntimeError(Exception):
    pass


class Session(object):
    def __init__(self, session):
        self._session = session

    def connect(self):
        return self._session.Connect()

    @staticmethod
    def _handle_multi_result(multi_result):
        raw_oid = multi_result.OID.strip('. ')

        oid = '.{0}'.format('.'.join(raw_oid.split('.')[0:-1]).strip('.'))
        oid_index = int(raw_oid.split('.')[-1])

        if multi_result.Type in ['noSuchInstance', 'noSuchObject', 'endOfMibView']:
            return SNMPVariable(
                oid=oid,
                oid_index=oid_index,
                snmp_type=multi_result.Type,
                value=None,
            )
        elif multi_result.Type in ['bool']:
            return SNMPVariable(
                oid=oid,
                oid_index=oid_index,
                snmp_type=multi_result.Type,
                value=multi_result.BoolValue,
            )
        elif multi_result.Type in ['int']:
            return SNMPVariable(
                oid=oid,
                oid_index=oid_index,
                snmp_type=multi_result.Type,
                value=multi_result.IntValue,
            )
        elif multi_result.Type in ['float']:
            return SNMPVariable(
                oid=oid,
                oid_index=oid_index,
                snmp_type=multi_result.Type,
                value=multi_result.FloatValue,
            )
        elif multi_result.Type in ['bytearray']:
            return SNMPVariable(
                oid=oid,
                oid_index=oid_index,
                snmp_type=multi_result.Type,
                value=''.join([chr(x) for x in multi_result.ByteArray]),
            )
        elif multi_result.Type in ['string']:
            return SNMPVariable(
                oid=oid,
                oid_index=oid_index,
                snmp_type=multi_result.Type,
                value=multi_result.StringValue,
            )

        raise UnknownSNMPTypeError('{0} represents an unknown SNMP type'.format(
            multi_result
        ))

    def _handle_exception(self, method, oid):
        try:
            return method(oid)
        except RuntimeError as e:
            raise GoRuntimeError('{0} raised on Go side while calling {1} with oid {2}'.format(
                repr(e), method, repr(oid),
            ))

    def get(self, oid):
        return self._handle_multi_result(
            self._handle_exception(self._session.Get, oid)
        )

    def get_next(self, oid):
        return self._handle_multi_result(
            self._handle_exception(self._session.GetNext, oid)
        )

    def close(self):
        return self._session.Close()


def create_snmpv1_session(hostname, community, port=161, timeout=5, retries=1):
    session = _new_session_v1(
        str(hostname),
        int(port),
        str(community),
        int(timeout),
        int(retries),
    )

    return Session(
        session=session,
    )


def create_snmpv2c_session(hostname, community, port=161, timeout=5, retries=1):
    session = _new_session_v2c(
        str(hostname),
        int(port),
        str(community),
        int(timeout),
        int(retries),
    )

    return Session(
        session=session,
    )


def create_snmpv3_session(hostname, security_username, security_level, auth_password, auth_protocol, privacy_password,
                          privacy_protocol, port=161, timeout=5, retries=1):
    session = _new_session_v3(
        str(hostname),
        int(port),
        str(security_username),
        str(privacy_password),
        str(auth_password),
        str(security_level),
        str(auth_protocol),
        str(privacy_protocol),
        int(timeout),
        int(retries),
    )

    return Session(
        session=session,
    )
