package gosnmp_python

import (
	"encoding/json"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/ftpsolutions/gosnmp"
)

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
	ByteArrayValue   []int
	StringValue      string
}

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
		OID: oid,
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
		multiResult.Type = "int"
		multiResult.IntValue = value.(int)
		return multiResult, nil
	case gosnmp.TimeTicks:
		multiResult.Type = "int"
		multiResult.IntValue = int(value.(uint))
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

		multiResult.ByteArrayValue = valueAsInts
		return multiResult, nil

	case gosnmp.ObjectIdentifier:
		fallthrough
	case gosnmp.IPAddress:
		multiResult.Type = "string"
		multiResult.StringValue = value.(string)
		return multiResult, nil

	}

	return multiResult, fmt.Errorf("Unknown type; oid=%v, type=%v, value=%v", oid, valueType, value)
}

type sessionInterface interface {
	getSNMP() *gosnmp.GoSNMP
	connect() error
	get(string) (multiResult, error)
	getJSON(string) (string, error)
	getNext(string) (multiResult, error)
	getNextJSON(string) (string, error)
	close() error
}

type session struct {
	snmp      wrappedSNMPInterface
	connected bool // used to avoid weird memory errors if the underlying connect fails (snmp object left in insane state)
}

func newSessionV1(hostname string, port int, community string, timeout, retries int) session {
	snmp := wrappedSNMP{
		&gosnmp.GoSNMP{
			Target:    hostname,
			Port:      uint16(port),
			Community: community,
			Version:   gosnmp.Version1,
			Timeout:   time.Duration(timeout) * time.Second,
			Retries:   retries,
			MaxOids:   math.MaxInt32,
		},
	}

	s := session{
		snmp: &snmp,
	}

	return s
}

func newSessionV2c(hostname string, port int, community string, timeout, retries int) session {
	snmp := wrappedSNMP{
		&gosnmp.GoSNMP{
			Target:    hostname,
			Port:      uint16(port),
			Community: community,
			Version:   gosnmp.Version2c,
			Timeout:   time.Duration(timeout) * time.Second,
			Retries:   retries,
			MaxOids:   math.MaxInt32,
		},
	}

	s := session{
		snmp: &snmp,
	}

	return s
}

func newSessionV3(hostname string, port int, contextName, securityUsername, privacyPassword, authPassword, securityLevel, authProtocol, privacyProtocol string, timeout, retries int) session {
	actualAuthPassword, actualAuthProtocol := getAuthenticationDetails(authPassword, authProtocol)
	actualPrivPassword, actualPrivProtocol := getPrivacyDetails(privacyPassword, privacyProtocol)

	snmp := wrappedSNMP{
		&gosnmp.GoSNMP{
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
			MaxOids:     math.MaxInt32,
			ContextName: contextName,
		},
	}

	s := session{
		snmp: &snmp,
	}

	return s
}

func (s *session) getSNMP() *gosnmp.GoSNMP {
	return s.snmp.getSNMP()
}

func (s *session) connect() error {
	if s.connected {
		return nil
	}

	err := s.snmp.connect()

	s.connected = err == nil

	return err
}

func (s *session) get(oid string) (multiResult, error) {
	emptyMultiResult := multiResult{}

	result, err := s.snmp.get([]string{oid})
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

func (s *session) getJSON(oid string) (string, error) {
	result, err := s.snmp.get([]string{oid})
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

func (s *session) getNext(oid string) (multiResult, error) {
	emptyMultiResult := multiResult{}

	result, err := s.snmp.getNext([]string{oid})
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

func (s *session) getNextJSON(oid string) (string, error) {
	result, err := s.snmp.getNext([]string{oid})
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

func (s *session) close() error {
	if s.snmp != nil {
		if s.snmp.getConn() != nil {
			s.snmp.close()
			s.snmp = nil
		}
	}

	s.connected = false

	return nil
}
