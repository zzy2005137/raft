/*
 * @Author: test1Tan
 * @GIthub: https://github.com/test1Tan-b-z
 * @Date: 2020-08-11 21:03:54
 * @LastEditors: Zhen
 * @LastEditTime: 2021-03-18 12:53
 */
package main

import (
	"os"

	// "time"
	"fmt"

	"github.com/raft/fabric_raft/sdkInit"

	// "github.com/raft/fabric_raft/model"
	"github.com/raft/fabric_raft/service"
	// "github.com/raft/fabric_raft/web"
	// "github.com/raft/fabric_raft/web/controller"

	// "github.com/raft/fabric_raft/goserial"
	"database/sql"
	"encoding/json"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

const (
	org1CfgPath = "./sdkConfig/org1_config.yaml"
	org2CfgPath = "./sdkConfig/org2_config.yaml"
	initialized = false
	ordererID   = "orderer1.example.com"
)

var (
	peer0Org1 = "peer0.org1.example.com"
	peer0Org2 = "peer0.org2.example.com"
)

func main() {

	//==============================================//
	//准备
	//==============================================//

	initInfo1 := &sdkInit.InitInfo{

		ChannelID:     "mychannel",
		ChannelConfig: os.Getenv("GOPATH") + "/src/github.com/raft/fabric_raft/channel-artifacts/mychannel.tx",

		OrgAdmin:       "Admin",
		OrgName:        "Org1",
		OrdererOrgName: "orderer1.example.com",

		ChaincodeID:     "mychaincode",
		ChaincodeGoPath: os.Getenv("GOPATH"),
		ChaincodePath:   "github.com/raft/fabric_raft/chaincode",
		UserName:        "User1",
	}

	initInfo2 := &sdkInit.InitInfo{

		ChannelID:     "mychannel",
		ChannelConfig: os.Getenv("GOPATH") + "/src/github.com/raft/fabric_raft/channel-artifacts/mychannel.tx",

		OrgAdmin:       "Admin",
		OrgName:        "Org2",
		OrdererOrgName: "orderer1.example.com",

		ChaincodeID:     "mychaincode",
		ChaincodeGoPath: os.Getenv("GOPATH"),
		ChaincodePath:   "github.com/raft/fabric_raft/chaincode/",
		UserName:        "User1",
	}

	sdk1, err := sdkInit.SetupSDK(org1CfgPath, initialized)
	if err != nil {
		fmt.Printf(err.Error())
		return
	}

	sdk2, err := sdkInit.SetupSDK(org2CfgPath, initialized)
	if err != nil {
		fmt.Printf(err.Error())
		return
	}
	defer sdk1.Close()
	defer sdk2.Close()

	//创建通道，加入节点
	err = sdkInit.CreateChannel(sdk1, initInfo1)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	err = sdkInit.NoCreateChannel(sdk2, initInfo2)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// join Channel
	err = sdkInit.Join(sdk1, initInfo1.OrgAdmin, initInfo1.OrgName, ordererID, initInfo1.ChannelID)
	err = sdkInit.Join(sdk2, initInfo2.OrgAdmin, initInfo2.OrgName, ordererID, initInfo2.ChannelID)

	// 安装和实例化链码
	err = sdkInit.InstallCC(sdk1, initInfo1)
	err = sdkInit.InstallCC(sdk2, initInfo2)
	if err != nil {
		fmt.Println("安装链码失败")
		fmt.Println(err.Error())
	}

	//实例化链代码
	channelClient, err := sdkInit.InstantiateCC(sdk1, initInfo1)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(channelClient)

	client2, _ := sdkInit.GetClient(sdk1, initInfo1)  //test
	fmt.Println(client2)

	//==============================================//
	//             服务层功能测试
	//==============================================//
	serviceSetup := service.ServiceSetup{
		ChaincodeID: initInfo1.ChaincodeID,
		Client:      channelClient,
	}

	m := service.Mechanic{
		Key:   "zs",
		Value: "101",
		Test:  "test",
	}

	fmt.Println(m)

	//添加数据
	msg, err := serviceSetup.Save(m)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("信息发布成功, 交易编号为: " + msg)
	}

	// 根据Key查询信息
	result, err := serviceSetup.FindInfo("zs")
	if err != nil {
		fmt.Println(err.Error())
	} else {
		var m service.Mechanic
		json.Unmarshal(result, &m)
		fmt.Println("根据Key查询信息成功：")
		fmt.Println(m)
	}

	// 富查询
	result, err = serviceSetup.FindInfoBy("101", "test")
	if err != nil {
		fmt.Println(err.Error())
	} else {
		var m service.Mechanic
		json.Unmarshal(result, &m)
		fmt.Println("富查询信息成功：")
		fmt.Println(m)
	}

	fmt.Println("===========服务层功能测试完成================")



}

//数据库测试，待用
func dbtest(serviceSetup service.ServiceSetup) {

	//=====================================
	//         连接数据库添加数据测试
	//=====================================

	mm := service.Measure{
		Id:    "null",
		No:    "2",
		Time:  "null",
		Ddata: [3]float32{0, 0, 0},
		Ldata: [3]float32{0, 0, 0},
	}

	db, err := sql.Open("mysql", "root:nuaalabMySQL@tcp(101.132.32.165:7249)/factoryNH")
	if err != nil {
		fmt.Printf("连接失败: %v\n", err)
	}
	defer db.Close()
	//连接测试
	err = db.Ping()
	if err != nil {
		fmt.Println(err)
	}

	stmt, err := db.Prepare("select * from measure_2")

	rows, err := stmt.Query()
	defer stmt.Close()
	//var time string
	//var d1,d2,d3,l1,l2,l3 float32
	i := 0
	for rows.Next() {
		i++

		if ers := rows.Scan(&mm.Time, &mm.Ddata[0], &mm.Ddata[1], &mm.Ddata[2], &mm.Ldata[0], &mm.Ldata[1], &mm.Ldata[2]); ers == nil {
			//fmt.Print("%s %s %s %s %s %s %s \n", time, d1, d2, d3, l1, l2, l3)
			mm.Id = strconv.Itoa(i)

			msg, err := serviceSetup.MeasureSave(mm)
			if err != nil {
				fmt.Println(err.Error())
			} else {
				fmt.Println("信息发布成功, 交易编号为: " + msg)
			}

		}
	}

	// 根据Id查询信息
	result, err := serviceSetup.MeasureFindInfo("1")
	if err != nil {
		fmt.Println(err.Error())
	} else {
		var mm service.Measure
		json.Unmarshal(result, &mm)
		fmt.Println("根据Id查询信息成功：")
		fmt.Println(mm)
	}

	// 富查询
	result, err = serviceSetup.MeasureFindInfoBy("2")
	if err != nil {
		fmt.Println(err.Error())
	} else {
		var m service.Measure
		resul := strings.Split(string(result), "-+-+-")
		for i := 0; i < len(resul); i++ {
			json.Unmarshal([]byte(resul[i]), &m)
			fmt.Println("富查询信息成功：")
			fmt.Println(m)
		}

	}
}
