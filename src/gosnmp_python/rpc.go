package gosnmp_python

import (
	"errors"
	"fmt"
)

import (
	"sync"
)

var sessionMutex sync.RWMutex
var sessions map[uint64]sessionInterface
var lastSessionID uint64

func init() {
	sessions = make(map[uint64]sessionInterface)
}

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

// NewRPCSessionV1 creates a new Session for SNMPv1 and returns the sessionID
func NewRPCSessionV1(hostname string, port int, community string, timeout, retries int) uint64 {
	if !GetPyPy() {
		tState := releaseGIL()
		defer reacquireGIL(tState)
	}

	sessionMutex.Lock()
	sessionID := lastSessionID
	lastSessionID++
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
	sessions[sessionID] = &session
	sessionMutex.Unlock()

	return sessionID
}

// NewRPCSessionV2c creates a new Session for SNMPv2c and returns the sessionID
func NewRPCSessionV2c(hostname string, port int, community string, timeout, retries int) uint64 {
	if !GetPyPy() {
		tState := releaseGIL()
		defer reacquireGIL(tState)
	}

	sessionMutex.Lock()
	sessionID := lastSessionID
	lastSessionID++
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
	sessions[sessionID] = &session
	sessionMutex.Unlock()

	return sessionID
}

// NewRPCSessionV3 creates a new Session for SNMPv3 and returns the sessionID
func NewRPCSessionV3(hostname string, port int, securityUsername, privacyPassword, authPassword, securityLevel, authProtocol, privacyProtocol string, timeout, retries int) uint64 {
	if !GetPyPy() {
		tState := releaseGIL()
		defer reacquireGIL(tState)
	}

	sessionMutex.Lock()
	sessionID := lastSessionID
	lastSessionID++
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
	sessions[sessionID] = &session
	sessionMutex.Unlock()

	return sessionID
}

// RPCConnect calls .connect on the Session identified by the sessionID
func RPCConnect(sessionID uint64) error {
	if !GetPyPy() {
		tState := releaseGIL()
		defer reacquireGIL(tState)
	}

	var err error

	sessionMutex.RLock()
	val, ok := sessions[sessionID]
	sessionMutex.RUnlock()

	if ok {
		pyPyLock()
		err = val.connect()
		pyPyUnlock()
	} else {
		err = fmt.Errorf("sessionID %v does not exist", sessionID)
	}

	return err
}

// RPCGet calls .get on the Session identified by the sessionID
func RPCGet(sessionID uint64, oid string) (string, error) {
	if !GetPyPy() {
		tState := releaseGIL()
		defer reacquireGIL(tState)
	}

	var err error
	var result string

	sessionMutex.RLock()
	val, ok := sessions[sessionID]
	sessionMutex.RUnlock()

	if ok {
		pyPyLock()
		result, err = val.getJSON(oid)
		pyPyUnlock()
	} else {
		err = errors.New(fmt.Sprintf("sessionID %v does not exist", sessionID))
	}

	return result, err
}

// RPCGetNext calls .getNext on the Session identified by the sessionID
func RPCGetNext(sessionID uint64, oid string) (string, error) {
	if !GetPyPy() {
		tState := releaseGIL()
		defer reacquireGIL(tState)
	}

	var err error
	var result string

	sessionMutex.RLock()
	val, ok := sessions[sessionID]
	sessionMutex.RUnlock()

	if ok {
		pyPyLock()
		result, err = val.getNextJSON(oid)
		pyPyUnlock()
	} else {
		err = errors.New(fmt.Sprintf("sessionID %v does not exist", sessionID))
	}

	return result, err
}

// RPCClose calls .close on the Session identified by the sessionID
func RPCClose(sessionID uint64) error {
	if !GetPyPy() {
		tState := releaseGIL()
		defer reacquireGIL(tState)
	}

	var err error

	sessionMutex.RLock()
	val, ok := sessions[sessionID]
	sessionMutex.RUnlock()

	if ok {
		pyPyLock()
		err = val.close()
		pyPyUnlock()
		sessionMutex.Lock()
		delete(sessions, sessionID)
		sessionMutex.Unlock()
	} else {
		err = errors.New(fmt.Sprintf("sessionID %v does not exist; only %v", sessionID, sessions))
	}

	return err
}
