package gosnmp_python_go

import (
	"fmt"
	"log"
	"runtime/debug"
	"sync"
	"time"
)

var sessionMutex sync.Mutex
var sessions map[uint64]sessionInterface
var lastSessionID uint64

func init() {
	sessions = make(map[uint64]sessionInterface)

	time.Sleep(time.Second) // give the Python side a little time to settle
}

// this is used to ensure the Go runtime keeps operating in the event of strange errors
func handlePanic(extra string, sessionID uint64, s sessionInterface, err error) {
	log.Printf(
		fmt.Sprintf(
			"handlePanic() for %v()\n\tSessionID: %v\n\tSession: %+v\n\tError: %v\n\nStack trace follows:\n\n%v",
			extra,
			sessionID,
			s,
			err,
			string(debug.Stack()),
		),
	)
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

	sessionMutex.Lock()
	val, ok := sessions[sessionID]
	sessionMutex.Unlock()

	// permit recovering from a panic but return the error
	defer func(s sessionInterface) {
		if r := recover(); r != nil {
			if handledError, _ := r.(error); handledError != nil {
				handlePanic("connect", sessionID, val, handledError)
				err = handledError
			}
		}
	}(val)

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

	sessionMutex.Lock()
	val, ok := sessions[sessionID]
	sessionMutex.Unlock()

	// permit recovering from a panic but return the error
	defer func(s sessionInterface) {
		if r := recover(); r != nil {
			if handledError, _ := r.(error); handledError != nil {
				handlePanic("getJSON", sessionID, val, handledError)
				err = handledError
			}
		}
	}(val)

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

	sessionMutex.Lock()
	val, ok := sessions[sessionID]
	sessionMutex.Unlock()

	// permit recovering from a panic but return the error
	defer func(s sessionInterface) {
		if r := recover(); r != nil {
			if handledError, _ := r.(error); handledError != nil {
				handlePanic("getNextJSON", sessionID, val, handledError)
				err = handledError
			}
		}
	}(val)

	if ok {
		result, err = val.getNextJSON(oid)
	} else {
		err = fmt.Errorf("sessionID %v does not exist", sessionID)
	}

	return result, err
}

// RPCSetString calls .setString on the Session identified by the sessionID
func RPCSetString(sessionID uint64, oid, value string) (string, error) {
	if !GetPyPy() {
		tState := releaseGIL()
		defer reacquireGIL(tState)
	}

	var err error
	var result string

	sessionMutex.Lock()
	val, ok := sessions[sessionID]
	sessionMutex.Unlock()

	// permit recovering from a panic but return the error
	defer func(s sessionInterface) {
		if r := recover(); r != nil {
			if handledError, _ := r.(error); handledError != nil {
				handlePanic("setStringJSON", sessionID, val, handledError)
				err = handledError
			}
		}
	}(val)

	if ok {
		result, err = val.setStringJSON(oid, value)
	} else {
		err = fmt.Errorf("sessionID %v does not exist", sessionID)
	}

	return result, err
}

// RPCSetInteger calls .SetInteger on the Session identified by the sessionID
func RPCSetInteger(sessionID uint64, oid string, value int) (string, error) {
	if !GetPyPy() {
		tState := releaseGIL()
		defer reacquireGIL(tState)
	}

	var err error
	var result string

	sessionMutex.Lock()
	val, ok := sessions[sessionID]
	sessionMutex.Unlock()

	// permit recovering from a panic but return the error
	defer func(s sessionInterface) {
		if r := recover(); r != nil {
			if handledError, _ := r.(error); handledError != nil {
				handlePanic("SetIntegerJSON", sessionID, val, handledError)
				err = handledError
			}
		}
	}(val)

	if ok {
		result, err = val.setIntegerJSON(oid, value)
	} else {
		err = fmt.Errorf("sessionID %v does not exist", sessionID)
	}

	return result, err
}

// RPCSetIPAddress calls .setIPAddress on the Session identified by the sessionID
func RPCSetIPAddress(sessionID uint64, oid, value string) (string, error) {
	if !GetPyPy() {
		tState := releaseGIL()
		defer reacquireGIL(tState)
	}

	var err error
	var result string

	sessionMutex.Lock()
	val, ok := sessions[sessionID]
	sessionMutex.Unlock()

	// permit recovering from a panic but return the error
	defer func(s sessionInterface) {
		if r := recover(); r != nil {
			if handledError, _ := r.(error); handledError != nil {
				handlePanic("setIPAddressJSON", sessionID, val, handledError)
				err = handledError
			}
		}
	}(val)

	if ok {
		result, err = val.setIPAddressJSON(oid, value)
	} else {
		err = fmt.Errorf("sessionID %v does not exist", sessionID)
	}

	return result, err
}

// RPCClose calls .close on the Session identified by the sessionID
func RPCClose(sessionID uint64) error {
	if !GetPyPy() {
		tState := releaseGIL()
		defer reacquireGIL(tState)
	}

	sessionMutex.Lock()
	val, ok := sessions[sessionID]
	sessionMutex.Unlock()

	if !ok {
		return nil
	}

	sessionMutex.Lock()
	delete(sessions, sessionID)
	sessionMutex.Unlock()

	// permit recovering from a panic silently (bury the error)
	defer func(s sessionInterface) {
		if r := recover(); r != nil {
			if handledError, _ := r.(error); handledError != nil {
				handlePanic("close", sessionID, val, handledError)
			}
		}
	}(val)

	return val.close()
}
