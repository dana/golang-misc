package main

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
)

var test_input = []byte("11:d=localhost{\".ipc_transit_meta\":{\"destination\":\"localhost\",\"ttl\":9,\"destination_qname\":\"test\",\"send_ts\":1447014248},\"1\":2}")

func TestParse(t *testing.T) {
	if false {
		t.Error()
	}
	fmt.Println("good to go")
}

func main() {
	parseWireHeader(test_input)
}

func parseWireHeader(test_input []byte) (map[string]string, error) {
	var ret = make(map[string]string)
	test_string := string(test_input)
	parts := strings.SplitN(test_string, ":", 2)
	header_length,atoiErr := strconv.Atoi(parts[0])
	if atoiErr != nil {
		return ret, atoiErr
	}
	header_string := parts[1][0:header_length]
	fmt.Println(header_string)
	return ret, nil
}
