

package main

import (
	/*
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	*/
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/peer"

	"fmt"

//	"encoding/json"
	"bytes"
)


type MechanicChaincode struct {

}


func (t *MechanicChaincode) Init(stub shim.ChaincodeStubInterface) peer.Response{

	return shim.Success(nil)
}

func (t *MechanicChaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response{
	fun, args := stub.GetFunctionAndParameters()

	/*
	var result string
	
	var result Mechanic
	var err error
	
	if fun == "set"{
		return set(stub, args)
	}else if fun == "get"{
		return get(stub, args)
	}
	*/
	if fun == "add"{
		return t.add(stub, args)
	}else if fun == "queryInfo"{
		return t.queryInfo(stub,args)
	}else if fun =="queryInfoBy"{
		return t.queryInfoBy(stub,args)
	}

	return shim.Error("指定的函数名称错误")
	/*
	if err != nil{
		return shim.Error(err.Error())
	}
	return t.addEdu(stub, args)
	return shim.Success([]byte(result))
	*/
}


func Put(stub shim.ChaincodeStubInterface, b []byte, Id string) ([]byte, bool) {

	/*
	b, err := json.Marshal(m)
	if err != nil {
		return nil, false
	}
	*/
	// 保存状态
	err := stub.PutState(Id, b)
	if err != nil {
		return nil, false
	}

	return b, true
}

func (t *MechanicChaincode) add(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	if len(args) != 3{
		return shim.Error("给定的参数个数不符合要求add" )
	}

	/*
	var m Mechanic
	err := json.Unmarshal([]byte(args[0]), &m)
	if err != nil {
		return shim.Error("反序列化信息时发生错误")
	}
	*/

	_, bl := Put(stub, []byte(args[0]), args[2])
	if !bl {
		return shim.Error("保存信息时发生错误")
	}

	err := stub.SetEvent(args[1], []byte{})
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success([]byte("信息添加成功"))
}

func (t *MechanicChaincode) queryInfo(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 1 {
		return shim.Error("给定的参数个数不符合要求queryInfo")
	}

	// 根据Key查询m状态
	b, err := stub.GetState(args[0])
	if err != nil {
		return shim.Error("根据Key查询信息失败")
	}

	if b == nil {
		return shim.Error("根据Key没有查询到相关的信息")
	}

	// 对查询到的状态进行反序列化
	/*
	var m Mechanic
	err = json.Unmarshal(b, &m)
	if err != nil {
		return  shim.Error("反序列化m信息失败")
	}

	// 返回
	result, err := json.Marshal(m)
	if err != nil {
		return shim.Error("序列化edu信息时发生错误")
	}
	*/
	return shim.Success(b)
}

// 根据Value及Test查询信息
// args: Value, Test
func (t *MechanicChaincode) queryInfoBy(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	if len(args) != 1 {
		return shim.Error("给定的参数个数不符合要求queryInfoBy")
	}
	/*
	Value := args[0]
	Test := args[1]
	*/
	// 拼装CouchDB所需要的查询字符串(是标准的一个JSON串)
	// queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"eduObj\", \"CertNo\":\"%s\"}}", CertNo)
	//queryString := fmt.Sprintf("{\"selector\":{ \"Value\":\"%s\", \"test\":\"%s\"}}",  Value, Test)

	// 查询数据
	result, err := getByQueryString(stub, args[0])
	if err != nil {
		return shim.Error("根据证书编号及姓名查询信息时发生错误")
	}
	if result == nil {
		return shim.Error("根据指定的证书编号及姓名没有查询到相关的信息")
	}
	return shim.Success(result)
}

// 根据指定的查询字符串实现富查询
func getByQueryString(stub shim.ChaincodeStubInterface, queryString string) ([]byte, error) {

	resultsIterator, err := stub.GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer  resultsIterator.Close()

	// buffer is a JSON array containing QueryRecords
	var buffer bytes.Buffer

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString("-+-+-")
		}

		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		bArrayMemberAlreadyWritten = true
	}

	fmt.Printf("- getQueryResultForQueryString queryResult:\n%s\n", buffer.String())

	return buffer.Bytes(), nil

}

func main(){
	/*
	err := shim.Start(new(SimpleChaincode))
	*/
	err := shim.Start(new(MechanicChaincode))
	if err != nil{
		fmt.Printf("启动SimpleChaincode时发生错误: %s", err)
	}
}

/*
func set(stub shim.ChaincodeStubInterface, args []string)(string, error){

	if len(args) != 3{
		return "", fmt.Errorf("给定的参数个数不符合要求")
	}

	err := stub.PutState(args[0], []byte(args[1]))
	if err != nil{
		return "", fmt.Errorf(err.Error())
	}

	err = stub.SetEvent(args[2], []byte{})
	if err != nil {
		return "", fmt.Errorf(err.Error())
	}

	return string(args[0]), nil

}

func get(stub shim.ChaincodeStubInterface, args []string)(string, error){
	if len(args) != 1{
		return "", fmt.Errorf("给定的参数个数不符合要求")
	}
	result, err := stub.GetState(args[0])
	if err != nil{
		return "", fmt.Errorf("获取数据发生错误")
	}
	if result == nil{
		return "", fmt.Errorf("根据 %s 没有获取到相应的数据", args[0])
	}
	return string(result), nil

}
*/