import json
from collections import namedtuple

SNMPVariable = namedtuple("SNMPVariable", ["oid", "oid_index", "snmp_type", "value"])

MultiResult = namedtuple(
    "MultiResult",
    [
        "OID",
        "Type",
        "IsNull",
        "IsUnknown",
        "IsNoSuchInstance",
        "IsNoSuchObject",
        "IsEndOfMibView",
        "BoolValue",
        "IntValue",
        "FloatValue",
        "ByteArrayValue",
        "StringValue",
    ],
)


class UnknownSNMPTypeError(Exception):
    pass


class GoRuntimeError(Exception):
    pass


def handle_exception(method, args, other=None):
    try:
        return method(*args)
    except Exception as e:
        e = e if not hasattr(e, "__cause__") and isinstance(e.__cause__, BaseException) else e.__cause__

        new_args = ("attempt to call {} on the Go side raised {}".format("{}(*{})".format(repr(method), repr(args)), repr(e)),)

        e.args = new_args

        raise GoRuntimeError(e)


def handle_multi_result_json(multi_result_json_string, session=None):
    if multi_result_json_string == "[]":
        return []

    try:
        multi_result_json = json.loads(multi_result_json_string)
    except ValueError as e:
        raise ValueError("{0} raised {1} while parsing {2}".format(session, e, repr(multi_result_json_string)))

    if isinstance(multi_result_json, (list, tuple)):
        return [MultiResult(**x) for x in multi_result_json]

    return MultiResult(**multi_result_json)


def _handle_multi_result(multi_result):
    raw_oid = multi_result.OID.strip(". ")

    oid = ".{0}".format(".".join(raw_oid.split(".")[0:-1]).strip("."))
    oid_index = int(raw_oid.split(".")[-1])

    if multi_result.Type in ["noSuchInstance", "noSuchObject", "endOfMibView"]:
        return SNMPVariable(
            oid=oid,
            oid_index=oid_index,
            snmp_type=multi_result.Type,
            value=None,
        )
    elif multi_result.Type in ["bool"]:
        return SNMPVariable(
            oid=oid,
            oid_index=oid_index,
            snmp_type=multi_result.Type,
            value=multi_result.BoolValue,
        )
    elif multi_result.Type in ["int"]:
        return SNMPVariable(
            oid=oid,
            oid_index=oid_index,
            snmp_type=multi_result.Type,
            value=multi_result.IntValue,
        )
    elif multi_result.Type in ["float"]:
        return SNMPVariable(
            oid=oid,
            oid_index=oid_index,
            snmp_type=multi_result.Type,
            value=multi_result.FloatValue,
        )
    elif multi_result.Type in ["bytearray"]:
        return SNMPVariable(
            oid=oid,
            oid_index=oid_index,
            snmp_type=multi_result.Type,
            value="".join([chr(x) for x in multi_result.ByteArrayValue]),
        )
    elif multi_result.Type in ["string"]:
        return SNMPVariable(
            oid=oid,
            oid_index=oid_index,
            snmp_type=multi_result.Type,
            value=multi_result.StringValue,
        )

    raise UnknownSNMPTypeError("{0} represents an unknown SNMP type".format(multi_result))


def handle_multi_result(multi_result_or_multi_results):
    if not isinstance(multi_result_or_multi_results, MultiResult):
        return [_handle_multi_result(x) for x in multi_result_or_multi_results]

    return _handle_multi_result(multi_result_or_multi_results)
