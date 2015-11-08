package sysvipc

import (
	"bufio"
	"fmt"
	"os"
	"encoding/json"
	"reflect"
	"strconv"
	"strings"
	"testing"
)

var test_qname string = "foo"


func Send(sendMessage map[string]interface {}, qname string) error {
	mq, err := getQueue(qname)
	if err != nil {
		return err
	}
	jsonBytes, marshalErr := json.Marshal(sendMessage)
	if marshalErr != nil {
		return marshalErr
	}
	sendErr := RawSend(jsonBytes, mq)
	return sendErr
}

func RawSend(rawBytes []byte, queue MessageQueue) error {
	err := queue.Send(1, rawBytes, nil)
	return err
}

func TestSendRcv(t *testing.T) {
	defer func() {
		os.Remove("/tmp/ipc_transit/" + test_qname)
	}()

// How to create this message: http://play.golang.org/p/13OSJHd5xe
// Info about seemingly fully dynamic marshal/unmarshal: http://stackoverflow.com/questions/19482612/go-golang-array-type-inside-struct-missing-type-composite-literal
	sendMessage := map[string]interface{}{
		"Name": "Wednesday",
		"Age": 6,
		"Parents": map[string]interface{}{
			"bee": "boo",
			"foo": map[string]interface{}{
				"hi": []string{"a","b"},
			},
		},
	}
	sendErr := Send(sendMessage, test_qname)
	if sendErr != nil {
		t.Error(sendErr)
		return
	}
	m, receiveErr := Receive(test_qname)
	if receiveErr != nil {
		t.Error(receiveErr)
		return
	}
	msg := m.(map[string]interface{})
	for k, v := range msg {
		fmt.Println(k, " -> ", reflect.TypeOf(v))
		switch vv := v.(type) {
		case string:
			if k == "Name" {
				if v != "Wednesday" {
					t.Error(receiveErr)
				}
			}
			fmt.Println(k, "is string", vv)
		case float64:
			if k == "Age" {
				if v != 6.0 {
					t.Error(receiveErr)
				}
			}
			fmt.Println(k, "is float64", vv)
		case map[string]interface {}:
			fmt.Println(k, "is an array:")
			for i, u := range vv {
				fmt.Println(i, u)
			}
		default:
			fmt.Println(k, "is of a type I don't know how to handle", vv)
		}
	}
}

func Receive(qname string) (interface{}, error) {
	var f interface{}
	mq, err := getQueue(qname)
	if err != nil {
		return f, err
	}
	rawBytes, receiveErr := RawReceive(mq)
	jsonErr := json.Unmarshal(rawBytes, &f)
	if jsonErr != nil {
		return f, jsonErr
	}

	return f, receiveErr
}
func RawReceive(queue MessageQueue) ([]byte, error) {
	rawBytes, _, err := queue.Receive(102400000, -1, nil)
	return rawBytes, err
}

//  $ cat /tmp/ipc_transit/foo
//  qid=17039435
//  qname=foo
type transitInfo struct {
	qid   int64
	qname string
}

func parseTransitFile(filePath string) (transitInfo, error) {
	info := transitInfo{0, ""}
	fi, err := os.Open(filePath)
	if err != nil {
		return info, err
	}
	defer func() {
		if err := fi.Close(); err != nil {
			panic(err)
		}
	}()
	scanner := bufio.NewScanner(fi)
	for scanner.Scan() {
		things := strings.Split(scanner.Text(), "=")
		key := string(things[0])
		value := things[1]
		switch key {
		case "qid":
			my_qid, atoiErr := strconv.Atoi(string(value))
			if atoiErr != nil {
				return info, atoiErr
			}
			info.qid = int64(my_qid)
		case "qname":
			info.qname = string(value)
		}
	}
	if err := scanner.Err(); err != nil {
		return info, err
	}
	return info, err
}

func makeNewQueue(qname string, queuePath string) error {
	fmt.Println("makeNewQueue: " + qname)
	fi, err := os.Create(queuePath)
	if err != nil {
		return err
	}
	defer func() {
		if err := fi.Close(); err != nil {
			panic(err)
		}
	}()
	fi.WriteString("qname=" + qname + "\n")
	key, ftokErr := Ftok(queuePath, 100)
	fi.WriteString("qid=" + strconv.Itoa(int(key)))
	return ftokErr
}

func getQueue(qname string) (MessageQueue, error) {
	transitInfoFilePath := "/tmp/ipc_transit/" + qname
	if _, statErr := os.Stat(transitInfoFilePath); os.IsNotExist(statErr) {
		makeErr := makeNewQueue(qname, transitInfoFilePath)
		if makeErr != nil {
			panic(makeErr)
		}
	}
	info, err := parseTransitFile(transitInfoFilePath)
	if err != nil {
		return MessageQueue(0), err
	}
	mq, err := GetMsgQueue(info.qid, &MQFlags{
		Create: true,
		//		Create:    false,
		//		Exclusive: true,
		//		Exclusive: false,
		Perms: 0666,
	})
	if false {
		fmt.Println("two")
	}
	return mq, err
}
