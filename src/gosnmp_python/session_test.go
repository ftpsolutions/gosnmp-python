package gosnmp_python

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/initialed85/gosnmp"
	"time"
)

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
		testPrivacyPassword,
		testAuthPassword,
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
			AuthenticationPassphrase: testAuthPassword,
			PrivacyProtocol:          gosnmp.AES,
			PrivacyPassphrase:        testPrivacyPassword,
		},
	)
	a.Equal(snmp.Version, gosnmp.Version3)
	a.Equal(snmp.Timeout, time.Second*5)
	a.Equal(snmp.Retries, 1)
}
