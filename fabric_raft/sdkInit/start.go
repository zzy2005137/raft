/**
  author: kevin
*/

package sdkInit

import (
	"fmt"

	mspclient "github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/fab/ccpackager/gopackager"
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/common/cauthdsl"
	// "github.com/hyperledger/fabric-protos-go/common"
)

const ChaincodeVersion = "1.0"

func SetupSDK(ConfigFile string, initialized bool) (*fabsdk.FabricSDK, error) {

	if initialized {
		return nil, fmt.Errorf("Fabric SDK已被实例化")
	}

	sdk, err := fabsdk.New(config.FromFile(ConfigFile))
	if err != nil {
		return nil, fmt.Errorf("实例化Fabric SDK失败: %v", err)
	}

	fmt.Println("Fabric SDK初始化成功")
	return sdk, nil
}

// 创建通道并将指定的Peers加入
func CreateChannel(sdk *fabsdk.FabricSDK, info *InitInfo) error {

	clientContext := sdk.Context(fabsdk.WithUser(info.OrgAdmin), fabsdk.WithOrg(info.OrgName))
	if clientContext == nil {
		return fmt.Errorf("根据指定的组织名称与管理员创建资源管理客户端Context失败")
	}

	// New returns a resource management client instance.
	resMgmtClient, err := resmgmt.New(clientContext)
	if err != nil {
		return fmt.Errorf("根据指定的资源管理客户端Context创建通道管理客户端失败: %v", err)
	}

	// New creates a new Client instance
	mspClient, err := mspclient.New(sdk.Context(), mspclient.WithOrg(info.OrgName))
	if err != nil {
		return fmt.Errorf("根据指定的 OrgName 创建 Org MSP 客户端实例失败: %v", err)
	}

	//  Returns: signing identity
	adminIdentity, err := mspClient.GetSigningIdentity(info.OrgAdmin)
	if err != nil {
		return fmt.Errorf("获取指定id的签名标识失败: %v", err)
	}

	// SaveChannelRequest holds parameters for save channel request
	channelReq := resmgmt.SaveChannelRequest{ChannelID: info.ChannelID,
		ChannelConfigPath: info.ChannelConfig,
		SigningIdentities: []msp.SigningIdentity{adminIdentity}}
	// save channel response with transaction ID
	_, err = resMgmtClient.SaveChannel(channelReq,
		resmgmt.WithRetry(retry.DefaultResMgmtOpts),
		resmgmt.WithOrdererEndpoint(info.OrdererOrgName))
	if err != nil {
		return fmt.Errorf("创建应用通道失败: %v", err)
	}

	fmt.Println("通道已成功创建，")

	info.OrgResMgmt = resMgmtClient
	/*
		// allows for peers to join existing channel with optional custom options (specific peers, filtered peers). If peer(s) are not specified in options it will default to all peers that belong to client's MSP.
		err = info.OrgResMgmt.JoinChannel(info.ChannelID, resmgmt.WithRetry(retry.DefaultResMgmtOpts), resmgmt.WithOrdererEndpoint(info.OrdererOrgName))
		if err != nil {
			return fmt.Errorf("Peers加入通道失败: %v", err)
		}

		fmt.Println("peers 已成功加入通道.")
	*/
	return nil
}

func InstallAndInstantiateCC(sdk *fabsdk.FabricSDK, info *InitInfo) (*channel.Client, error) {
	fmt.Println("开始安装链码......")
	// creates new go lang chaincode package
	ccPkg, err := gopackager.NewCCPackage(info.ChaincodePath, info.ChaincodeGoPath)
	if err != nil {
		return nil, fmt.Errorf("创建链码包失败: %v", err)
	}

	// contains install chaincode request parameters
	installCCReq := resmgmt.InstallCCRequest{Name: info.ChaincodeID, Path: info.ChaincodePath, Version: ChaincodeVersion, Package: ccPkg}
	// allows administrators to install chaincode onto the filesystem of a peer
	_, err = info.OrgResMgmt.InstallCC(installCCReq, resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	if err != nil {
		return nil, fmt.Errorf("安装链码失败: %v", err)
	}

	fmt.Println("指定的链码安装成功")
	fmt.Println("开始实例化链码......")

	//  returns a policy that requires one valid
	ccPolicy := cauthdsl.SignedByAnyMember([]string{"org1.kevin.kongyixueyuan.com"})

	instantiateCCReq := resmgmt.InstantiateCCRequest{Name: info.ChaincodeID, Path: info.ChaincodePath, Version: ChaincodeVersion, Args: [][]byte{[]byte("init")}, Policy: ccPolicy}
	// instantiates chaincode with optional custom options (specific peers, filtered peers, timeout). If peer(s) are not specified
	_, err = info.OrgResMgmt.InstantiateCC(info.ChannelID, instantiateCCReq, resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	if err != nil {
		return nil, fmt.Errorf("实例化链码失败: %v", err)
	}

	fmt.Println("链码实例化成功")

	clientChannelContext := sdk.ChannelContext(info.ChannelID, fabsdk.WithUser(info.UserName), fabsdk.WithOrg(info.OrgName))
	// returns a Client instance. Channel client can query chaincode, execute chaincode and register/unregister for chaincode events on specific channel.
	channelClient, err := channel.New(clientChannelContext)
	if err != nil {
		return nil, fmt.Errorf("创建应用通道客户端失败: %v", err)
	}

	fmt.Println("通道客户端创建成功，可以利用此客户端调用链码进行查询或执行事务.")

	return channelClient, nil
}

//默认加入所有属于组织的节点
func Join(sdk *fabsdk.FabricSDK, orgAdmin string, orgName string, ordererID string, channelID string) error {
	// join Channel
	orgContext := sdk.Context(fabsdk.WithUser(orgAdmin), fabsdk.WithOrg(orgName))

	orgResMgmt, err := resmgmt.New(orgContext)
	if err != nil {
		fmt.Printf("failed to create org1ResMgmt: %s \n", err)
	}
	fmt.Println("开始加入频道")
	if err = orgResMgmt.JoinChannel(channelID, resmgmt.WithRetry(retry.DefaultResMgmtOpts), resmgmt.WithOrdererEndpoint(ordererID)); err != nil {
		fmt.Printf("%s peers 加入频道失败: %s \n", orgName, err)
	}
	fmt.Println(orgName, "加入频道成功")

	return nil

}

func InstallCC(sdk *fabsdk.FabricSDK, info *InitInfo) error {
	fmt.Println("开始安装链码......")
	// creates new go lang chaincode package
	ccPkg, err := gopackager.NewCCPackage(info.ChaincodePath, info.ChaincodeGoPath)
	if err != nil {
		return fmt.Errorf("创建链码包失败: %v", err)
	}

	// contains install chaincode request parameters
	installCCReq := resmgmt.InstallCCRequest{
		Name:    info.ChaincodeID,
		Path:    info.ChaincodePath,
		Version: ChaincodeVersion,
		Package: ccPkg,
	}
	// allows administrators to install chaincode onto the filesystem of a peer
	_, err = info.OrgResMgmt.InstallCC(installCCReq, resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	if err != nil {
		return fmt.Errorf("安装链码失败: %v", err)
	}

	fmt.Println("指定的链码安装成功")

	return nil
}

func InstantiateCC(sdk *fabsdk.FabricSDK, info *InitInfo) (*channel.Client, error) {

	fmt.Println("开始实例化链码......")

	//  returns a policy that requires one valid
	org1OrOrg2 := "OR('Org1MSP.member','Org2MSP.member')"
	// ccPolicy, err := genPolicy(org1OrOrg2)
	ccPolicy, err := cauthdsl.FromString(org1OrOrg2)
	if err != nil {
		return nil, err
	}
	// args := packArgs([]string{"init", "zs", "101", "test","123"})
	args := packArgs([]string{"init"})

	instantiateCCReq := resmgmt.InstantiateCCRequest{
		Name:    info.ChaincodeID,
		Version: ChaincodeVersion,
		Path:    info.ChaincodePath,
		// Args: [][]byte{[]byte("init")},
		Args:   args,
		Policy: ccPolicy,
	}
	// instantiates chaincode with optional custom options (specific peers, filtered peers, timeout). If peer(s) are not specified
	//reqPeers := resmgmt.WithTargetEndpoints(peer)
	_, err = info.OrgResMgmt.InstantiateCC(info.ChannelID,
		instantiateCCReq,
		resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	if err != nil {
		return nil, fmt.Errorf("实例化链码失败: %v", err)
	}

	fmt.Println("链码实例化成功")

	clientChannelContext := sdk.ChannelContext(info.ChannelID, fabsdk.WithUser(info.UserName), fabsdk.WithOrg(info.OrgName))
	// returns a Client instance. Channel client can query chaincode, execute chaincode and register/unregister for chaincode events on specific channel.
	channelClient, err := channel.New(clientChannelContext)
	if err != nil {
		return nil, fmt.Errorf("创建应用通道客户端失败: %v", err)
	}

	fmt.Println("通道客户端创建成功，可以利用此客户端调用链码进行查询或执行事务.")

	return channelClient, nil
}
func NoCreateChannel(sdk *fabsdk.FabricSDK, info *InitInfo) error {

	clientContext := sdk.Context(fabsdk.WithUser(info.OrgAdmin), fabsdk.WithOrg(info.OrgName))
	if clientContext == nil {
		return fmt.Errorf("根据指定的组织名称与管理员创建资源管理客户端Context失败")
	}

	// New returns a resource management client instance.
	resMgmtClient, err := resmgmt.New(clientContext)
	if err != nil {
		return fmt.Errorf("根据指定的资源管理客户端Context创建通道管理客户端失败: %v", err)
	}

	info.OrgResMgmt = resMgmtClient
	return nil
}

func packArgs(paras []string) [][]byte {
	var args [][]byte
	for _, k := range paras {
		args = append(args, []byte(k))
	}
	return args
}


func GetClient(sdk *fabsdk.FabricSDK, info *InitInfo) (*channel.Client, error) {
	
	clientChannelContext := sdk.ChannelContext(info.ChannelID, fabsdk.WithUser(info.UserName), fabsdk.WithOrg(info.OrgName))
	// returns a Client instance. Channel client can query chaincode, execute chaincode and register/unregister for chaincode events on specific channel.
	channelClient, err := channel.New(clientChannelContext)
	if err != nil {
		return nil, fmt.Errorf("创建应用通道客户端失败: %v", err)
	}

	fmt.Println("另一个通道客户端创建成功，可以利用此客户端调用链码进行查询或执行事务.")

	return channelClient, nil
}