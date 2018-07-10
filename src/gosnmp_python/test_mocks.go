//+build !test

package gosnmp_python

import (
	"encoding/json"
	"net"

	"github.com/ftpsolutions/gosnmp"
)

type mockSession struct{}

func (m *mockSession) getSNMP() *gosnmp.GoSNMP {
	return nil
}

func (m *mockSession) connect() error {
	return nil
}

func (m *mockSession) get(oid string) (multiResult, error) {
	return multiResult{OID: oid}, nil
}

func (m *mockSession) getJSON(oid string) (string, error) {
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

func (m *mockSession) getNext(oid string) (multiResult, error) {
	return multiResult{OID: oid}, nil
}

func (m *mockSession) getNextJSON(oid string) (string, error) {
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

func (m *mockSession) close() error {
	return nil
}

type mockWrappedSNMP struct{}

func (m *mockWrappedSNMP) getSNMP() *gosnmp.GoSNMP {
	return &gosnmp.GoSNMP{}
}

func (m *mockWrappedSNMP) getConn() net.Conn {
	return nil
}

func (m *mockWrappedSNMP) connect() error {
	return nil
}

func (m *mockWrappedSNMP) get(oids []string) (result *gosnmp.SnmpPacket, err error) {
	return nil, nil
}

func (m *mockWrappedSNMP) getNext(oids []string) (result *gosnmp.SnmpPacket, err error) {
	return nil, nil
}

func (m *mockWrappedSNMP) close() error {
	return nil
}
