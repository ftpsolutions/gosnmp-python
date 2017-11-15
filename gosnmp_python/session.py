from sys import version as python_version

from common import handle_exception, handle_multi_result

if 'pypy' not in python_version.strip().lower():
    from py2.gosnmp_python import NewSessionV1, NewSessionV2c, NewSessionV3
else:
    from cffi.gosnmp_python import NewSessionV1, NewSessionV2c, NewSessionV3


def _new_session_v1(*args):
    return NewSessionV1(*args)


def _new_session_v2c(*args):
    return NewSessionV2c(*args)


def _new_session_v3(*args):
    return NewSessionV3(*args)


class Session(object):
    def __init__(self, session):
        self._session = session

    def __repr__(self):
        return '{0}(session={1})'.format(
            self.__class__.__name__,
            repr(self._session)
        )

    def connect(self):
        return self._session.Connect()

    def get(self, oid):
        return handle_multi_result(
            handle_exception(self._session.Get, (oid,))
        )

    def get_next(self, oid):
        return handle_multi_result(
            handle_exception(self._session.GetNext, (oid,))
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
