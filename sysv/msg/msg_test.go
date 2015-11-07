package sysvipc

import (
	"syscall"
	"testing"
)

func TestSendRcv(t *testing.T) {
	msgSetup(t)
	defer msgTeardown(t)

	if err := q.Send(-1, nil, nil); err != syscall.EINVAL {
		t.Error("msgsnd with negative mtyp should fail", err)
	}

	if _, _, err := MessageQueue(5).Receive(64, -1, nil); err != syscall.EINVAL {
		t.Error("msgrcv with bad msqid should fail", err)
	}

	q.Send(6, []byte("test message body"), nil)
	msg, mtyp, err := q.Receive(64, -100, nil)
	if err != nil {
		t.Error(err)
	}

	if string(msg) != "test message body" || mtyp != 6 {
		t.Errorf("%q %v", string(msg), mtyp)
	}
}

func TestNonBlockingSend(t *testing.T) {
	msgSetup(t)
	defer msgTeardown(t)

	info, err := q.Stat()
	if err != nil {
		t.Fatal(err)
	}

	err = q.Set(&MQInfo{
		Perms: IpcPerms{
			OwnerUID: info.Perms.OwnerUID,
			OwnerGID: info.Perms.OwnerGID,
			Mode:     info.Perms.Mode,
		},
		MaxBytes: 8,
	})
	if err != nil {
		t.Fatal(err)
	}

	if err := q.Send(3, []byte("more than 8"), &MQSendFlags{DontWait: true}); err != syscall.EAGAIN {
		t.Error("too-long write should have failed", err)
	}
}

func TestNonBlockingReceive(t *testing.T) {
	msgSetup(t)
	defer msgTeardown(t)

	_, _, err := q.Receive(64, -99, &MQRecvFlags{DontWait: true})
	if err != syscall.EAGAIN && err != syscall.ENOMSG {
		t.Error("non-blocking read against empty queue should fail", err)
	}
}

func TestMSGNOERR(t *testing.T) {
	msgSetup(t)
	defer msgTeardown(t)

	if err := q.Send(3, []byte("this message is pretty long"), nil); err != nil {
		t.Fatal(err)
	}

	msg, mtyp, err := q.Receive(22, -99, &MQRecvFlags{Truncate: true})
	if err != nil {
		t.Fatal(err)
	}
	if mtyp != 3 {
		t.Error("wrong type?")
	}
	if string(msg) != "this message is pretty" {
		t.Errorf("not properly truncated. '%s'", string(msg))
	}
}

func TestMSGSet(t *testing.T) {
	msgSetup(t)
	defer msgTeardown(t)

	info, err := q.Stat()
	if err != nil {
		t.Fatal(err)
	}

	mqi := &MQInfo{
		Perms: IpcPerms{
			OwnerUID: info.Perms.OwnerUID,
			OwnerGID: info.Perms.OwnerGID,
			Mode:     0644,
		},
	}

	err = q.Set(mqi)
	if err != nil {
		t.Fatal(err)
	}

	info, err = q.Stat()
	if err != nil {
		t.Fatal(err)
	}

	if info.Perms.Mode != 0644 {
		t.Error("perms change didn't take")
	}
}

var q MessageQueue

func msgSetup(t *testing.T) {
	mq, err := GetMsgQueue(0xDA7ABA5E, &MQFlags{
		Create:    true,
		Exclusive: true,
		Perms:     0600,
	})
	if err != nil {
		t.Fatal(err)
	}
	q = mq
}

func msgTeardown(t *testing.T) {
	if err := q.Remove(); err != nil {
		t.Fatal(err)
	}
}
