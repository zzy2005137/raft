
 package main

 import (
	 "os"
	 "fmt"
 
	 "github.com/raft/fabric_raft/sdkInit"
 
	 // "github.com/raft/fabric_raft/model"
	 "github.com/raft/fabric_raft/service"
	 "github.com/raft/fabric_raft/web"
	 "github.com/raft/fabric_raft/web/controller"


 )
 
 const (
	 org1CfgPath = "../sdkConfig/org1_config.yaml"
	 org2CfgPath = "../sdkConfig/org2_config.yaml"
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
	 sdk1, err := sdkInit.SetupSDK(org1CfgPath, initialized)
	 if err != nil {
		 fmt.Printf(err.Error())
		 return
	 }
	 defer sdk1.Close()

	 channelClient, _ := sdkInit.GetClient(sdk1, initInfo1)  //test
	 fmt.Println(channelClient)
 
	 //==============================================//
	 //             服务层功能测试
	 //==============================================//
	 serviceSetup := service.ServiceSetup{
		 ChaincodeID: initInfo1.ChaincodeID,
		 Client:      channelClient,
	 }
 
	//  m := service.Mechanic{
	// 	 Key:   "zs",
	// 	 Value: "101",
	// 	 Test:  "test",
	//  }
 
	//  fmt.Println(m)

 
	 fmt.Println("===========服务层功能测试完成================")
 
	 fmt.Println("===========web 服务启动===================")
 
	 app := controller.Application{
	 	Fabric: &serviceSetup,
	 }
	 web.WebStart(&app)
 
 }