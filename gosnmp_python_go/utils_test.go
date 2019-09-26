package gosnmp_python_go

import (
	"github.com/stretchr/testify/assert"
	"testing"
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
