package gosnmp_python

import (
	"fmt"
	"log"
	"sync"
	"time"
)

var sessionMutex sync.RWMutex
var sessions map[uint64]sessionInterface
var lastSessionID uint64

func init() {
	sessions = make(map[uint64]sessionInterface)

	time.Sleep(time.Second) // give the Python side a little time to settle
}

func handlePanic(extra string, sessionID uint64, s sessionInterface, err error) {
	log.Printf(fmt.Sprintf("Handled \"%v\" in %v for %v - %+v\n", err, extra, sessionID, s.getSNMP()))
}

// NewRPCSessionV1 creates a new Session for SNMPv1 and returns the sessionID
func NewRPCSessionV1(hostname string, port int, community string, timeout, retries int) uint64 {
	if !GetPyPy() {
		tState := releaseGIL()
		defer reacquireGIL(tState)
	}

	session := newSessionV1(
		hostname,
		port,
		community,
		timeout,
		retries,
	)

	sessionMutex.Lock()
	sessionID := lastSessionID
	lastSessionID++
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

	session := newSessionV2c(
		hostname,
		port,
		community,
		timeout,
		retries,
	)

	sessionMutex.Lock()
	sessionID := lastSessionID
	lastSessionID++
	sessions[sessionID] = &session
	sessionMutex.Unlock()

	return sessionID
}

// NewRPCSessionV3 creates a new Session for SNMPv3 and returns the sessionID
func NewRPCSessionV3(hostname string, port int, contextName, securityUsername, privacyPassword, authPassword, securityLevel, authProtocol, privacyProtocol string, timeout, retries int) uint64 {
	if !GetPyPy() {
		tState := releaseGIL()
		defer reacquireGIL(tState)
	}

	session := newSessionV3(
		hostname,
		port,
		contextName,
		securityUsername,
		privacyPassword,
		authPassword,
		securityLevel,
		authProtocol,
		privacyProtocol,
		timeout,
		retries,
	)

	sessionMutex.Lock()
	sessionID := lastSessionID
	lastSessionID++
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
		err = val.connect()
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
		result, err = val.getJSON(oid)
	} else {
		err = fmt.Errorf("sessionID %v does not exist", sessionID)
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
		result, err = val.getNextJSON(oid)
	} else {
		err = fmt.Errorf("sessionID %v does not exist", sessionID)
	}

	return result, err
}

// RPCClose calls .close on the Session identified by the sessionID
func RPCClose(sessionID uint64) (err error) {
	if !GetPyPy() {
		tState := releaseGIL()
		defer reacquireGIL(tState)
	}

	sessionMutex.RLock()
	val, ok := sessions[sessionID]
	sessionMutex.RUnlock()

	if !ok {
		return fmt.Errorf("sessionID %v does not exist; only %v", sessionID, sessions)
	}

	sessionMutex.Lock()
	delete(sessions, sessionID)
	sessionMutex.Unlock()

	defer func(s sessionInterface) {
		if r := recover(); r != nil {
			if handledError, _ := r.(error); handledError != nil {
				handlePanic("RPCClose", sessionID, val, handledError)
			}
		}
	}(val)

	err = val.close()

	return err
}
