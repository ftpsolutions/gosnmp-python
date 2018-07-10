package gosnmp_python

import (
	"testing"
	"time"

	"github.com/ftpsolutions/gosnmp"
	"github.com/stretchr/testify/assert"
)

func init() {
	setTesting()
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

	a.Equal(testHostname, snmp.Target)
	a.Equal(uint16(161), snmp.Port)
	a.Equal(testCommunity, snmp.Community)
	a.Equal(gosnmp.Version1, snmp.Version)
	a.Equal(time.Second*5, snmp.Timeout)
	a.Equal(1, snmp.Retries)

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

	a.Equal(testHostname, snmp.Target)
	a.Equal(uint16(161), snmp.Port)
	a.Equal(testCommunity, snmp.Community)
	a.Equal(gosnmp.Version2c, snmp.Version)
	a.Equal(time.Second*5, snmp.Timeout)
	a.Equal(1, snmp.Retries)

	delete(sessions, sessionID)
	lastSessionID--
}

func TestNewRPCSessionV3(t *testing.T) {
	a := assert.New(t)

	sessionID := NewRPCSessionV3(
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

	a.Equal(uint64(0), sessionID)

	snmp := sessions[sessionID].getSNMP()

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

	result, err := RPCGet(sessionID, testOID)

	a.Equal(result, testJSON)
	a.Equal(err, nil)
}

func TestRPGetNext(t *testing.T) {
	a := assert.New(t)

	sessionID := uint64(0)

	subject := mockSession{}

	sessions[sessionID] = &subject

	result, err := RPCGetNext(sessionID, testOID)

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
