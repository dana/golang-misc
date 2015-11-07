package sysvipc

import (
	//"syscall"
	"testing"
	"fmt"
)

func TestSendRcv(t *testing.T) {
	msgSetup(t)
//	defer msgTeardown(t)

	fmt.Println("Before Send");
//	q.Send(1, []byte("test message body"), nil)
	RawSend([]byte("test message body"), q)
	fmt.Println("Before Receive");
//	msg, mtyp, err := q.Receive(64, -10, nil)
	msg,err := RawReceive(q)
	if err != nil {
		t.Error(err)
	}

//	if string(msg) != "test message body" || mtyp != 1 {
//		t.Errorf("%q %v", string(msg), mtyp)
	if string(msg) != "test message body" {
		t.Errorf("%q", string(msg))
	}
}

func RawSend(rawBytes []byte, queue MessageQueue) {
	queue.Send(1, rawBytes, nil)
}

func RawReceive(queue MessageQueue) ([]byte, error) {
	msg, _, err := queue.Receive(64, -10, nil)
//	if err != nil {
//		t.Error(err)
//	}
	
	return msg, err
}
//func TestNonBlockingSend(t *testing.T) {
//	msgSetup(t)
//	defer msgTeardown(t)
//
//	info, err := q.Stat()
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	err = q.Set(&MQInfo{
//		Perms: IpcPerms{
//			OwnerUID: info.Perms.OwnerUID,
//			OwnerGID: info.Perms.OwnerGID,
//			Mode:     info.Perms.Mode,
//		},
//		MaxBytes: 8,
//	})
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	if err := q.Send(3, []byte("more than 8"), &MQSendFlags{DontWait: true}); err != syscall.EAGAIN {
//		t.Error("too-long write should have failed", err)
//	}
//}
//
//func TestNonBlockingReceive(t *testing.T) {
//	msgSetup(t)
//	defer msgTeardown(t)
//
//	_, _, err := q.Receive(64, -99, &MQRecvFlags{DontWait: true})
//	if err != syscall.EAGAIN && err != syscall.ENOMSG {
//		t.Error("non-blocking read against empty queue should fail", err)
//	}
//}
//
var q MessageQueue

func msgSetup(t *testing.T) {
	fmt.Println("one")
//	mq, err := GetMsgQueue(0xDA7ABA5E, &MQFlags{
	mq, err := GetMsgQueue(17039435, &MQFlags{
		Create:    true,
//		Create:    false,
//		Exclusive: true,
//		Exclusive: false,
		Perms:     0600,
	})
	fmt.Println("two");
	if err != nil {
		t.Fatal(err)
	}
	q = mq
}

//func msgTeardown(t *testing.T) {
//	return
//	fmt.Println(t)
//	if err := q.Remove(); err != nil {
//		t.Fatal(err)
//	}
//}
