package sysvipc

import (
	"testing"
	"fmt"
	"os"
)

func TestSendRcv(t *testing.T) {
	mq, setupErr := msgSetup("foo")
	if setupErr != nil {
		t.Error(setupErr)
		return
	}

	err := RawSend([]byte("test message body"), mq)
	if err != nil {
		t.Error(err)
		return
	}
	
	msg,err := RawReceive(mq)
	if err != nil {
		t.Error(err)
		return
	}

	if string(msg) != "test message body" {
		t.Errorf("%q", string(msg))
		return
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
type transitInfo struct {
	qid int64
	qname string
}

func parseTransitFile(filePath string) (transitInfo, error) {
	info := transitInfo{17039435, "foo"}
	fi, err := os.Open(filePath)
	if err != nil {
		return info, err
	}
	defer func() {
		if err := fi.Close(); err != nil {
			panic(err)
		}
	}()
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
