package gosnmp_python

// #cgo pkg-config: python2
// #include <Python.h>
import "C"
import "sync"

var runningPyPy = false
var runningPyPyMutex sync.RWMutex

func releaseGIL() *C.PyThreadState {
	var tState *C.PyThreadState

	if C.PyEval_ThreadsInitialized() == 0 {
		return nil
	}

	tState = C.PyEval_SaveThread()

	return tState
}

func reacquireGIL(tState *C.PyThreadState) {
	if tState == nil {
		return
	}

	C.PyEval_RestoreThread(tState)
}

// SetPyPy is used by the Python side to declare whether or not we're running under PyPy (can't be discovered on the Go side)
func SetPyPy() {
	runningPyPyMutex.Lock()
	runningPyPy = true
	runningPyPyMutex.Unlock()
}

// GetPyPy returns true if we're running under PyPy
func GetPyPy() bool {
	runningPyPyMutex.RLock()
	val := runningPyPy
	runningPyPyMutex.RUnlock()

	return val
}
