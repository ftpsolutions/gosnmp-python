package gosnmp_python

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/initialed85/gosnmp"
	"time"
	"fmt"
	"errors"
)

func TestGetSecurityLevel(t *testing.T) {
	a := assert.New(t)

	tests := []struct {
		expected gosnmp.SnmpV3MsgFlags
		input    string
	}{
		{gosnmp.NoAuthNoPriv, "noAuthNoPriv"},
		{gosnmp.AuthNoPriv, "AuthNoPriv"},
		{gosnmp.AuthPriv, "AuthPriv"},
	}

	for _, test := range tests {
		a.Equal(test.expected, getSecurityLevel(test.input))
	}
}

func TestGetPrivacyDetails(t *testing.T) {
	a := assert.New(t)

	tests := []struct {
		expectedPassword string
		expectedProtocol gosnmp.SnmpV3PrivProtocol
		inputProtocol    string
	}{
		{"", gosnmp.NoPriv, ""},
		{testPassword, gosnmp.DES, "DES"},
		{testPassword, gosnmp.AES, "AES"},
	}

	for _, test := range tests {
		password, protocol := getPrivacyDetails(testPassword, test.inputProtocol)
		a.Equal(test.expectedPassword, password)
		a.Equal(test.expectedProtocol, protocol)
	}

}

func TestGetAuthenticationDetails(t *testing.T) {
	a := assert.New(t)

	tests := []struct {
		expectedPassword string
		expectedProtocol gosnmp.SnmpV3AuthProtocol
		inputProtocol    string
	}{
		{"", gosnmp.NoAuth, "",},
		{testPassword, gosnmp.MD5, "MD5"},
		{testPassword, gosnmp.SHA, "SHA"},
	}

	for _, test := range tests {
		password, protocol := getAuthenticationDetails(testPassword, test.inputProtocol)
		a.Equal(test.expectedPassword, password)
		a.Equal(test.expectedProtocol, protocol)
	}
}

func TestBuildMultiResult(t *testing.T) {
	a := assert.New(t)

	tests := []struct {
		expectedMultiResult multiResult
		expectedError       error
		inputOID            string
		inputValueType      gosnmp.Asn1BER
		inputValue          interface{}
	}{
		{
			multiResult{OID: testOID, Type: "noSuchInstance", IsNoSuchInstance: true},
			nil,
			testOID,
			gosnmp.Null,
			nil,
		},
		{
			multiResult{OID: testOID, Type: "noSuchInstance", IsNoSuchInstance: true},
			nil,
			testOID,
			gosnmp.UnknownType,
			nil,
		},
		{
			multiResult{OID: testOID, Type: "noSuchObject", IsNoSuchObject: true},
			nil,
			testOID,
			gosnmp.NoSuchObject,
			nil,
		},
		{
			multiResult{OID: testOID, Type: "endOfMibView", IsEndOfMibView: true},
			nil,
			testOID,
			gosnmp.EndOfMibView,
			nil,
		},
		{
			multiResult{OID: testOID, Type: "bool", BoolValue: true},
			nil,
			testOID,
			gosnmp.Boolean,
			true,
		},
		{
			multiResult{OID: testOID, Type: "int", IntValue: 1337},
			nil,
			testOID,
			gosnmp.Counter32,
			uint(1337),
		},
		{
			multiResult{OID: testOID, Type: "int", IntValue: 1337},
			nil,
			testOID,
			gosnmp.Gauge32,
			uint(1337),
		},
		{
			multiResult{OID: testOID, Type: "int", IntValue: 1337},
			nil,
			testOID,
			gosnmp.Uinteger32,
			uint(1337),
		},
		{
			multiResult{OID: testOID, Type: "int", IntValue: 1337},
			nil,
			testOID,
			gosnmp.Counter64,
			uint64(1337),
		},
		{
			multiResult{OID: testOID, Type: "int", IntValue: 1337},
			nil,
			testOID,
			gosnmp.Integer,
			int(1337),
		},
		{
			multiResult{OID: testOID, Type: "int", IntValue: 1337},
			nil,
			testOID,
			gosnmp.TimeTicks,
			int(1337),
		},
		{
			multiResult{OID: testOID, Type: "float", FloatValue: 1337.1337},
			nil,
			testOID,
			gosnmp.Opaque,
			float64(1337.1337),
		},
		{
			multiResult{OID: testOID, Type: "bytearray", ByteArrayValue: []int{6, 2, 9, 1}},
			nil,
			testOID,
			gosnmp.OctetString,
			[]byte{6, 2, 9, 1},
		},
		{
			multiResult{OID: testOID, Type: "string", StringValue: "something"},
			nil,
			testOID,
			gosnmp.ObjectIdentifier,
			"something",
		},
		{
			multiResult{OID: testOID, Type: "string", StringValue: "something"},
			nil,
			testOID,
			gosnmp.IPAddress,
			"something",
		},
		{
			multiResult{OID: testOID},
			errors.New("Unknown type; oid=1.2.3.4, type=69, value=what even is this?"),
			testOID,
			gosnmp.NsapAddress,
			"what even is this?",
		},
	}

	for _, test := range tests {

		fmt.Printf("\t%v\n", test)

		result, err := buildMultiResult(test.inputOID, test.inputValueType, test.inputValue)

		a.Equal(test.expectedMultiResult, result)
		a.Equal(test.expectedError, err)
	}
}

func TestNewSessionV1(t *testing.T) {
	a := assert.New(t)

	session := newSessionV1(
		testHostname,
		testPort,
		testCommunity,
		testTimeout,
		testRetries,
	)

	snmp := session.snmp.getSNMP()

	a.Equal("some_hostname", snmp.Target)
	a.Equal(uint16(161), snmp.Port)
	a.Equal("public", snmp.Community)
	a.Equal(snmp.Version, gosnmp.Version1)
	a.Equal(snmp.Timeout, time.Second*5)
	a.Equal(snmp.Retries, 1)
}

func TestNewSessionV2c(t *testing.T) {
	a := assert.New(t)

	session := newSessionV2c(
		testHostname,
		testPort,
		testCommunity,
		testTimeout,
		testRetries,
	)

	snmp := session.snmp.getSNMP()

	a.Equal("some_hostname", snmp.Target)
	a.Equal(uint16(161), snmp.Port)
	a.Equal("public", snmp.Community)
	a.Equal(snmp.Version, gosnmp.Version2c)
	a.Equal(snmp.Timeout, time.Second*5)
	a.Equal(snmp.Retries, 1)
}

func TestNewSessionV3(t *testing.T) {
	a := assert.New(t)

	session := newSessionV3(
		testHostname,
		testPort,
		testSecurityUsername,
		testPassword,
		testPassword,
		testSecurityLevel,
		testAuthProtocol,
		testPrivacyProtocol,
		testTimeout,
		testRetries,
	)

	snmp := session.snmp.getSNMP()

	a.Equal("some_hostname", snmp.Target)
	a.Equal(uint16(161), snmp.Port)
	a.Equal(snmp.SecurityModel, gosnmp.UserSecurityModel)
	a.Equal(snmp.MsgFlags, gosnmp.AuthPriv)
	a.Equal(
		snmp.SecurityParameters,
		&gosnmp.UsmSecurityParameters{
			UserName:                 testSecurityUsername,
			AuthenticationProtocol:   gosnmp.SHA,
			AuthenticationPassphrase: testPassword,
			PrivacyProtocol:          gosnmp.AES,
			PrivacyPassphrase:        testPassword,
		},
	)
	a.Equal(snmp.Version, gosnmp.Version3)
	a.Equal(snmp.Timeout, time.Second*5)
	a.Equal(snmp.Retries, 1)
}

func TestConnect(t *testing.T) {
	a := assert.New(t)

	subject := mockSession{}

	a.Equal(nil, subject.connect())
}

func TestGet(t *testing.T) {
	a := assert.New(t)

	subject := mockSession{}

	result, err := subject.get(testOID)
	a.Equal(multiResult{OID: testOID}, result)
	a.Equal(nil, err)
}

func TestGetJSON(t *testing.T) {
	a := assert.New(t)

	subject := mockSession{}

	result, err := subject.getJSON(testOID)
	a.Equal(testJSON, result)
	a.Equal(nil, err)
}

func TestGetNext(t *testing.T) {
	a := assert.New(t)

	subject := mockSession{}

	result, err := subject.getNext(testOID)
	a.Equal(multiResult{OID: testOID}, result)
	a.Equal(nil, err)
}

func TestGetNextJSON(t *testing.T) {
	a := assert.New(t)

	subject := mockSession{}

	result, err := subject.getNextJSON(testOID)
	a.Equal(testJSON, result)
	a.Equal(nil, err)
}
