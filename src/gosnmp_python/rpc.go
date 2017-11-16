package gosnmp_python

import (
	"errors"
	"fmt"
)

import (
	_ "net/http/pprof"
	"sync"
)

// globals

var sessionMutex sync.RWMutex
var sessions map[uint64]session

var lastSessionId uint64 = 0

// initialiser

func init() {
	sessions = make(map[uint64]session)
}

// private functions

func pyPyLock() {
	if GetPyPy() {
		sessionMutex.Lock()
	}
}

func pyPyUnlock() {
	if GetPyPy() {
		sessionMutex.Unlock()
	}
}

// public functions

func RPCConnect(sessionId uint64) error {
	if !GetPyPy() {
		tState := releaseGIL()
		defer reacquireGIL(tState)
	}

	var err error

	sessionMutex.RLock()
	val, ok := sessions[sessionId]
	sessionMutex.RUnlock()

	if ok {
		pyPyLock()
		err = val.connect()
		pyPyUnlock()
	} else {
		err = errors.New(fmt.Sprintf("sessionId %v does not exist", sessionId))
	}

	return err
}

func RPCGet(sessionId uint64, oid string) (string, error) {
	if !GetPyPy() {
		tState := releaseGIL()
		defer reacquireGIL(tState)
	}

	var err error
	var result string

	sessionMutex.RLock()
	val, ok := sessions[sessionId]
	sessionMutex.RUnlock()

	if ok {
		pyPyLock()
		result, err = val.getJSON(oid)
		pyPyUnlock()
	} else {
		err = errors.New(fmt.Sprintf("sessionId %v does not exist", sessionId))
	}

	return result, err
}

func RPCGetNext(sessionId uint64, oid string) (string, error) {
	if !GetPyPy() {
		tState := releaseGIL()
		defer reacquireGIL(tState)
	}

	var err error
	var result string

	sessionMutex.RLock()
	val, ok := sessions[sessionId]
	sessionMutex.RUnlock()

	if ok {
		pyPyLock()
		result, err = val.getNextJSON(oid)
		pyPyUnlock()
	} else {
		err = errors.New(fmt.Sprintf("sessionId %v does not exist", sessionId))
	}

	return result, err
}

func RPCClose(sessionId uint64) error {
	if !GetPyPy() {
		tState := releaseGIL()
		defer reacquireGIL(tState)
	}

	var err error

	sessionMutex.RLock()
	val, ok := sessions[sessionId]
	sessionMutex.RUnlock()

	if ok {
		pyPyLock()
		err = val.close()
		pyPyUnlock()
		sessionMutex.Lock()
		delete(sessions, sessionId)
		sessionMutex.Unlock()
	} else {
		err = errors.New(fmt.Sprintf("sessionId %v does not exist; only %v", sessionId, sessions))
	}

	return err
}

// constructors

func NewRPCSessionV1(hostname string, port int, community string, timeout, retries int) uint64 {
	if !GetPyPy() {
		tState := releaseGIL()
		defer reacquireGIL(tState)
	}

	sessionMutex.Lock()
	sessionId := lastSessionId
	lastSessionId++
	sessionMutex.Unlock()

	sessionMutex.Lock()
	session := newSessionV1(
		hostname,
		port,
		community,
		timeout,
		retries,
	)
	sessionMutex.Unlock()

	sessionMutex.Lock()
	sessions[sessionId] = session
	sessionMutex.Unlock()

	return sessionId
}

func NewRPCSessionV2c(hostname string, port int, community string, timeout, retries int) uint64 {
	if !GetPyPy() {
		tState := releaseGIL()
		defer reacquireGIL(tState)
	}

	sessionMutex.Lock()
	sessionId := lastSessionId
	lastSessionId++
	sessionMutex.Unlock()

	sessionMutex.Lock()
	session := newSessionV2c(
		hostname,
		port,
		community,
		timeout,
		retries,
	)
	sessionMutex.Unlock()

	sessionMutex.Lock()
	sessions[sessionId] = session
	sessionMutex.Unlock()

	return sessionId
}

func NewRPCSessionV3(hostname string, port int, securityUsername, privacyPassword, authPassword, securityLevel, authProtocol, privacyProtocol string, timeout, retries int) uint64 {
	if !GetPyPy() {
		tState := releaseGIL()
		defer reacquireGIL(tState)
	}

	sessionMutex.Lock()
	sessionId := lastSessionId
	lastSessionId++
	sessionMutex.Unlock()

	sessionMutex.Lock()
	session := newSessionV3(
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
	sessionMutex.Unlock()

	sessionMutex.Lock()
	sessions[sessionId] = session
	sessionMutex.Unlock()

	return sessionId
}
