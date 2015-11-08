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
	var headerTest = make(map[string]string)
	headerTest["a"] = "this"
	headerTest["hello"] = "goodbye"
	retBytes,_ := createWireHeader(headerTest)
	fmt.Println(string(retBytes))
}

func createWireHeader(headerMap map[string]string) ([]byte, error) {
//	ret := []byte("19:d=localhost")
	ret := []byte("")
	for key, value := range headerMap {
//		byteKey := []byte(key)
//		byteValue := []byte(value)
//		fmt.Println(byteKey, byteValue)
		ret = append(ret, key...)
		ret = append(ret, "="...)
		ret = append(ret, value...)
		ret = append(ret, ","...)
	}
	return ret, nil
}

func parseWireHeader(testInput []byte) (map[string]string, error) {
	var retMap = make(map[string]string)
	testString := string(testInput)
	fullHeaderParts := strings.SplitN(testString, ":", 2)
	headerLength,atoiErr := strconv.Atoi(fullHeaderParts[0])
	if atoiErr != nil {
		return retMap, atoiErr
	}
	headerString := fullHeaderParts[1][0:headerLength]
	headerParts := strings.Split(headerString, ",")
	for _, part := range headerParts {
		fields := strings.Split(part, "=")
		key := fields[0]
		value := fields[1]
		retMap[key] = value
	}
	return retMap, nil
}
