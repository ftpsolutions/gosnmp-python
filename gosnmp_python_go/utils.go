package gosnmp_python_go

// #cgo pkg-config: python2
// #include <Python.h>
import "C"
import (
	"fmt"
	"github.com/ftpsolutions/gosnmp"
	"math"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// releaseGIL releases (unlocks) the Python GIL
func releaseGIL() *C.PyThreadState {
	var tState *C.PyThreadState

	tState = C.PyEval_SaveThread()

	return tState
}

// reacquireGIL reacquires (locks) the Python GIL
func reacquireGIL(tState *C.PyThreadState) {
	if tState == nil {
		return
	}

	C.PyEval_RestoreThread(tState)
}

// formatOID ensures an OID string is in the right format
func formatOID(oid string) string {
	return fmt.Sprintf(".%v", strings.Trim(oid, "."))
}

// formatOIDs ensures each OID in a slice is in the right format
func formatOIDs(oids []string) []string {
	formattedOIDs := make([]string, 0)

	for _, oid := range oids {
		formattedOIDs = append(formattedOIDs, formatOID(oid))
	}

	return formattedOIDs
}

// splitOID splits an OID string into a slice
func splitOID(oid string) []string {
	return strings.Split(strings.Trim(oid, "."), ".")
}

// oidToFloat converts an OID to a float for comparison
func oidToFloat(oid string) (float64, error) {
	oidParts := splitOID(oid)

	oidFloat := 0.0
	for i := 0; i < len(oidParts); i++ {
		oidPartInt64, err := strconv.ParseInt(oidParts[i], 10, 64)
		if err != nil {
			return 0, err
		}

		oidFloat += float64(oidPartInt64) * math.Pow10(128-i)
	}

	return oidFloat, nil
}

// getVariableByName returns variables by their OID for comparison, sorting etc
func getVariableByName(result *gosnmp.SnmpPacket) map[string]gosnmp.SnmpPDU {
	variablesByName := make(map[string]gosnmp.SnmpPDU, 0)

	if result == nil {
		return variablesByName
	}

	if len(result.Variables) == 0 {
		return variablesByName
	}

	for _, variable := range result.Variables {
		_, ok := variablesByName[variable.Name]
		if ok {
			continue
		}

		variablesByName[variable.Name] = variable
	}

	return variablesByName
}

// sortVariables sorts the result.Variables in-place
func sortVariables(result *gosnmp.SnmpPacket) {
	if result == nil {
		return
	}

	if len(result.Variables) == 0 {
		return
	}

	sort.SliceStable(
		result.Variables,
		func(i, j int) bool {
			var err error
			var oidI, oidJ float64

			// TODO: probably safe- .Name will never not be an OID
			oidI, err = oidToFloat(result.Variables[i].Name)
			if err != nil {
				panic(err)
			}

			// TODO: probably safe- .Name will never not be an OID
			oidJ, err = oidToFloat(result.Variables[j].Name)
			if err != nil {
				panic(err)
			}

			return oidI < oidJ
		},
	)
}

// deduplicateResult deduplicates result.Variables in-place
func deduplicateResult(result *gosnmp.SnmpPacket) {
	if result == nil {
		return
	}

	if len(result.Variables) == 0 {
		return
	}

	deduplicatedVariables := make([]gosnmp.SnmpPDU, 0)
	for _, variable := range getVariableByName(result) {
		deduplicatedVariables = append(deduplicatedVariables, variable)
	}

	result.Variables = deduplicatedVariables

	sortVariables(result)
}

// isThisAnEndVariable returns true if this is EndOfMibView, etc
func isThisAnEndVariable(pdu gosnmp.SnmpPDU) bool {
	switch pdu.Type {
	case gosnmp.NoSuchInstance:
		fallthrough
	case gosnmp.NoSuchObject:
		fallthrough
	case gosnmp.EndOfMibView:
		fallthrough
	case gosnmp.EndOfContents:
		return true
	}

	return false
}

// hasOIDPrefix return true if the given oid matches the given oidPrefix (e.g. oid=.1.3.6.1.3.69 will match oidPrefix=.1.3.6.1.3)
func hasOIDPrefix(oid string, oidPrefix string) bool {
	patternEscapedOidPrefix := strings.ReplaceAll(oidPrefix, ".", "\\.")

	pattern := fmt.Sprintf("^%v($|\\.)", patternEscapedOidPrefix)

	compiledPattern := regexp.MustCompile(pattern)

	match := compiledPattern.MatchString(oid)

	return match
}
