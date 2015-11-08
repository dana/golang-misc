package main

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
)

var testInput = []byte("19:d=localhost,foo=bar{\".ipc_transit_meta\":{\"destination\":\"localhost\",\"ttl\":9,\"destination_qname\":\"test\",\"send_ts\":1447014248},\"1\":2}")

func TestParse(t *testing.T) {
	if false {
		t.Error()
	}
	fmt.Println("good to go")
}

func main() {
	wireHeader, _ := parseWireHeader(testInput)
	fmt.Println(wireHeader)
}

func parseWireHeader(testInput []byte) (map[string]string, error) {
	var ret = make(map[string]string)
	testString := string(testInput)
	headerParts := strings.SplitN(testString, ":", 2)
	headerLength,atoiErr := strconv.Atoi(headerParts[0])
	if atoiErr != nil {
		return ret, atoiErr
	}
	headerString := headerParts[1][0:headerLength]
	fmt.Println(headerString)
	return ret, nil
}
