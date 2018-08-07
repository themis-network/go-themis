package escrow

import (
	"math/big"
	"fmt"
	"github.com/themis-network/go-themis/common"
	"github.com/themis-network/go-themis"
	"github.com/themis-network/go-themis/ethclient"
	"context"
	"github.com/themis-network/go-themis/core/types"
	"github.com/themis-network/go-themis/cmd/stub"
	"github.com/themis-network/go-themis/accounts/abi/bind"
)

const (
	ContractAddr  = "368790bd905Adf4962e36BDec79EEF7dB7F0BEA3" //trade contract address
	testRawurl = "ws://192.168.1.213:8546"
	judgeTopic = "0x15c344b2775b6729564ceb0bd0971860f1f1d150ba24d1e4791336e3de69a186"
	uploadSecretTopic = "0x8a59d01dda427123e224b10a5103435e6a94ce386bd3d81052074263f9defce8"
	nodeProtocol = "ws://"
)

var(
	contractAddress = common.HexToAddress(ContractAddr)
)

type Monitor struct{
	escrow EscrowNode
}

type ContractClient struct{
	rawClient *ethclient.Client
	ctx context.Context
	traderCaller *stub.TradeCaller
	traderTransactor *stub.Trade
}


func getClient(nodeEndpoint string) (*ethclient.Client, error){
	nodeWsUrl := nodeProtocol + nodeEndpoint
	client, _ := ethclient.Dial(nodeWsUrl)
	return client, nil
}

func GetContractData(){
	client, _ := getClient("")
	ctx := context.Background()
	hash := common.BigToHash(big.NewInt(1))
	x, _ := client.StorageAt(ctx, contractAddress, hash, nil)
	fmt.Println(x)
}

func (t *EscrowNode) monitor(){

	logger.Println("start monitor")

	ctx := context.Background()
	contractAddress := common.HexToAddress(ContractAddr)

	rawClient, _ := getClient(t.config.Nodes)

	query := ethereum.FilterQuery{
		//FromBlock: big.NewInt(1431798),
		Addresses: []common.Address{contractAddress},
		//Topics: [][]common.Hash{{common.HexToHash(judgeTopic)}},
	}

	var ch = make(chan types.Log)
	sub, err := rawClient.SubscribeFilterLogs(ctx, query, ch)

	if err != nil {
		logger.Println("Subscribe error:", err)
		return
	}

	for {
		select {
		case err := <-sub.Err():
			logger.Fatal(err)
		case eventLog := <-ch:
			logger.Println("Log:", eventLog.Address.Hex(), eventLog.Data)

			for i := 0; i<len(eventLog.Topics); i++ {
				logger.Println("topic:", i, eventLog.Topics[i].Hex())
			}
			t.processLog(eventLog)
		}
	}
}

/**
 process receive event from contract
 */
func (t *EscrowNode)processLog(eventLog types.Log){

	topic := eventLog.Topics[0].Hex()

	if topic == judgeTopic { //judge event
		length := len(eventLog.Data)
		orderId := BytesToInt64(eventLog.Data[length-8:])
		winner := eventLog.Topics[1].Bytes()
		juedge := eventLog.Topics[2].Hex()
		t.orderWinner[orderId] = BytesToUint32(winner)

		logger.Println("Process Log, event judgeTopic, orderId:{}, winner:{}, judge:{}", orderId, winner, juedge)

		secret, err := t.getFragment(orderId, t.orderWinner[orderId])
		if err != nil {
			logger.Println("Error, getFragment error: ", err)
		}

		decrypt, err := t.decrypt(secret)
		if err != nil {
			logger.Println("Error, Decrypt error: ", err)
		}

		t.secrets[orderId] = decrypt
	}else if topic == uploadSecretTopic { //upload secret event

		orderIdBytes := eventLog.Topics[1].Bytes()
		orderId := big.NewInt(0)
		orderId.SetBytes(orderIdBytes)

		logger.Println("Process Log, event uploadSecretTopic, orderId:{}", orderId)

		status, err := t.contractClient.traderCaller.GetOrderStatus(t.getCallOpts(), orderId)
		logger.Println("order status", status)
		if err != nil{
			logger.Println("failed to get order status")
			return
		}

		if status == SecretUploaded {
			vb, vs, err := t.verify(orderId)
			if err != nil {
				logger.Println("verify error")
			}
			t.sendVerifyResults(orderId, vb, vs)
		}
	}else {
		logger.Println("Process Log, event unknow")
	}
}

func (t *EscrowNode)verify(orderId *big.Int) (bool, bool, error){

	from := t.escrowAddr

	opts := &bind.CallOpts{
		Pending: true,
		From: from,
		Context: t.contractClient.ctx,
	}

	buyer, err:= t.contractClient.traderCaller.GetOrderBuyer(opts, orderId)
	seller, err:= t.contractClient.traderCaller.GetOrderBuyer(opts, orderId)

	verfifyDataB, err := t.contractClient.traderCaller.GetVerifyData(opts, orderId, uint32(buyer.Uint64()))
	verfifyDataS, err := t.contractClient.traderCaller.GetVerifyData(opts, orderId, uint32(seller.Uint64()))

	fragmentBuyer, err := t.getFragment(orderId.Int64(), uint32(buyer.Uint64()))
	fragmentSeller, err := t.getFragment(orderId.Int64(), uint32(seller.Uint64()))

	vb, err:= verifyFragment(verfifyDataB, fragmentBuyer)
	vs, err := verifyFragment(fragmentSeller, verfifyDataS)
	if err != nil {
		return false, false, err
	}else{
		return vb, vs, nil
	}
}

func (t *EscrowNode)sendVerifyResults(orderId *big.Int, vb bool, vs bool){

	nonce, _ := t.contractClient.rawClient.PendingNonceAt(context.Background(), t.escrowAddr)
	gasPrice, _ := t.contractClient.rawClient.SuggestGasPrice(context.Background())

	logger.Println(nonce, gasPrice, )

	auth := bind.NewKeyedTransactor(t.privKey.PrivateKey)
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)     // in wei
	auth.GasLimit = uint64(500000) // in units
	auth.GasPrice = gasPrice

	tx, err:= t.contractClient.traderTransactor.SendVerifyResult(auth, orderId, vb, vs)
	if err!= nil {
		logger.Println(err)
	}
	logger.Println("sendVerifyResults", orderId, tx.Hash().Hex())
}


//从合约中获取碎片
func (t *EscrowNode) getFragment(order int64, user uint32) (string, error){

	from := t.escrowAddr

	opts := &bind.CallOpts{
		Pending: true,
		From: from,
		Context: t.contractClient.ctx,
	}

	str, err := t.contractClient.traderCaller.GetSecret(opts, big.NewInt(order), from, user)
	if err != nil{
		return "", nil
	}
	return str, nil
}

//获取订单仲裁结果
func (t *EscrowNode) getWinner(order int64) (uint32, error){

	from := t.escrowAddr

	opts := &bind.CallOpts{
		Pending: false,
		From: from,
		Context: t.contractClient.ctx,
	}

	winner, err := t.contractClient.traderCaller.GetWinner(opts, big.NewInt(order))
	if err != nil{
		return 0, err
	}

	return winner, nil
}

func getContractClient(nodeEndpoint string) (*ContractClient, error){

	logger.Println("Connecting to themis rpc service, nodeEndpoint:", nodeEndpoint)

	nodeWsUrl := nodeProtocol + nodeEndpoint
	rawClient, err := ethclient.Dial(nodeWsUrl)
	if err != nil {
		return nil, err
	}

	c := context.Background()

	addr := common.HexToAddress(ContractAddr)
	tCaller, err := stub.NewTradeCaller(addr, rawClient)
	tTran, err := stub.NewTrade(addr, rawClient)
	if err != nil {
		logger.Println("error NewTrade")
		return nil, err
	}

	contractClient := &ContractClient{
		rawClient: rawClient,
		ctx: c,
		traderCaller: tCaller,
		traderTransactor: tTran,
	}
	return contractClient, nil
}

func (t *EscrowNode) getCallOpts() * bind.CallOpts{
	from := t.escrowAddr
	opts := &bind.CallOpts{
		Pending: false,
		From: from,
		Context: t.contractClient.ctx,
	}
	return opts
}
