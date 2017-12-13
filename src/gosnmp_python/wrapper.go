package gosnmp_python

import (
	"github.com/initialed85/gosnmp"
	"net"
)

type wrappedSNMPInterface interface {
	getSNMP() *gosnmp.GoSNMP
	getConn() net.Conn
	Connect() error
	Get(oids []string) (result *gosnmp.SnmpPacket, err error)
	GetNext(oids []string) (result *gosnmp.SnmpPacket, err error)
}

type mockWrappedSNMP struct{}

func (m *mockWrappedSNMP) getSNMP() *gosnmp.GoSNMP {
	return &gosnmp.GoSNMP{}
}

func (m *mockWrappedSNMP) getConn() net.Conn {
	return nil
}

func (m *mockWrappedSNMP) Connect() error {
	return nil
}

func (m *mockWrappedSNMP) Get(oids []string) (result *gosnmp.SnmpPacket, err error) {
	return nil, nil
}

func (m *mockWrappedSNMP) GetNext(oids []string) (result *gosnmp.SnmpPacket, err error) {
	return nil, nil
}

type wrappedSNMP struct {
	snmp *gosnmp.GoSNMP
}

func (w *wrappedSNMP) getSNMP() *gosnmp.GoSNMP {
	return w.snmp
}

func (w *wrappedSNMP) getConn() net.Conn {
	return w.snmp.Conn
}

func (w *wrappedSNMP) Connect() error {
	return w.snmp.Connect()
}

func (w *wrappedSNMP) Get(oids []string) (result *gosnmp.SnmpPacket, err error) {
	return w.snmp.Get(oids)
}

func (w *wrappedSNMP) GetNext(oids []string) (result *gosnmp.SnmpPacket, err error) {
	return w.snmp.GetNext(oids)
}
