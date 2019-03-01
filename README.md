# go_fabric_ledger
使用ledger 查询区块交易等


设置ledger  注册链码和实例化链码得时候操作 ledger定义

	"github.com/hyperledger/fabric-sdk-go/pkg/client/ledger"
	
	Ledger          *ledger.Client

//A、准备通道上下文
//B、创建分类帐客户端
org1AdminChannelContext := setup.sdk.ChannelContext(setup.ChannelID[i], fabsdk.WithUser(setup.OrgAdmin), fabsdk.WithOrg(setup.OrgName))

// Ledger client
setup.Ledger err = ledger.New(org1AdminChannelContext)
if err != nil {
return errors.WithMessage(err, "Failed to create new resource management client: %s")
}