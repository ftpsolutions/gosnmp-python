package gosnmp_python

// #cgo pkg-config: python2
// #include <Python.h>

import (
	"time"
	"math"
	"github.com/initialed85/gosnmp"
	"fmt"
	"strings"
	"errors"
	"encoding/json"
)

// structs

type session struct {
	snmp *gosnmp.GoSNMP // this has to be private
}

type multiResult struct {
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
	ByteArray        []int
	StringValue      string
}

// helper functions

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

func buildMultiResult(oid string, valueType gosnmp.Asn1BER, value interface{}) (multiResult, error) {
	multiResult := multiResult{
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
		ByteArray:        []int{},
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

		valueAsBytes := value.([]byte)
		valueAsInts := make([]int, len(valueAsBytes), len(valueAsBytes))

		for i, c := range valueAsBytes {
			valueAsInts[i] = int(c)
		}

		multiResult.ByteArray = valueAsInts
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

// public methods

func (self *session) connect() error {
	return self.snmp.Connect()
}

func (self *session) get(oid string) (multiResult, error) {
	emptyMultiResult := multiResult{}

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

func (self *session) getJSON(oid string) (string, error) {
	result, err := self.snmp.Get([]string{oid})
	if err != nil {
		return "{}", err
	}

	multiResult, err := buildMultiResult(
		result.Variables[0].Name,
		result.Variables[0].Type,
		result.Variables[0].Value,
	)
	if err != nil {
		return "{}", err
	}

	multiResultBytes, err := json.Marshal(multiResult)
	if err != nil {
		return "{}", err
	}

	return string(multiResultBytes), nil
}

func (self *session) getNext(oid string) (multiResult, error) {
	emptyMultiResult := multiResult{}

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

func (self *session) getNextJSON(oid string) (string, error) {
	result, err := self.snmp.GetNext([]string{oid})
	if err != nil {
		return "{}", err
	}

	multiResult, err := buildMultiResult(
		result.Variables[0].Name,
		result.Variables[0].Type,
		result.Variables[0].Value,
	)
	if err != nil {
		return "{}", err
	}

	multiResultBytes, err := json.Marshal(multiResult)
	if err != nil {
		return "{}", err
	}

	return string(multiResultBytes), nil
}

func (self *session) close() error {
	err := self.snmp.Conn.Close()

	self.snmp = nil // seems to be important for Go's garbage collector

	return err
}

// constructors

func newSessionV1(hostname string, port int, community string, timeout, retries int) session {
	snmp := &gosnmp.GoSNMP{
		Target:    hostname,
		Port:      uint16(port),
		Community: community,
		Version:   gosnmp.Version1,
		Timeout:   time.Duration(timeout) * time.Second,
		Retries:   retries,
		MaxOids:   math.MaxInt32,
	}

	self := session{
		snmp: snmp,
	}

	return self
}

func newSessionV2c(hostname string, port int, community string, timeout, retries int) (session) {
	snmp := &gosnmp.GoSNMP{
		Target:    hostname,
		Port:      uint16(port),
		Community: community,
		Version:   gosnmp.Version2c,
		Timeout:   time.Duration(timeout) * time.Second,
		Retries:   retries,
		MaxOids:   math.MaxInt32,
	}

	self := session{
		snmp: snmp,
	}

	return self
}

func newSessionV3(hostname string, port int, securityUsername, privacyPassword, authPassword, securityLevel, authProtocol, privacyProtocol string, timeout, retries int) session {
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

	self := session{
		snmp: snmp,
	}

	return self
}
