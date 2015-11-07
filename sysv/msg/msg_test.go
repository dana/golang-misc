package sysvipc

import (
	"testing"
	"fmt"
)

func TestSendRcv(t *testing.T) {
	mq, _ := msgSetup("foo")

	err := RawSend([]byte("test message body"), mq)
	if err != nil {
		t.Error(err)
	}
	
	msg,err := RawReceive(mq)
	if err != nil {
		t.Error(err)
	}

	if string(msg) != "test message body" {
		t.Errorf("%q", string(msg))
	}
}

func RawSend(rawBytes []byte, queue MessageQueue) (error) {
	err := queue.Send(1, rawBytes, nil)
	return err
}

func RawReceive(queue MessageQueue) ([]byte, error) {
	msg, _, err := queue.Receive(102400000, -1, nil)
	
	return msg, err
}

//  $ cat /tmp/ipc_transit/foo 
//  qid=17039435
//  qname=foo

func msgSetup(qname string) (MessageQueue, error) {
	mq, err := GetMsgQueue(17039435, &MQFlags{
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
