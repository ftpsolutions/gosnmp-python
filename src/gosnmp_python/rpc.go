package gosnmp_python

// #cgo pkg-config: python2
// #include <Python.h>
import "C"
import (
	"sync"
	"errors"
	"fmt"
) // this has to follow the above comments with no spaces

var mutex sync.Mutex
var sessions map[uint64]*Session
var lastSessionId uint64 = 0

// initialiser

func init() {
	sessions = make(map[uint64]*Session)
}

// constructors

func NewRPCSessionV1(hostname string, port int, community string, timeout, retries int) uint64 {
	if C.PyEval_ThreadsInitialized() != 0 {
		tState := C.PyEval_SaveThread()
		defer C.PyEval_RestoreThread(tState)
	}

	session := NewSessionV1(
		hostname,
		port,
		community,
		timeout,
		retries,
	)

	sessionId := lastSessionId

	mutex.Lock()
	sessions[sessionId] = &session
	mutex.Unlock()

	lastSessionId++

	return sessionId
}

func NewRPCSessionV2c(hostname string, port int, community string, timeout, retries int) uint64 {
	if C.PyEval_ThreadsInitialized() != 0 {
		tState := C.PyEval_SaveThread()
		defer C.PyEval_RestoreThread(tState)
	}

	session := NewSessionV2c(
		hostname,
		port,
		community,
		timeout,
		retries,
	)

	sessionId := lastSessionId

	mutex.Lock()
	sessions[sessionId] = &session
	mutex.Unlock()

	lastSessionId++

	return sessionId
}

func NewRPCSessionV3(hostname string, port int, securityUsername, privacyPassword, authPassword, securityLevel, authProtocol, privacyProtocol string, timeout, retries int) uint64 {
	if C.PyEval_ThreadsInitialized() != 0 {
		tState := C.PyEval_SaveThread()
		defer C.PyEval_RestoreThread(tState)
	}

	session := NewSessionV3(
		hostname,
		port,
		securityUsername,
		privacyPassword,
		authPassword,
		securityLevel,
		authProtocol,
		privacyProtocol,
		timeout,
		retries,
	)

	sessionId := lastSessionId

	mutex.Lock()
	sessions[sessionId] = &session
	mutex.Unlock()

	lastSessionId++

	return sessionId
}

// public functions

func RPCConnect(sessionId uint64) error {
	if C.PyEval_ThreadsInitialized() != 0 {
		tState := C.PyEval_SaveThread()
		defer C.PyEval_RestoreThread(tState)
	}

	mutex.Lock()
	defer mutex.Unlock()
	if val, ok := sessions[sessionId]; ok {
		return val.Connect()
	}

	return errors.New(fmt.Sprintf("sessionId %v does not exist", sessionId))
}

func RPCGet(sessionId uint64, oid string) (MultiResult, error) {
	if C.PyEval_ThreadsInitialized() != 0 {
		tState := C.PyEval_SaveThread()
		defer C.PyEval_RestoreThread(tState)
	}

	mutex.Lock()
	defer mutex.Unlock()
	if val, ok := sessions[sessionId]; ok {
		return val.Get(oid)
	}
	return MultiResult{}, errors.New(fmt.Sprintf("sessionId %v does not exist", sessionId))
}

func RPCGetNext(sessionId uint64, oid string) (MultiResult, error) {
	if C.PyEval_ThreadsInitialized() != 0 {
		tState := C.PyEval_SaveThread()
		defer C.PyEval_RestoreThread(tState)
	}

	mutex.Lock()
	defer mutex.Unlock()
	if val, ok := sessions[sessionId]; ok {
		return val.GetNext(oid)
	}

	return MultiResult{}, errors.New(fmt.Sprintf("sessionId %v does not exist", sessionId))
}

func RPCClose(sessionId uint64) error {
	if C.PyEval_ThreadsInitialized() != 0 {
		tState := C.PyEval_SaveThread()
		defer C.PyEval_RestoreThread(tState)
	}

	mutex.Lock()
	defer mutex.Unlock()
	if val, ok := sessions[sessionId]; ok {
		err := val.Close()
		sessions[sessionId] = nil
		delete(sessions, sessionId)
		return err
	}

	return errors.New(fmt.Sprintf("sessionId %v does not exist", sessionId))
}
