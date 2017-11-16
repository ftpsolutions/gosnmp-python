package gosnmp_python

// #cgo pkg-config: python2
// #include <Python.h>

import (
	"errors"
	"fmt"
)

import (
	_ "net/http/pprof"
	"sync"
)

var mutex sync.RWMutex
var sessions map[uint64]*Session
var lastSessionId uint64 = 0

// initialiser

func init() {
	sessions = make(map[uint64]*Session)
}

// constructors

func NewRPCSessionV1(hostname string, port int, community string, timeout, retries int) uint64 {
	tState := releaseGIL()
	defer reacquireGIL(tState)

	mutex.Lock()
	session := NewSessionV1(
		hostname,
		port,
		community,
		timeout,
		retries,
	)
	mutex.Unlock()

	mutex.Lock()
	sessionId := lastSessionId
	lastSessionId++
	mutex.Unlock()

	mutex.RLock()
	sessions[sessionId] = &session
	mutex.RUnlock()

	return sessionId
}

func NewRPCSessionV2c(hostname string, port int, community string, timeout, retries int) uint64 {
	tState := releaseGIL()
	defer reacquireGIL(tState)

	mutex.Lock()
	session := NewSessionV2c(
		hostname,
		port,
		community,
		timeout,
		retries,
	)
	mutex.Unlock()

	mutex.Lock()
	sessionId := lastSessionId
	lastSessionId++
	mutex.Unlock()

	mutex.Lock()
	sessions[sessionId] = &session
	mutex.Unlock()

	return sessionId
}

func NewRPCSessionV3(hostname string, port int, securityUsername, privacyPassword, authPassword, securityLevel, authProtocol, privacyProtocol string, timeout, retries int) uint64 {
	tState := releaseGIL()
	defer reacquireGIL(tState)

	mutex.Lock()
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
	mutex.Unlock()

	mutex.Lock()
	sessionId := lastSessionId
	lastSessionId++
	mutex.Unlock()

	mutex.RLock()
	sessions[sessionId] = &session
	mutex.RUnlock()

	return sessionId
}

// public functions

func RPCConnect(sessionId uint64) error {
	tState := releaseGIL()
	defer reacquireGIL(tState)

	mutex.RLock()
	val, ok := sessions[sessionId]
	mutex.RUnlock()

	if ok {
		return val.Connect()
	}

	return errors.New(fmt.Sprintf("sessionId %v does not exist", sessionId))
}

func RPCGet(sessionId uint64, oid string) (MultiResult, error) {
	tState := releaseGIL()
	defer reacquireGIL(tState)

	mutex.RLock()
	val, ok := sessions[sessionId]
	mutex.RUnlock()

	if ok {
		return val.Get(oid)
	}

	return MultiResult{}, errors.New(fmt.Sprintf("sessionId %v does not exist", sessionId))
}

func RPCGetJSON(sessionId uint64, oid string) (string, error) {
	tState := releaseGIL()
	defer reacquireGIL(tState)

	mutex.RLock()
	val, ok := sessions[sessionId]
	mutex.RUnlock()

	if ok {
		return val.GetJSON(oid)
	}

	return "{}", errors.New(fmt.Sprintf("sessionId %v does not exist", sessionId))
}

func RPCGetNext(sessionId uint64, oid string) (MultiResult, error) {
	tState := releaseGIL()
	defer reacquireGIL(tState)

	mutex.RLock()
	val, ok := sessions[sessionId]
	mutex.RUnlock()

	if ok {
		return val.GetNext(oid)
	}

	return MultiResult{}, errors.New(fmt.Sprintf("sessionId %v does not exist", sessionId))
}

func RPCGetNextJSON(sessionId uint64, oid string) (string, error) {
	tState := releaseGIL()
	defer reacquireGIL(tState)

	mutex.RLock()
	val, ok := sessions[sessionId]
	mutex.RUnlock()

	if ok {
		return val.GetNextJSON(oid)
	}

	return "{}", errors.New(fmt.Sprintf("sessionId %v does not exist", sessionId))
}

func RPCClose(sessionId uint64) error {
	tState := releaseGIL()
	defer reacquireGIL(tState)

	mutex.RLock()
	val, ok := sessions[sessionId]
	mutex.RUnlock()

	if ok {
		err := val.Close()

		mutex.Lock()
		delete(sessions, sessionId)
		mutex.Unlock()

		return err
	}

	return errors.New(fmt.Sprintf("sessionId %v does not exist; only %v", sessionId, sessions))
}
