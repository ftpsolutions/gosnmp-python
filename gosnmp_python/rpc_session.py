from __future__ import (absolute_import, division, print_function,
                        unicode_literals)

import re
import os
from sys import version as python_version
from sys import version_info
from threading import RLock

from builtins import *
from builtins import object, str
from future import standard_library

from .common import (handle_exception, handle_multi_result,
                     handle_multi_result_json)

standard_library.install_aliases()

is_pypy = 'pypy' in python_version.lower()
version = version_info[:3]

# needed for CFFI under Python2
if version < (3, 0, 0):
    from past.types.oldstr import oldstr as str

if not is_pypy and version < (3, 0, 0):  # for Python2
    from .py2.gosnmp_python_go import SetPyPy, NewRPCSessionV1, NewRPCSessionV2c, NewRPCSessionV3, RPCConnect, RPCGet, \
        RPCGetNext, RPCSetInteger, RPCSetIPAddress, RPCSetString, RPCClose
else:  # for all versions of PyPy and also Python3
    raise ValueError('PyPy and Python3 is not supported. Waiting for gopy CFFI support')
#     from .cffi.gosnmp_python import SetPyPy, NewRPCSessionV1, NewRPCSessionV2c, NewRPCSessionV3, RPCConnect, RPCGet, \
#         RPCGetNext, RPCSetInteger, RPCSetIPAddress, RPCSetString, RPCClose
#
#     SetPyPy()
#
#     print('WARNING: PyPy or Python3 detected, will use CFFI- be prepared for very odd behaviour')

_new_session_lock = RLock()


def _new_rpc_session_v1(*args):
    with _new_session_lock:
        return handle_exception(
            NewRPCSessionV1, args
        )


def _new_rpc_session_v2c(*args):
    with _new_session_lock:
        return handle_exception(
            NewRPCSessionV2c, args
        )


def _new_rpc_session_v3(*args):
    with _new_session_lock:
        return handle_exception(
            NewRPCSessionV3, args
        )


_IP_ADDRESS = re.compile(
    '^[0-2]?[0-9]?[0-9]{1}\.[0-2]?[0-9]?[0-9]{1}\.[0-2]?[0-9]?[0-9]{1}\.[0-2]?[0-9]?[0-9]{1}$'
)


class RPCSession(object):
    def __init__(self, session_id, **kwargs):
        self._session_id = session_id
        self._kwargs = kwargs

    def __del__(self):
        try:
            self.close()
        except BaseException:
            pass

    def __repr__(self):
        return '{0}(session_id={1}, {2})'.format(
            self.__class__.__name__,
            repr(self._session_id),
            ', '.join('{0}={1}'.format(k, repr(v)) for k, v in list(self._kwargs.items()))
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

    def set(self, oid, value, is_ip_address=None):
        if not isinstance(value, (int, long, str, unicode)):
            raise TypeError('gosnmp_python only supports SNMP set for integers and strings')

        if isinstance(value, basestring):
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
        'hostname': hostname,
        'community': community,
        'port': port,
        'timeout': timeout,
        'retries': retries,
    }

    return RPCSession(
        session_id=session_id,
        **kwargs
    )


def create_snmpv2c_session(hostname, community, port=161, timeout=5, retries=1):
    session_id = _new_rpc_session_v2c(
        str(hostname),
        int(port),
        str(community),
        int(timeout),
        int(retries),
    )

    kwargs = {
        'hostname': hostname,
        'community': community,
        'port': port,
        'timeout': timeout,
        'retries': retries,
    }

    return RPCSession(
        session_id=session_id,
        **kwargs
    )


def create_snmpv3_session(hostname, security_username, security_level, auth_password, auth_protocol, privacy_password,
                          privacy_protocol, context_name=None, port=161, timeout=5, retries=1):
    context_name = context_name if context_name is not None else ''

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
        'hostname': hostname,
        'security_username': security_username,
        'security_level': security_level,
        'auth_password': auth_password,
        'auth_protocol': auth_protocol,
        'privacy_password': privacy_password,
        'privacy_protocol': privacy_protocol,
        'context_name': context_name,
        'port': port,
        'timeout': timeout,
        'retries': retries,
    }

    return RPCSession(
        session_id=session_id,
        **kwargs
    )
