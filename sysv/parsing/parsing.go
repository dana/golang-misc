package main

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
)

var test_string string = `11:d=localhost{".ipc_transit_meta":{"destination":"localhost","ttl":9,"destination_qname":"test","send_ts":1447014248},"1":2}`

func TestParse(t *testing.T) {
	if false {
		t.Error()
	}
	fmt.Println("good to go")
}

func main() {
//	fmt.Println("Main", test_string)
	z := strings.SplitN(test_string, ":", 2)[0]
	fmt.Println(z)
}
