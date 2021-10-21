import json
import re
from threading import RLock

from gosnmp_python.built.gosnmp_python_go import (
    NewRPCSessionV1,
    NewRPCSessionV2c,
    NewRPCSessionV3,
    RPCConnect,
    RPCGet,
    RPCGetNext,
    RPCGetBulk,
    RPCWalk,
    RPCWalkBulk,
    RPCSetInteger,
    RPCSetIPAddress,
    RPCSetString,
    RPCClose,
)
from gosnmp_python.common import handle_exception, handle_multi_result, handle_multi_result_json

_new_session_lock = RLock()


def _new_rpc_session_v1(*args):
    with _new_session_lock:
        return handle_exception(NewRPCSessionV1, args)


def _new_rpc_session_v2c(*args):
    with _new_session_lock:
        return handle_exception(NewRPCSessionV2c, args)


def _new_rpc_session_v3(*args):
    with _new_session_lock:
        return handle_exception(NewRPCSessionV3, args)


_IP_ADDRESS = re.compile(r"^[0-2]?[0-9]?[0-9]{1}\.[0-2]?[0-9]?[0-9]{1}\.[0-2]?[0-9]?[0-9]{1}\.[0-2]?[0-9]?[0-9]{1}$")

_V1 = "v1"
_V2C = "v2c"
_V3 = "v3"


class RPCSession(object):
    def __init__(self, session_id, version, **kwargs):
        self._session_id = session_id
        self._version = version
        self._kwargs = kwargs

    def __del__(self):
        try:
            self.close()
        except BaseException:
            pass

    def __repr__(self):
        return "{0}(session_id={1}, {2})".format(
            self.__class__.__name__, repr(self._session_id), ", ".join("{0}={1}".format(k, repr(v)) for k, v in list(self._kwargs.items()))
        )

    def connect(self):
        return handle_exception(RPCConnect, (self._session_id,), self)

    def get(self, oid):
        oid = str(oid)

        return handle_multi_result(
            handle_multi_result_json(
                handle_exception(RPCGet, (self._session_id, oid), self),
                self,
            )
        )

    def get_next(self, oid):
        oid = str(oid)

        return handle_multi_result(
            handle_multi_result_json(
                handle_exception(RPCGetNext, (self._session_id, oid), self),
                self,
            ),
        )

    def get_bulk(self, oids, non_repeaters, max_repetitions):
        if self._version == _V1:
            raise NotImplementedError("cannot call GETBULK with SNMPv1")

        if not isinstance(oids, (list, tuple)):
            oids = [oids]

        # TODO: fix this hack- gopy not happy receiving lists
        oids = json.dumps(oids)

        return handle_multi_result(
            handle_multi_result_json(
                handle_exception(RPCGetBulk, (self._session_id, oids, non_repeaters, max_repetitions), self),
                self,
            ),
        )

    def walk(self, oid):
        oid = str(oid)

        return handle_multi_result(
            handle_multi_result_json(
                handle_exception(RPCWalk, (self._session_id, oid), self),
                self,
            ),
        )

    def walk_bulk(self, oid):
        if self._version == _V1:
            raise NotImplementedError("cannot call BULKWALK with SNMPv1")

        oid = str(oid)

        return handle_multi_result(
            handle_multi_result_json(
                handle_exception(RPCWalkBulk, (self._session_id, oid), self),
                self,
            ),
        )

    def set(self, oid, value, is_ip_address=None):
        if not isinstance(value, (int, str)):
            raise TypeError("gosnmp_python only supports SNMP set for integers and strings")

        if isinstance(value, str):
            if is_ip_address is True or is_ip_address is None and _IP_ADDRESS.match(value) is not None:
                method = RPCSetIPAddress
            else:
                method = RPCSetString
        else:
            method = RPCSetInteger

        return handle_multi_result(
            handle_multi_result_json(
                handle_exception(method, (self._session_id, oid, value), self),
                self,
            ),
        )

    def close(self):
        return handle_exception(RPCClose, (self._session_id,), self)


def create_snmpv1_session(hostname, community, port=161, timeout=5, retries=1):
    session_id = _new_rpc_session_v1(
        str(hostname),
        int(port),
        str(community),
        int(timeout),
        int(retries),
    )

    kwargs = {
        "hostname": hostname,
        "community": community,
        "port": port,
        "timeout": timeout,
        "retries": retries,
    }

    return RPCSession(session_id=session_id, version=_V1, **kwargs)


def create_snmpv2c_session(hostname, community, port=161, timeout=5, retries=1):
    session_id = _new_rpc_session_v2c(
        str(hostname),
        int(port),
        str(community),
        int(timeout),
        int(retries),
    )

    kwargs = {
        "hostname": hostname,
        "community": community,
        "port": port,
        "timeout": timeout,
        "retries": retries,
    }

    return RPCSession(session_id=session_id, version=_V2C, **kwargs)


def create_snmpv3_session(
    hostname,
    security_username,
    security_level,
    auth_password,
    auth_protocol,
    privacy_password,
    privacy_protocol,
    context_name=None,
    port=161,
    timeout=5,
    retries=1,
):
    context_name = context_name if context_name is not None else ""

    session_id = _new_rpc_session_v3(
        str(hostname),
        int(port),
        str(context_name),
        str(security_username),
        str(privacy_password),
        str(auth_password),
        str(security_level),
        str(auth_protocol),
        str(privacy_protocol),
        int(timeout),
        int(retries),
    )

    kwargs = {
        "hostname": hostname,
        "security_username": security_username,
        "security_level": security_level,
        "auth_password": auth_password,
        "auth_protocol": auth_protocol,
        "privacy_password": privacy_password,
        "privacy_protocol": privacy_protocol,
        "context_name": context_name,
        "port": port,
        "timeout": timeout,
        "retries": retries,
    }

    return RPCSession(session_id=session_id, version=_V3, **kwargs)
