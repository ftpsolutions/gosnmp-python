from common import handle_exception, handle_multi_result
from gosnmp_python import NewRPCSessionV1, NewRPCSessionV2c, NewRPCSessionV3, RPCConnect, RPCGet, RPCGetNext, \
    RPCClose


def _new_rpc_session_v1(*args):
    return NewRPCSessionV1(*args)


def _new_rpc_session_v2c(*args):
    return NewRPCSessionV2c(*args)


def _new_rpc_session_v3(*args):
    return NewRPCSessionV3(*args)


class RPCSession(object):
    def __init__(self, session_id):
        self._session_id = session_id

    def __repr__(self):
        return '{0}(session_id={1})'.format(
            self.__class__.__name__,
            repr(self._session_id)
        )

    def connect(self):
        return RPCConnect(self._session_id)

    def get(self, oid):
        return handle_multi_result(
            handle_exception(RPCGet, (self._session_id, oid,))
        )

    def get_next(self, oid):
        return handle_multi_result(
            handle_exception(RPCGetNext, (self._session_id, oid,))
        )

    def close(self):
        return RPCClose(self._session_id)


def create_snmpv1_session(hostname, community, port=161, timeout=5, retries=1):
    session_id = _new_rpc_session_v1(
        str(hostname),
        int(port),
        str(community),
        int(timeout),
        int(retries),
    )

    return RPCSession(
        session_id=session_id,
    )


def create_snmpv2c_session(hostname, community, port=161, timeout=5, retries=1):
    session_id = _new_rpc_session_v2c(
        str(hostname),
        int(port),
        str(community),
        int(timeout),
        int(retries),
    )

    return RPCSession(
        session_id=session_id,
    )


def create_snmpv3_session(hostname, security_username, security_level, auth_password, auth_protocol, privacy_password,
                          privacy_protocol, port=161, timeout=5, retries=1):
    session_id = _new_rpc_session_v3(
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

    return RPCSession(
        session_id=session_id,
    )
