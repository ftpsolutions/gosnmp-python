package gosnmp_python

import (
	"net"

	"github.com/ftpsolutions/gosnmp"
)

type wrappedSNMPInterface interface {
	getSNMP() *gosnmp.GoSNMP
	getConn() net.Conn
	connect() error
	get(oids []string) (result *gosnmp.SnmpPacket, err error)
	getNext(oids []string) (result *gosnmp.SnmpPacket, err error)
	close() error
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

func (w *wrappedSNMP) connect() error {
	return w.snmp.Connect()
}

func (w *wrappedSNMP) get(oids []string) (result *gosnmp.SnmpPacket, err error) {
	return w.snmp.Get(oids)
}

func (w *wrappedSNMP) getNext(oids []string) (result *gosnmp.SnmpPacket, err error) {
	return w.snmp.GetNext(oids)
}

func (w *wrappedSNMP) close() error {
	return w.snmp.Conn.Close()
}
