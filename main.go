package main

import (
	"gosnmp_python"
	"fmt"
)

func main() {

	sessionID := gosnmp_python.NewRPCSessionV2c(
		"10.10.0.2",
		161,
		"public",
		5,
		1,
	)

	err := gosnmp_python.RPCConnect(sessionID)
	if err != nil {
		panic(err)
	}

	jsonResult, err := gosnmp_python.RPCGet(sessionID, ".1.3.6.1.2.1.1.5.0")
	if err != nil {
		panic(err)
	}

	fmt.Println(jsonResult)

	err = gosnmp_python.RPCClose(sessionID)
	if err != nil {
		panic(err)
	}

}
