package sysvipc

import (
	"bufio"
	"fmt"
	"os"
	"encoding/json"
	"strconv"
	"strings"
	"testing"
)

var test_qname string = "foo"

func Send(rawBytes []byte, qname string) error {
	mq, err := getQueue(qname)
	if err != nil {
		return err
	}
	sendErr := RawSend(rawBytes, mq)
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
	err := Send([]byte(`{"Name":"Wednesday","Age":6,"Parents":["Gomez","Morticia"]}`), test_qname)
	if err != nil {
		t.Error(err)
		return
	}
	m, err := Receive(test_qname)
	msg := m.(map[string]interface{})
	for k, v := range msg {
		switch vv := v.(type) {
		case string:
			fmt.Println(k, "is string", vv)
		case float64:
			fmt.Println(k, "is float64", vv)
		case []interface{}:
			fmt.Println(k, "is an array:")
			for i, u := range vv {
				fmt.Println(i, u)
			}
		default:
			fmt.Println(k, "is of a type I don't know how to handle")
		}
	}
	if err != nil {
		t.Error(err)
		return
	}

//	if string(msg) != `{"Name":"Wednesday","Age":6,"Parents":["Gomez","Morticia"]}` {
//		t.Error(string(msg))
//		return
//	}
}

func Receive(qname string) (interface{}, error) {
	mq, err := getQueue(qname)
	if err != nil {
		return []byte(""), err
	}
	rawBytes, receiveErr := RawReceive(mq)
	var f interface{}
	jsonErr := json.Unmarshal(rawBytes, &f)
	if jsonErr != nil {
		return rawBytes, jsonErr
	}

	return f, receiveErr
}
func RawReceive(queue MessageQueue) ([]byte, error) {
	msg, _, err := queue.Receive(102400000, -1, nil)
	return msg, err
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
			my_qid, _ := strconv.Atoi(string(value))
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
