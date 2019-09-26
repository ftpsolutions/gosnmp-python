package gosnmp_python_go

import (
	"testing"

	"github.com/ftpsolutions/gosnmp"
	"github.com/stretchr/testify/assert"
)

func TestWrappedSNMPGetSNMP(t *testing.T) {
	a := assert.New(t)

	goSNMP := &gosnmp.GoSNMP{}

	subject := wrappedSNMP{
		goSNMP,
	}

	a.Equal(goSNMP, subject.getSNMP())
}

func TestWrappedSNMPGetConn(t *testing.T) {
	a := assert.New(t)

	goSNMP := &gosnmp.GoSNMP{}

	subject := wrappedSNMP{
		goSNMP,
	}

	a.Equal(nil, subject.getConn())
}

func TestWrappedSNMPConnect(t *testing.T) {
	a := assert.New(t)

	subject := mockWrappedSNMP{}

	a.Equal(nil, subject.connect())
}

func TestWrappedSNMPGet(t *testing.T) {
	a := assert.New(t)

	subject := mockWrappedSNMP{}

	result, err := subject.get([]string{"1.2.3.4"})

	a.Equal((*gosnmp.SnmpPacket)(nil), result)
	a.Equal(nil, err)
}

func TestWrappedSNMPGetNext(t *testing.T) {
	a := assert.New(t)

	subject := mockWrappedSNMP{}

	result, err := subject.getNext([]string{"1.2.3.4"})

	a.Equal((*gosnmp.SnmpPacket)(nil), result)
	a.Equal(nil, err)
}

func TestWrappedSNMPSet(t *testing.T) {
	a := assert.New(t)

	subject := mockWrappedSNMP{}

	result, err := subject.set([]gosnmp.SnmpPDU{})

	a.Equal((*gosnmp.SnmpPacket)(nil), result)
	a.Equal(nil, err)
}

func TestWrappedSNMPClose(t *testing.T) {
	a := assert.New(t)

	subject := mockWrappedSNMP{}

	a.Equal(nil, subject.close())
}
