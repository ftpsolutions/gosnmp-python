package gosnmp_python_go

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/ftpsolutions/gosnmp"
	"github.com/stretchr/testify/assert"
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
		{"", gosnmp.NoAuth, ""},
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
			uint(1337),
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
			errors.New("Unknown type; oid=1.2.3.4, type=0, value=what even is this?"),
			testOID,
			gosnmp.UnknownType,
			"what even is this?",
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

	a.Equal(testHostname, snmp.Target)
	a.Equal(uint16(161), snmp.Port)
	a.Equal(testCommunity, snmp.Community)
	a.Equal(gosnmp.Version1, snmp.Version)
	a.Equal(time.Second*5, snmp.Timeout)
	a.Equal(1, snmp.Retries)
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

	a.Equal(testHostname, snmp.Target)
	a.Equal(uint16(161), snmp.Port)
	a.Equal(testCommunity, snmp.Community)
	a.Equal(gosnmp.Version2c, snmp.Version)
	a.Equal(time.Second*5, snmp.Timeout)
	a.Equal(1, snmp.Retries)
}

func TestNewSessionV3(t *testing.T) {
	a := assert.New(t)

	session := newSessionV3(
		testHostname,
		testPort,
		testContextName,
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

	a.Equal(testHostname, snmp.Target)
	a.Equal(uint16(161), snmp.Port)
	a.Equal(testContextName, snmp.ContextName)
	a.Equal(gosnmp.UserSecurityModel, snmp.SecurityModel)
	a.Equal(gosnmp.AuthPriv, snmp.MsgFlags)
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
	a.Equal(gosnmp.Version3, snmp.Version)
	a.Equal(time.Second*5, snmp.Timeout)
	a.Equal(1, snmp.Retries)
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

func TestSetString(t *testing.T) {
	a := assert.New(t)

	subject := mockSession{}

	result, err := subject.setString(testOID, "")
	a.Equal(multiResult{OID: testOID}, result)
	a.Equal(nil, err)
}

func TestSetStringJSON(t *testing.T) {
	a := assert.New(t)

	subject := mockSession{}

	result, err := subject.setStringJSON(testOID, "")
	a.Equal(testJSON, result)
	a.Equal(nil, err)
}

func TestSetInteger(t *testing.T) {
	a := assert.New(t)

	subject := mockSession{}

	result, err := subject.setInteger(testOID, 0)
	a.Equal(multiResult{OID: testOID}, result)
	a.Equal(nil, err)
}

func TestSetIntegerJSON(t *testing.T) {
	a := assert.New(t)

	subject := mockSession{}

	result, err := subject.setIntegerJSON(testOID, 0)
	a.Equal(testJSON, result)
	a.Equal(nil, err)
}

func TestSetIPAddress(t *testing.T) {
	a := assert.New(t)

	subject := mockSession{}

	result, err := subject.setIPAddress(testOID, "")
	a.Equal(multiResult{OID: testOID}, result)
	a.Equal(nil, err)
}

func TestSetIPAddressJSON(t *testing.T) {
	a := assert.New(t)

	subject := mockSession{}

	result, err := subject.setIPAddressJSON(testOID, "")
	a.Equal(testJSON, result)
	a.Equal(nil, err)
}

func TestClose(t *testing.T) {
	a := assert.New(t)

	subject := mockSession{}

	a.Equal(nil, subject.close())
}
