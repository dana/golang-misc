package sysvipc

import (
	"strconv"
	"testing"
	"strings"
	"fmt"
	"os"
	"bufio"
)
var test_qname string = "foo"

func Send(rawBytes []byte, qname string) (error) {
	mq, err := msgSetup(qname)
	if err != nil {
		return err
	}
	sendErr := RawSend(rawBytes, mq);
	return sendErr
}

func RawSend(rawBytes []byte, queue MessageQueue) (error) {
	err := queue.Send(1, rawBytes, nil)
	return err
}

func TestSendRcv(t *testing.T) {
//	mq, setupErr := msgSetup(test_qname)
//	if setupErr != nil {
//		t.Error(setupErr)
//		return
//	}

//	err := RawSend([]byte("test message body"), mq)
	err := Send([]byte("test message body"), test_qname)
	if err != nil {
		t.Error(err)
		return
	}
	
	msg,err := Receive(test_qname)
	if err != nil {
		t.Error(err)
		return
	}

	if string(msg) != "test message body" {
		t.Error(string(msg))
		return
	}
}

func Receive(qname string) ([]byte, error) {
	mq, err := msgSetup(qname)
	if err != nil {
		return []byte(""), err
	}
	msg, receiveErr := RawReceive(mq);
	return msg, receiveErr
}
func RawReceive(queue MessageQueue) ([]byte, error) {
	msg, _, err := queue.Receive(102400000, -1, nil)
	
	return msg, err
}

//  $ cat /tmp/ipc_transit/foo 
//  qid=17039435
//  qname=foo
type transitInfo struct {
	qid int64
	qname string
}

func parseTransitFile(filePath string) (transitInfo, error) {
//	info := transitInfo{17039435, "foo"}
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
	
	fmt.Println(info)
	return info, err
}

func msgSetup(qname string) (MessageQueue, error) {
	info, err := parseTransitFile("/tmp/ipc_transit/" + qname)
	if err != nil {
		return MessageQueue(0), err
	}
	mq, err := GetMsgQueue(info.qid, &MQFlags{
		Create:    true,
//		Create:    false,
//		Exclusive: true,
//		Exclusive: false,
		Perms:     0600,
	})
	if false {
		fmt.Println("two");
	}
	return mq, err
}
