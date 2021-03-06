package service

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
)

/*
func (t *ServiceSetup) SetInfo(name, num string) (string, error) {

	eventID := "eventSetInfo"
	reg, notifier := regitserEvent(t.Client, t.ChaincodeID, eventID)
	defer t.Client.UnregisterChaincodeEvent(reg)

	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "set", Args: [][]byte{[]byte(name), []byte(num), []byte(eventID)}}
	respone, err := t.Client.Execute(req)
	if err != nil {
		return "", err
	}

	err = eventResult(notifier, eventID)
	if err != nil {
		return "", err
	}

	return string(respone.TransactionID), nil
}


func (t *ServiceSetup) GetInfo(name string) (String, error){

	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "get", Args: [][]byte{[]byte(name)}}
	respone, err := t.Client.Query(req)
	if err != nil {
		return "", err
	}

	return string(respone.Payload), nil
}
*/

func (t *ServiceSetup) Save(m Mechanic) (string, error) {

	eventID := "eventAdd"
	reg, notifier := regitserEvent(t.Client, t.ChaincodeID, eventID)
	defer t.Client.UnregisterChaincodeEvent(reg)

	// 数据处理，将m序列化成为json对象
	b, err := json.Marshal(m)
	if err != nil {
		return "", fmt.Errorf("指定的m对象序列化时发生错误")
	}

	//调用API生成请求
	req := channel.Request{
		ChaincodeID: t.ChaincodeID,
		Fcn:         "add",
		Args:        [][]byte{b, []byte(eventID), []byte(m.Key)}, //与链代码的参数顺序有关，第一个是value， 第三个是key
	}

	//执行请求
	response, err := t.Client.Execute(req)
	if err != nil {
		return "", err
	}

	err = eventResult(notifier, eventID)
	if err != nil {
		return "", err
	}

	return string(response.TransactionID), nil
}

func (t *ServiceSetup) FindInfo(Key string) ([]byte, error) {

	//生成请求

	req := channel.Request{
		ChaincodeID: t.ChaincodeID,
		Fcn:         "queryInfo",
		Args:        [][]byte{[]byte(Key)},
	}
	//执行请求
	response, err := t.Client.Query(req)
	if err != nil {
		return []byte{0x00}, err
	}

	return response.Payload, nil
}
func (t *ServiceSetup) FindInfoDetails(Key string) ([]byte, error) {

	//生成请求

	req := channel.Request{
		ChaincodeID: t.ChaincodeID,
		Fcn:         "queryInfo",
		Args:        [][]byte{[]byte(Key)},
	}
	//执行请求
	response, err := t.Client.Query(req)
	if err != nil {
		return []byte{0x00}, err
	}

	fmt.Println("Proposal.TxnID ======   ", response.Proposal.TxnID)

	return response.Payload, nil
}

func (t *ServiceSetup) FindInfoBy(Value, Test string) ([]byte, error) {

	queryString := fmt.Sprintf("{\"selector\":{ \"Value\":\"%s\", \"test\":\"%s\"}}", Value, Test)

	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "queryInfoBy", Args: [][]byte{[]byte(queryString)}}
	respone, err := t.Client.Query(req)
	if err != nil {
		return []byte{0x00}, err
	}

	return respone.Payload, nil
}

//-----------------------------------------------------------------------------------------//

func (t *ServiceSetup) MeasureSave(m Measure) (string, error) {

	eventID := "eventAdd"
	reg, notifier := regitserEvent(t.Client, t.ChaincodeID, eventID)
	defer t.Client.UnregisterChaincodeEvent(reg)

	// 将m对象序列化成为字节数组
	b, err := json.Marshal(m)
	if err != nil {
		return "", fmt.Errorf("指定的m对象序列化时发生错误")
	}

	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "add", Args: [][]byte{b, []byte(eventID), []byte(m.Id)}}
	respone, err := t.Client.Execute(req)
	if err != nil {
		return "", err
	}

	err = eventResult(notifier, eventID)
	if err != nil {
		return "", err
	}

	return string(respone.TransactionID), nil
}

func (t *ServiceSetup) MeasureFindInfo(Key string) ([]byte, error) {

	req := channel.Request{
		ChaincodeID: t.ChaincodeID,
		Fcn:         "queryInfo",
		Args:        [][]byte{[]byte(Key)},
	}

	respone, err := t.Client.Query(req)
	if err != nil {
		return []byte{0x00}, err
	}

	return respone.Payload, nil
}

func (t *ServiceSetup) MeasureFindInfoBy(No string) ([]byte, error) {

	queryString := fmt.Sprintf("{\"selector\":{ \"No\":\"%s\"}}", No)

	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "queryInfoBy", Args: [][]byte{[]byte(queryString)}}
	respone, err := t.Client.Query(req)
	if err != nil {
		return []byte{0x00}, err
	}

	return respone.Payload, nil
}
