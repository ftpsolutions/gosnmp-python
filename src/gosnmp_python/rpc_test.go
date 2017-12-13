package gosnmp_python

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"reflect"
	"github.com/initialed85/gosnmp"
	"time"
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

	snmp := sessions[sessionID].getSNMP()

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

	snmp := sessions[sessionID].getSNMP()

	a.Equal("some_hostname", snmp.Target)
	a.Equal(uint16(161), snmp.Port)
	a.Equal("public", snmp.Community)
	a.Equal(snmp.Version, gosnmp.Version2c)
	a.Equal(snmp.Timeout, time.Second*5)
	a.Equal(snmp.Retries, 1)

	delete(sessions, sessionID)
	lastSessionID--
}

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

	snmp := sessions[sessionID].getSNMP()

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

func TestRPCConnect(t *testing.T) {
	a := assert.New(t)

	sessionID := uint64(0)

	subject := mockSession{}

	sessions[sessionID] = &subject

	a.Equal(RPCConnect(sessionID), nil)
}

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
