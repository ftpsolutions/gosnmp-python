package gosnmp_python

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"reflect"
	"encoding/json"
	"github.com/initialed85/gosnmp"
	"time"
	"github.com/stretchr/testify/mock"
)

func TestPyPyLock(t *testing.T) {
	a := assert.New(t)

	runningPyPy = true

	pyPyLock()

	a.Equal(
		int64(1),
		reflect.ValueOf(sessionMutex).FieldByName("w").FieldByName("state").Int(),
	)
}

func TestPyPyUnLock(t *testing.T) {
	a := assert.New(t)

	runningPyPy = true

	pyPyUnlock()

	a.Equal(
		int64(0),
		reflect.ValueOf(sessionMutex).FieldByName("w").FieldByName("state").Int(),
	)
}

var testHostname = "some_hostname"
var testPort = 161
var testCommunity = "public"
var testTimeout = 5
var testRetries = 1

func TestNewRPCSessionV1(t *testing.T) {
	a := assert.New(t)

	sessionID := NewRPCSessionV1(
		testHostname,
		testPort,
		testCommunity,
		testTimeout,
		testRetries,
	)

	a.Equal(uint64(0), sessionID)

	snmp := sessions[sessionID].(session).snmp

	a.Equal("some_hostname", snmp.Target)
	a.Equal(uint16(161), snmp.Port)
	a.Equal("public", snmp.Community)
	a.Equal(snmp.Version, gosnmp.Version1)
	a.Equal(snmp.Timeout, time.Second*5)
	a.Equal(snmp.Retries, 1)

	delete(sessions, sessionID)
	lastSessionID--

}

func TestNewRPCSessionV2c(t *testing.T) {
	a := assert.New(t)

	sessionID := NewRPCSessionV2c(
		testHostname,
		testPort,
		testCommunity,
		testTimeout,
		testRetries,
	)

	a.Equal(uint64(0), sessionID)

	snmp := sessions[sessionID].(session).snmp

	a.Equal("some_hostname", snmp.Target)
	a.Equal(uint16(161), snmp.Port)
	a.Equal("public", snmp.Community)
	a.Equal(snmp.Version, gosnmp.Version2c)
	a.Equal(snmp.Timeout, time.Second*5)
	a.Equal(snmp.Retries, 1)

	delete(sessions, sessionID)
	lastSessionID--
}

var testSecurityUsername = "some_security_username"
var testPrivacyPassword = "some_privacy_password"
var testAuthPassword = "some_auth_password"
var testSecurityLevel = "authPriv"
var testAuthProtocol = "SHA"
var testPrivacyProtocol = "AES"

func TestNewRPCSessionV3(t *testing.T) {
	a := assert.New(t)

	sessionID := NewRPCSessionV3(
		testHostname,
		testPort,
		testSecurityUsername,
		testPrivacyPassword,
		testAuthPassword,
		testSecurityLevel,
		testAuthProtocol,
		testPrivacyProtocol,
		testTimeout,
		testRetries,
	)

	a.Equal(uint64(0), sessionID)

	snmp := sessions[sessionID].(session).snmp

	a.Equal("some_hostname", snmp.Target)
	a.Equal(uint16(161), snmp.Port)
	a.Equal(snmp.SecurityModel, gosnmp.UserSecurityModel)
	a.Equal(snmp.MsgFlags, gosnmp.AuthPriv)
	a.Equal(
		snmp.SecurityParameters,
		&gosnmp.UsmSecurityParameters{
			UserName:                 testSecurityUsername,
			AuthenticationProtocol:   gosnmp.SHA,
			AuthenticationPassphrase: testAuthPassword,
			PrivacyProtocol:          gosnmp.AES,
			PrivacyPassphrase:        testPrivacyPassword,
		},
	)
	a.Equal(snmp.Version, gosnmp.Version3)
	a.Equal(snmp.Timeout, time.Second*5)
	a.Equal(snmp.Retries, 1)

	delete(sessions, sessionID)
	lastSessionID--
}

type mockSession struct {
	mock.Mock
}

func (m mockSession) connect() error {
	return nil
}

func (m mockSession) get(oid string) (multiResult, error) {
	return multiResult{OID: oid}, nil
}

func (m mockSession) getJSON(oid string) (string, error) {
	snmpResult, err := m.get(oid)
	if err != nil {
		return "", err
	}

	jsonResult, err := json.Marshal(snmpResult)
	if err != nil {
		return "", err
	}

	return string(jsonResult), err
}

func (m mockSession) getNext(oid string) (multiResult, error) {
	return multiResult{OID: oid}, nil
}

func (m mockSession) getNextJSON(oid string) (string, error) {
	snmpResult, err := m.getNext(oid)
	if err != nil {
		return "", err
	}

	jsonResult, err := json.Marshal(snmpResult)
	if err != nil {
		return "", err
	}

	return string(jsonResult), err
}

func (m mockSession) close() error {
	return nil
}

func TestRPCConnect(t *testing.T) {
	a := assert.New(t)

	sessionID := uint64(0)

	subject := mockSession{}

	sessions[sessionID] = &subject

	a.Equal(RPCConnect(sessionID), nil)
}

var testJSON = "{\"OID\":\"1.2.3.4\",\"Type\":\"\",\"IsNull\":false,\"IsUnknown\":false,\"IsNoSuchInstance\":false,\"IsNoSuchObject\":false,\"IsEndOfMibView\":false,\"BoolValue\":false,\"IntValue\":0,\"FloatValue\":0,\"ByteArray\":null,\"StringValue\":\"\"}"

func TestRPGet(t *testing.T) {
	a := assert.New(t)

	sessionID := uint64(0)

	subject := mockSession{}

	sessions[sessionID] = &subject

	result, err := RPCGet(sessionID, "1.2.3.4")

	a.Equal(result, testJSON)
	a.Equal(err, nil)
}

func TestRPGetNext(t *testing.T) {
	a := assert.New(t)

	sessionID := uint64(0)

	subject := mockSession{}

	sessions[sessionID] = &subject

	result, err := RPCGetNext(sessionID, "1.2.3.4")

	a.Equal(result, testJSON)
	a.Equal(err, nil)
}

func TestRPCClose(t *testing.T) {
	a := assert.New(t)

	sessionID := uint64(0)

	subject := mockSession{}

	sessions[sessionID] = &subject

	a.Equal(RPCClose(sessionID), nil)

	a.Equal(sessions[sessionID], nil)
}
