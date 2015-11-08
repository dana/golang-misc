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

func parseWireHeader(test_input []byte) {
	test_string := string(test_input)
//	fmt.Println("Main", test_string)
	parts := strings.SplitN(test_string, ":", 2)
//	fmt.Println(parts)
	header_length,_ := strconv.Atoi(parts[0])
	fmt.Println(header_length)
}
