package sysvipc

import (
	//"syscall"
	"testing"
	"fmt"
)

func TestSendRcv(t *testing.T) {
	mq, _ := msgSetup(t)

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

func msgSetup(t *testing.T) (MessageQueue, error) {
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

//func msgTeardown(t *testing.T) {
//	return
//	fmt.Println(t)
//	if err := q.Remove(); err != nil {
//		t.Fatal(err)
//	}
//}
