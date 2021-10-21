package gosnmp_python_go

import (
	"fmt"
	"github.com/ftpsolutions/gosnmp"
	"net"
	"time"
)

const defaultMaxRepetitions = 20        // start out asking for 20 OIDs in a GetBulkRequest
const updateInterval = time.Second * 30 // reassess if it's been this long since an update
const updateCallThreshold = 30          // reassess if it's been this many successful calls since an update

type wrappedSNMPInterface interface {
	getSNMP() *gosnmp.GoSNMP
	getConn() net.PacketConn
	connect() error
	get(oids []string) (result *gosnmp.SnmpPacket, err error)
	getNext(oids []string) (result *gosnmp.SnmpPacket, err error)
	getBulk(oids []string, nonRepeaters uint8, maxRepetitions uint8) (result *gosnmp.SnmpPacket, err error)
	walk(oids []string) (result *gosnmp.SnmpPacket, err error)
	walkBulk(oids []string) (result *gosnmp.SnmpPacket, err error)
	set(pdus []gosnmp.SnmpPDU) (result *gosnmp.SnmpPacket, err error)
	close() error
}

type wrappedSNMP struct {
	snmp                               *gosnmp.GoSNMP
	defaultMaxRepetitions              uint8
	optimalMaxRepetitions              uint8
	lastMaxRepetitionsUpdate           time.Time
	callsSinceLastMaxRepetitionsUpdate int64
}

func (w *wrappedSNMP) getSNMP() *gosnmp.GoSNMP {
	return w.snmp
}

func (w *wrappedSNMP) getConn() net.PacketConn {
	return w.snmp.Conn
}

func (w *wrappedSNMP) connect() error {
	// starting reps value to work down from
	w.defaultMaxRepetitions = defaultMaxRepetitions

	// our current optimal reps based on which value gets the most responses
	w.optimalMaxRepetitions = w.defaultMaxRepetitions

	// force an update off the bat
	w.lastMaxRepetitionsUpdate = time.Now().Add(-updateInterval).Add(-time.Second)
	w.callsSinceLastMaxRepetitionsUpdate = updateCallThreshold + 1

	return w.snmp.Connect()
}

func (w *wrappedSNMP) get(oids []string) (result *gosnmp.SnmpPacket, err error) {
	return w.snmp.Get(formatOIDs(oids))
}

func (w *wrappedSNMP) getNext(oids []string) (result *gosnmp.SnmpPacket, err error) {
	return w.snmp.GetNext(formatOIDs(oids))
}

func (w *wrappedSNMP) getBulk(oids []string, nonRepeaters uint8, maxRepetitions uint8) (result *gosnmp.SnmpPacket, err error) {
	if w.snmp.Version == gosnmp.Version1 {
		return nil, fmt.Errorf("cannot call BULKWALK with SNMPv1")
	}

	result, err = w.snmp.GetBulk(formatOIDs(oids), nonRepeaters, maxRepetitions)

	return result, err
}

func (w *wrappedSNMP) isThereMoreToWalk(oid string, originalOID string) (bool, gosnmp.SnmpPDU) {
	// if get next returns nothing, we're good
	nextResult, _ := w.getNext([]string{formatOID(oid)})
	if nextResult == nil || len(nextResult.Variables) == 0 {
		return false, gosnmp.SnmpPDU{}
	}

	// can only be one variable from getNext
	nextVariable := nextResult.Variables[0]

	// we've reached the end, we're good
	if isThisAnEndVariable(nextVariable) {
		return false, gosnmp.SnmpPDU{}
	}

	// if get next returned the same oid, we're as good as we can be without chopping off OIDs
	if nextVariable.Name == oid {
		return false, gosnmp.SnmpPDU{}
	}

	// if get next returned something outside the parent tree, we're good
	if !hasOIDPrefix(nextVariable.Name, originalOID) {
		return false, gosnmp.SnmpPDU{}
	}

	return true, nextResult.Variables[0]
}

func (w *wrappedSNMP) specialWalk(oid string, originalOID string) (result *gosnmp.SnmpPacket, err error) {
	/*
	   There's a bit of special stuff in this function that deviates slightly from a normal walk in that
	   it'll try to find the next sibling and keeping walking from there (as long as it doesn't leave the
	   parent tree).
	*/

	result = &gosnmp.SnmpPacket{
		Variables: make([]gosnmp.SnmpPDU, 0),
	}

	ok := false
	var thisPDU gosnmp.SnmpPDU
	var lastPDU gosnmp.SnmpPDU

	for {
		// this is GetNext underneath, with some logic
		ok, thisPDU = w.isThereMoreToWalk(formatOID(oid), originalOID)

		// no more walking required
		if !ok {
			break
		}

		// got same result- so no more walking required
		if thisPDU.Name == lastPDU.Name {
			break
		}

		// record what we got
		result.Variables = append(result.Variables, thisPDU)

		// and move our OID cursor
		oid = thisPDU.Name
		lastPDU = thisPDU
	}

	deduplicateResult(result)

	return result, err
}

// TODO: slice not actually needed, but keeping the interface consistent
func (w *wrappedSNMP) walk(oids []string) (result *gosnmp.SnmpPacket, err error) {
	oids = formatOIDs(oids)

	if len(oids) != 1 {
		return nil, fmt.Errorf("oids length must be exactly 1")
	}

	// TODO: can't use WalkAll, suffers when the OID tree is a funny shape- use a variant on GetNext
	// pdus, err := w.snmp.WalkAll(oids[0])
	// if err != nil {
	//     return nil, err
	// }

	oid := oids[0]
	originalOID := oid

	return w.specialWalk(oid, originalOID)
}

// TODO: slice not actually needed, but keeping the interface consistent
func (w *wrappedSNMP) walkBulk(oids []string) (result *gosnmp.SnmpPacket, err error) {
	/*
	   There's some special stuff in here to allow us to get the benefit from bulk commands without
	   being a slave to the way the target device walks itself. In summary:

	   - start out with optimalMaxRepetitions
	   - decrement optimalMaxRepetitions if timeouts are experienced
	   - finish off with a specialWalk if all we can't decrement optimalMaxRepetitions any further
	   - record progress throughout and adjust optimalMaxRepetitions value (for future use)
	   - reassess that value every 30s or 30 successful responses
	*/

	// help us fail fast if there's nothing at this OID or if this device is offline
	_, err = w.getNext(oids)
	if err != nil {
		return result, err
	}

	oids = formatOIDs(oids)

	if w.snmp.Version == gosnmp.Version1 {
		return nil, fmt.Errorf("cannot call BULKWALK with SNMPv1")
	}

	if len(oids) != 1 {
		return nil, fmt.Errorf("oids length must be exactly 1")
	}

	// TODO: can't use BulkWalkAll, suffers when the OID tree is a funny shape- use a variant on GetBulk
	// pdus, err := w.snmp.BulkWalkAll(oids[0])
	// if err != nil {
	//     return nil, err
	// }

	result = &gosnmp.SnmpPacket{
		Variables: make([]gosnmp.SnmpPDU, 0),
	}

	var nonRepeaters uint8 = 0

	originalRetries := w.snmp.Retries
	w.snmp.Retries = 0

	originalOID := oids[0]
	oid := originalOID
	var thisResult *gosnmp.SnmpPacket
	exhaustedRetries := false

	for {
		// if we're due to reassess defaultMaxRepetitions
		if time.Now().After(w.lastMaxRepetitionsUpdate.Add(updateInterval)) || w.callsSinceLastMaxRepetitionsUpdate > updateCallThreshold {
			if w.optimalMaxRepetitions < w.defaultMaxRepetitions {
				// e.g. 1, 2, 3, 4, 5, 10, 15, ... 100
				if w.optimalMaxRepetitions < 5 {
					w.optimalMaxRepetitions += 1
				} else {
					w.optimalMaxRepetitions += 5
				}

				w.lastMaxRepetitionsUpdate = time.Now()
				w.callsSinceLastMaxRepetitionsUpdate = 0
			}
		}

		if w.snmp.Logger != nil {
			w.snmp.Logger.Printf(
				"lastMaxRepetitionsUpdate=%v, callsSinceLastMaxRepetitionsUpdate=%v, optimalMaxRepetitions=%v",
				time.Now().Sub(w.lastMaxRepetitionsUpdate),
				w.callsSinceLastMaxRepetitionsUpdate,
				w.optimalMaxRepetitions,
			)
		}

		thisResult, err = w.getBulk([]string{oid}, nonRepeaters, w.optimalMaxRepetitions)
		if err != nil {
			// if we've failed- check that decrementing would even help; if this fails the device is offline / we're at the end of the tree
			_, err = w.getNext([]string{oid})
			if err != nil {
				exhaustedRetries = true
				break
			}

			// if we get here, we've exhausted all our retries without any response so give up
			if w.optimalMaxRepetitions <= 1 {
				exhaustedRetries = true
				break
			}

			// 100, 95, 90, ... 5, 4, 3, 2, 1
			if w.optimalMaxRepetitions > 5 {
				w.optimalMaxRepetitions -= 5
			} else {
				w.optimalMaxRepetitions -= 1
			}

			w.lastMaxRepetitionsUpdate = time.Now()
			w.callsSinceLastMaxRepetitionsUpdate = 0

			continue
		}

		// likely won't happen, but for completeness
		if thisResult == nil || len(thisResult.Variables) == 0 {
			err = fmt.Errorf("nothing returned for GetBulk oid=%v, nonRepeaters=%v, optimalMaxRepetitions=%v", oid, nonRepeaters, w.optimalMaxRepetitions)

			continue
		}

		w.callsSinceLastMaxRepetitionsUpdate++

		// filter anything out of our tree (GetBulk will just keep returning what's next)
		filteredVariables := make([]gosnmp.SnmpPDU, 0)
		for _, variable := range thisResult.Variables {
			if !hasOIDPrefix(variable.Name, originalOID) {
				continue
			}

			filteredVariables = append(filteredVariables, variable)
		}

		// if we got nothing in our tree, then we're done
		if len(filteredVariables) == 0 {
			break
		}

		// record what we got
		result.Variables = append(result.Variables, filteredVariables...)

		// and move our OID cursor
		oid = result.Variables[len(result.Variables)-1].Name

		// we've reached the end, we're good
		if isThisAnEndVariable(thisResult.Variables[len(thisResult.Variables)-1]) {
			break
		}
	}

	if exhaustedRetries {
		thisResult, err = w.specialWalk(oid, originalOID)
		if err == nil && thisResult != nil && len(thisResult.Variables) > 0 {
			result.Variables = append(result.Variables, thisResult.Variables...)
		}
	}

	deduplicateResult(result)

	w.snmp.Retries = originalRetries

	return result, err

}

func (w *wrappedSNMP) set(pdus []gosnmp.SnmpPDU) (result *gosnmp.SnmpPacket, err error) {
	return w.snmp.Set(pdus)
}

func (w *wrappedSNMP) close() error {
	return w.snmp.Conn.Close()
}
