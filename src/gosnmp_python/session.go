package gosnmp_python

import (
	"encoding/json"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/ftpsolutions/gosnmp"
	"log"
	"os"
	"strconv"
)

const (
	usmStatsBaseOID = ".1.3.6.1.6.3.15.1.1"
)

var usmStatsLookup = map[string]string{
	".1.3.6.1.6.3.15.1.1.1.0": "usmStatsUnsupportedSecLevels",
	".1.3.6.1.6.3.15.1.1.2.0": "usmStatsNotInTimeWindows",
	".1.3.6.1.6.3.15.1.1.3.0": "usmStatsUnknownUserNames",
	".1.3.6.1.6.3.15.1.1.4.0": "usmStatsUnknownEngineIDs",
	".1.3.6.1.6.3.15.1.1.5.0": "usmStatsWrongDigests",
	".1.3.6.1.6.3.15.1.1.6.0": "usmStatsDecryptionErrors",
}

var usmIgnoreTypes = map[gosnmp.Asn1BER]bool{
	gosnmp.EndOfContents:  false,
	gosnmp.NoSuchObject:   false,
	gosnmp.NoSuchInstance: false,
	gosnmp.EndOfMibView:   false,
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

	case gosnmp.OpaqueFloat:
		multiResult.Type = "float"
		multiResult.FloatValue = float64(value.(float32))
		return multiResult, nil
	case gosnmp.OpaqueDouble:
		fallthrough
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

func checkVariableCount(oid string, result *gosnmp.SnmpPacket) error {
	if len(result.Variables) != 1 {
		return fmt.Errorf(
			"get(%+v) returned %+v variables; this is unexpected (should have been 1)",
			oid,
			len(result.Variables),
		)
	}

	return nil
}

func translateUsmStats(oid string) string {
	val, ok := usmStatsLookup[oid]
	if !ok {
		return "unknown; no mapping for this OID"
	}

	return val
}

func checkForSNMPv3Issues(oid string, result *gosnmp.SnmpPacket) error {
	if result.Version != gosnmp.Version3 {
		return nil
	}

	_, ok := usmIgnoreTypes[result.Variables[0].Type]
	if ok {
		return nil
	}

	if strings.HasPrefix(result.Variables[0].Name, usmStatsBaseOID) {
		return fmt.Errorf(
			"requested %+v and got %+v (type %+v, value %+v); likely an SNMPv3 auth/priv issue- resolves to '%+v'",
			oid,
			result.Variables[0].Name,
			result.Variables[0].Type,
			result.Variables[0].Value,
			translateUsmStats(result.Variables[0].Name),
		)
	}

	return nil
}

type sessionInterface interface {
	getSNMP() *gosnmp.GoSNMP
	connect() error
	get(string) (multiResult, error)
	getJSON(string) (string, error)
	getNext(string) (multiResult, error)
	getNextJSON(string) (string, error)
	setString(string, string) (multiResult, error)
	setStringJSON(string, string) (string, error)
	setInteger(string, int) (multiResult, error)
	setIntegerJSON(string, int) (string, error)
	setIPAddress(string, string) (multiResult, error)
	setIPAddressJSON(string, string) (string, error)
	close() error
}

type session struct {
	snmp      wrappedSNMPInterface
	connected bool // used to avoid weird memory errors if the underlying connect fails (snmp object left in insane state)
}

func getLogger(snmpProtocol, hostname string, port int) *log.Logger {
	envDebug := os.Getenv("GOSNMP_PYTHON_DEBUG")
	if len(envDebug) <= 0 {
		return nil
	}

	debugEnabled, err := strconv.ParseBool(envDebug)
	if err != nil {
		return nil
	}

	if !debugEnabled {
		return nil
	}

	return log.New(
		os.Stdout,
		fmt.Sprintf("%v:%v:%v\t", snmpProtocol, hostname, port),
		0,
	)
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

	logger := getLogger("SNMPv1", hostname, port)
	if logger != nil {
		snmp.snmp.Logger = logger
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

	logger := getLogger("SNMPv2c", hostname, port)
	if logger != nil {
		snmp.snmp.Logger = logger
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
				AuthenticationPassphrase: actualAuthPassword,
				AuthenticationProtocol:   actualAuthProtocol,
				PrivacyPassphrase:        actualPrivPassword,
				PrivacyProtocol:          actualPrivProtocol,
			},
			MaxOids:     math.MaxInt32,
			ContextName: contextName,
		},
	}

	logger := getLogger("SNMPv3", hostname, port)
	if logger != nil {
		snmp.snmp.Logger = logger
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

	err = checkForSNMPv3Issues(oid, result)
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
	multiResult, err := s.get(oid)
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

	err = checkVariableCount(oid, result)
	if err != nil {
		return emptyMultiResult, err
	}

	err = checkForSNMPv3Issues(oid, result)
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
	multiResult, err := s.getNext(oid)
	if err != nil {
		return "{}", err
	}

	multiResultBytes, err := json.Marshal(multiResult)
	if err != nil {
		return "{}", err
	}

	return string(multiResultBytes), nil
}

func (s *session) setString(oid, value string) (multiResult, error) {
	emptyMultiResult := multiResult{}

	pdu := gosnmp.SnmpPDU{
		Name:  oid,
		Type:  gosnmp.OctetString,
		Value: value,
	}

	result, err := s.snmp.set([]gosnmp.SnmpPDU{pdu})
	if err != nil {
		return emptyMultiResult, nil
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

func (s *session) setStringJSON(oid, value string) (string, error) {
	multiResult, err := s.setString(oid, value)
	if err != nil {
		return "{}", err
	}

	multiResultBytes, err := json.Marshal(multiResult)
	if err != nil {
		return "{}", err
	}

	return string(multiResultBytes), nil
}

func (s *session) setInteger(oid string, value int) (multiResult, error) {
	emptyMultiResult := multiResult{}

	pdu := gosnmp.SnmpPDU{
		Name:  oid,
		Type:  gosnmp.Integer,
		Value: value,
	}

	result, err := s.snmp.set([]gosnmp.SnmpPDU{pdu})
	if err != nil {
		return emptyMultiResult, nil
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

func (s *session) setIntegerJSON(oid string, value int) (string, error) {
	multiResult, err := s.setInteger(oid, value)
	if err != nil {
		return "{}", err
	}

	multiResultBytes, err := json.Marshal(multiResult)
	if err != nil {
		return "{}", err
	}

	return string(multiResultBytes), nil
}

func (s *session) setIPAddress(oid, value string) (multiResult, error) {
	emptyMultiResult := multiResult{}

	pdu := gosnmp.SnmpPDU{
		Name:  oid,
		Type:  gosnmp.IPAddress,
		Value: value,
	}

	result, err := s.snmp.set([]gosnmp.SnmpPDU{pdu})
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

func (s *session) setIPAddressJSON(oid, value string) (string, error) {
	multiResult, err := s.setIPAddress(oid, value)
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
