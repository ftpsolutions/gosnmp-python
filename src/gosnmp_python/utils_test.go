package gosnmp_python

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestGetPyPy(t *testing.T) {
	a := assert.New(t)

	runningPyPy = false

	a.Equal(false, GetPyPy())

	runningPyPy = true

	a.Equal(true, GetPyPy())
}

func TestSetPyPy(t *testing.T) {
	a := assert.New(t)

	runningPyPy = false

	SetPyPy()

	a.Equal(true, runningPyPy)
}
