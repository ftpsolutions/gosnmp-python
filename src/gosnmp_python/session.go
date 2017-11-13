package gosnmp_python

// #cgo pkg-config: python2
// #include <Python.h>
import "C" // this has to follow the above comments with no spaces

import (
	"time"
	"math"
	"github.com/initialed85/gosnmp"
	"fmt"
	"strings"
	"errors"
)

//
// structs
//

type Session struct {
	snmp *gosnmp.GoSNMP // this has to be private
}

type MultiResult struct {
	OID              string
	Type             string
	IsNull           bool
	IsUnknown        bool
	IsNoSuchInstance bool
	IsNoSuchObject   bool
	IsEndOfMibView   bool
	BoolValue        bool
	IntValue         int
	FloatValue       float64
	ByteArray        []byte
	StringValue      string
}

//
// helper methods
//

func getSecurityLevel(securityLevel string) gosnmp.SnmpV3MsgFlags {
	securityLevel = strings.ToLower(securityLevel)
	actualSecurityLevel := gosnmp.NoAuthNoPriv

	switch securityLevel {

	case "authnopriv":
		actualSecurityLevel = gosnmp.AuthNoPriv
	case "authpriv":
		actualSecurityLevel = gosnmp.AuthPriv

	}

	return actualSecurityLevel
}

func getPrivacyDetails(privacyPassword, privacyProtocol string) (string, gosnmp.SnmpV3PrivProtocol) {
	if privacyProtocol == "" {
		privacyPassword = ""
	}

	actualPrivacyProtocol := gosnmp.NoPriv

	switch privacyProtocol {

	case "DES":
		actualPrivacyProtocol = gosnmp.DES
	case "AES":
		actualPrivacyProtocol = gosnmp.AES

	}

	return privacyPassword, actualPrivacyProtocol
}

func getAuthenticationDetails(AuthenticationPassword, AuthenticationProtocol string) (string, gosnmp.SnmpV3AuthProtocol) {
	if AuthenticationProtocol == "" {
		AuthenticationPassword = ""
	}

	actualAuthenticationProtocol := gosnmp.NoAuth

	switch AuthenticationProtocol {

	case "MD5":
		actualAuthenticationProtocol = gosnmp.MD5
	case "SHA":
		actualAuthenticationProtocol = gosnmp.SHA

	}

	return AuthenticationPassword, actualAuthenticationProtocol
}

func buildMultiResult(oid string, valueType gosnmp.Asn1BER, value interface{}) (MultiResult, error) {
	multiResult := MultiResult{
		OID:              oid,
		Type:             "",
		IsNull:           false,
		IsUnknown:        false,
		IsNoSuchInstance: false,
		IsNoSuchObject:   false,
		IsEndOfMibView:   false,
		BoolValue:        false,
		IntValue:         0,
		FloatValue:       0.0,
		ByteArray:        []byte{},
		StringValue:      "",
	}

	switch valueType {

	case gosnmp.Null:
		fallthrough
	case gosnmp.UnknownType:
		fallthrough
	case gosnmp.NoSuchInstance:
		multiResult.Type = "noSuchInstance"
		multiResult.IsNoSuchInstance = true
		return multiResult, nil

	case gosnmp.NoSuchObject:
		multiResult.Type = "noSuchObject"
		multiResult.IsNoSuchObject = true
		return multiResult, nil

	case gosnmp.EndOfMibView:
		multiResult.Type = "endOfMibView"
		multiResult.IsEndOfMibView = true
		return multiResult, nil

	case gosnmp.Boolean:
		multiResult.Type = "bool"
		multiResult.BoolValue = value.(bool)
		return multiResult, nil

	case gosnmp.Counter32:
		fallthrough
	case gosnmp.Gauge32:
		fallthrough
	case gosnmp.Uinteger32:
		multiResult.Type = "int"
		multiResult.IntValue = int(value.(uint))
		return multiResult, nil

	case gosnmp.Counter64:
		multiResult.Type = "int"
		multiResult.IntValue = int(value.(uint64))
		return multiResult, nil

	case gosnmp.Integer:
		fallthrough
	case gosnmp.TimeTicks:
		multiResult.Type = "int"
		multiResult.IntValue = value.(int)
		return multiResult, nil

	case gosnmp.Opaque:
		multiResult.Type = "float"
		multiResult.FloatValue = value.(float64)
		return multiResult, nil

	case gosnmp.OctetString:
		multiResult.Type = "bytearray"
		multiResult.ByteArray = value.([]byte)
		return multiResult, nil

	case gosnmp.ObjectIdentifier:
		fallthrough
	case gosnmp.IPAddress:
		multiResult.Type = "string"
		multiResult.StringValue = value.(string)
		return multiResult, nil

	}

	return multiResult, errors.New(
		fmt.Sprintf("Unknown type; oid=%v, type=%v, value=%v", oid, valueType, value),
	)
}

//
// public methods
//

func (self *Session) Connect() error {
	tState := C.PyEval_SaveThread()
	defer C.PyEval_RestoreThread(tState)

	return self.snmp.Connect()
}

func (self *Session) Get(oid string) (MultiResult, error) {
	tState := C.PyEval_SaveThread()
	defer C.PyEval_RestoreThread(tState)

	emptyMultiResult := MultiResult{}

	result, err := self.snmp.Get([]string{oid})
	if err != nil {
		return emptyMultiResult, err
	}

	multiResult, err := buildMultiResult(
		result.Variables[0].Name,
		result.Variables[0].Type,
		result.Variables[0].Value,
	)
	if err != nil {
		return emptyMultiResult, err
	}

	return multiResult, nil
}

func (self *Session) GetNext(oid string) (MultiResult, error) {
	tState := C.PyEval_SaveThread()
	defer C.PyEval_RestoreThread(tState)

	emptyMultiResult := MultiResult{}

	result, err := self.snmp.GetNext([]string{oid})
	if err != nil {
		return emptyMultiResult, err
	}

	multiResult, err := buildMultiResult(
		result.Variables[0].Name,
		result.Variables[0].Type,
		result.Variables[0].Value,
	)
	if err != nil {
		return emptyMultiResult, err
	}

	return multiResult, nil
}

func (self *Session) Close() error {
	tState := C.PyEval_SaveThread()
	defer C.PyEval_RestoreThread(tState)

	return self.snmp.Conn.Close()
}

//
// constructors
//

func NewSessionV1(hostname string, port int, community string, timeout, retries int) Session {
	tState := C.PyEval_SaveThread()
	defer C.PyEval_RestoreThread(tState)

	snmp := &gosnmp.GoSNMP{
		Target:    hostname,
		Port:      uint16(port),
		Community: community,
		Version:   gosnmp.Version1,
		Timeout:   time.Duration(timeout) * time.Second,
		Retries:   retries,
		MaxOids:   math.MaxInt32,
	}

	self := Session{
		snmp: snmp,
	}

	return self
}

func NewSessionV2c(hostname string, port int, community string, timeout, retries int) (Session) {
	tState := C.PyEval_SaveThread()
	defer C.PyEval_RestoreThread(tState)

	snmp := &gosnmp.GoSNMP{
		Target:    hostname,
		Port:      uint16(port),
		Community: community,
		Version:   gosnmp.Version1,
		Timeout:   time.Duration(timeout) * time.Second,
		Retries:   retries,
		MaxOids:   math.MaxInt32,
	}

	self := Session{
		snmp: snmp,
	}

	return self
}

func NewSessionV3(hostname string, port int, securityUsername, privacyPassword, authPassword, securityLevel, authProtocol, privacyProtocol string, timeout, retries int) Session {
	tState := C.PyEval_SaveThread()
	defer C.PyEval_RestoreThread(tState)

	actualAuthPassword, actualAuthProtocol := getAuthenticationDetails(authPassword, authProtocol)
	actualPrivPassword, actualPrivProtocol := getPrivacyDetails(privacyPassword, privacyProtocol)

	snmp := &gosnmp.GoSNMP{
		Target:        hostname,
		Port:          uint16(port),
		Version:       gosnmp.Version3,
		Timeout:       time.Duration(timeout) * time.Second,
		Retries:       retries,
		SecurityModel: gosnmp.UserSecurityModel,
		MsgFlags:      getSecurityLevel(securityLevel),
		SecurityParameters: &gosnmp.UsmSecurityParameters{
			UserName:                 securityUsername,
			AuthenticationProtocol:   actualAuthProtocol,
			AuthenticationPassphrase: actualAuthPassword,
			PrivacyProtocol:          actualPrivProtocol,
			PrivacyPassphrase:        actualPrivPassword,
		},
		MaxOids: math.MaxInt32,
	}

	self := Session{
		snmp: snmp,
	}

	return self
}
