package sysvipc

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"testing"
)

var test_qname string = "ipc-transit-test-queue"
var transitPath string = "/tmp/ipc_transit/"

//func createWireHeader(headerMap map[string]string) ([]byte, error) {
//        headerBytes = append(headerBytes, key...)
func Send(sendMessage map[string]interface{}, qname string) error {
    var wireHeader = make(map[string]string)
	wireHeader["qname"] = qname
	sendBytes, createWireHeaderErr := createWireHeader(wireHeader)
	if createWireHeaderErr != nil {
		return createWireHeaderErr
	}
	mq, err := getQueue(qname)
	if err != nil {
		return err
	}
	jsonBytes, marshalErr := json.Marshal(sendMessage)
	if marshalErr != nil {
		return marshalErr
	}
	sendBytes = append(sendBytes, jsonBytes...)
	fmt.Println(string(sendBytes))
	sendErr := RawSend(sendBytes, mq)
	return sendErr
}

func RawSend(rawBytes []byte, queue MessageQueue) error {
	err := queue.Send(1, rawBytes, nil)
	return err
}

func TestSendRcv(t *testing.T) {
	defer func() {
		os.Remove(transitPath + test_qname)
	}()

	// How to create this message: http://play.golang.org/p/13OSJHd5xe
	// Info about seemingly fully dynamic marshal/unmarshal: http://stackoverflow.com/questions/19482612/go-golang-array-type-inside-struct-missing-type-composite-literal
	sendMessage := map[string]interface{}{
		"Name": "Wednesday",
		"Age":  6,
		"Parents": map[string]interface{}{
			"bee": "boo",
			"foo": map[string]interface{}{
				"hi": []string{"a", "b"},
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
				if v != 6.0 {  //very strange, even though it's an int in the json, it unmarshalled as a float
					t.Error(receiveErr)
				}
			}
			fmt.Println(k, "is float64", vv)
		case map[string]interface{}:
			fmt.Println(k, "is an array:")
			for i, u := range vv {
				fmt.Println(i, u)
			}
		default:
			fmt.Println(k, "is of a type I don't know how to handle", vv)
		}
	}
}

//func parseWireHeader(testInput []byte) (map[string]string, []byte, error) {
func Receive(qname string) (interface{}, error) {
	var f interface{}
	mq, err := getQueue(qname)
	if err != nil {
		return nil, err
	}
	rawBytes, receiveErr := RawReceive(mq)
	if receiveErr != nil {
		return nil, receiveErr
	}
	wireHeader, payload, parseErr := parseWireHeader(rawBytes)
	fmt.Println("recieved from qname = " + wireHeader["qname"])
	if parseErr != nil {
		return nil, parseErr
	}
	jsonErr := json.Unmarshal(payload, &f)
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
	if _, statErr := os.Stat(transitPath); os.IsNotExist(statErr) {
		//dir does not exisdt
		mkdirErr := os.Mkdir(transitPath, 0777)
		if mkdirErr != nil {
			return mkdirErr
		}
	}
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
	transitInfoFilePath := transitPath + qname
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
	if false {	//this is about always having a use of fmt.Println so I never
				//have to take it out of imports
		fmt.Println("whatevah")
	}
	return mq, err
}


func createWireHeader(headerMap map[string]string) ([]byte, error) {
    headerBytes := []byte("")
    for key, value := range headerMap {
        headerBytes = append(headerBytes, key...)
        headerBytes = append(headerBytes, "="...)
        headerBytes = append(headerBytes, value...)
        headerBytes = append(headerBytes, ","...)
    }
    if len(headerBytes) > 0 {
        headerBytes = headerBytes[:len(headerBytes)-1]
    }
    ret := []byte(strconv.Itoa(len(headerBytes)))
    ret = append(ret, ":"...)
    ret = append(ret, headerBytes...)
    return ret, nil
}

func parseWireHeader(testInput []byte) (map[string]string, []byte, error) {
    var retMap = make(map[string]string)
    testString := string(testInput)
    fullHeaderParts := strings.SplitN(testString, ":", 2)
    headerLength,atoiErr := strconv.Atoi(fullHeaderParts[0])
    if atoiErr != nil {
        return retMap, nil, atoiErr
    }
    headerString := fullHeaderParts[1][0:headerLength]
    payload := testInput[len(fullHeaderParts[0])+headerLength+1:]
    headerParts := strings.Split(headerString, ",")
    for _, part := range headerParts {
        fields := strings.Split(part, "=")
        key := fields[0]
        value := fields[1]
        retMap[key] = value
    }
    return retMap, payload, nil
}


/*
Sat Nov  7 18:09:43 PST 2015
TODO:
Locking around the IPC Transit file manipulation
Nonblocking flags
Handle (and respect) the custom IPC::Transit header
Obviously turn this into a proper package
Large message handling
Remote transit
queue stats
internal(local) queues
testing with alternate directories
*/

/*
Sun Nov  8 12:24:56 PST 2015

Full example of message including header

11:d=localhost{".ipc_transit_meta":{"destination":"localhost","ttl":9,"destination_qname":"test","send_ts":1447014248},"1":2}
*/


