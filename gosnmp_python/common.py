import json

from collections import namedtuple

SNMPVariable = namedtuple(
    'SNMPVariable', [
        'oid',
        'oid_index',
        'snmp_type',
        'value'
    ]
)

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


class UnknownSNMPTypeError(Exception):
    pass


class GoRuntimeError(Exception):
    pass


def handle_exception(method, args):
    try:
        return method(*args)
    except RuntimeError as e:
        raise GoRuntimeError('{0} raised on Go side while calling {1} with args {2}'.format(
            repr(e), method, repr(args),
        ))


def handle_multi_result_json(multi_result_json_string, session=None):
    try:
        multi_result_json = json.loads(multi_result_json_string)
    except ValueError as e:
        raise ValueError(
            '{0} raised {1} while parsing {2}'.format(
                session,
                e, 
                repr(multi_result_json_string)
                )
            )

    return MultiResult(**multi_result_json)


def handle_multi_result(multi_result):
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
            value=''.join([chr(x) for x in multi_result.ByteArrayValue]),
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
