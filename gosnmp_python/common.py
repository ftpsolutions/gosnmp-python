from collections import namedtuple

SNMPVariable = namedtuple('SNMPVariable', ['oid', 'oid_index', 'snmp_type', 'value'])

class UnknownSNMPTypeError(Exception):
    pass


class GoRuntimeError(Exception):
    pass


