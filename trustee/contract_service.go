package trustee

import (
	"math/big"
	"fmt"
	"github.com/themis-network/go-themis/common"
	"github.com/themis-network/go-themis"
	"github.com/themis-network/go-themis/ethclient"
	"context"
	"log"
	"github.com/themis-network/go-themis/core/types"
	"encoding/binary"
	"github.com/themis-network/go-themis/cmd/stub"
	"github.com/themis-network/go-themis/accounts/abi/bind"
)

const (
	ContractAddr  = "AAa91587531b304B117e367bBAb75ecD9B77cE15" //trade contract address
	trusteeAddr = "3cbcd06204c1df807f942f9edab069934fc14140" //trustee's address
	rawurl = "ws://192.168.1.213:8546"
	judgeTopic = "0x15c344b2775b6729564ceb0bd0971860f1f1d150ba24d1e4791336e3de69a186"
	uploadSecretTopic = ""
)

var(
	contractAddress = common.HexToAddress(ContractAddr)
)

type Monitor struct{
	trustee TrusteeNode
}

type ContractClient struct{
	rawClient *ethclient.Client
	ctx context.Context
	trader *stub.TradeCaller
}


func getClient() (*ethclient.Client, error){
	client, _ := ethclient.Dial(rawurl)
	return client, nil
}

func GetContractData(){
	client, _ := getClient()
	ctx := context.Background()
	hash := common.BigToHash(big.NewInt(1))
	x, _ := client.StorageAt(ctx, contractAddress, hash, nil)
	fmt.Println(x)
}

func (t *TrusteeNode) monitor(){

	log.Println("start monitor")

	ctx := context.Background()
	contractAddress := common.HexToAddress(ContractAddr)

	rawClient, _ := getClient()

	query := ethereum.FilterQuery{
		//FromBlock: big.NewInt(1431798),
		Addresses: []common.Address{contractAddress},
		//Topics: [][]common.Hash{{common.HexToHash(judgeTopic)}},
	}

	var ch = make(chan types.Log)
	sub, err := rawClient.SubscribeFilterLogs(ctx, query, ch)

	if err != nil {
		log.Println("Subscribe error:", err)
		return
	}

	for {
		select {
		case err := <-sub.Err():
			log.Fatal(err)
		case eventLog := <-ch:
			fmt.Println("Log:", eventLog.Address.Hex(), eventLog.Data)

			for i := 0; i<len(eventLog.Topics); i++ {
				log.Println(eventLog.Topics[i].Hex())
			}
			t.processLog(eventLog)
		}
	}
	log.Println("finish monitor")
}

/**
 process receive event from contract
 */
func (t *TrusteeNode)processLog(eventLog types.Log){

	topic := eventLog.Topics[0].Hex()

	if topic == judgeTopic { //judge event
		length := len(eventLog.Data)
		orderId := BytesToInt64(eventLog.Data[length-8:])
		winner := eventLog.Topics[1].Big()
		juedge := eventLog.Topics[2].Hex()
		t.orderWinner[orderId] = winner

		log.Println("Process Log, event judgeTopic, orderId:{}, winner:{}, judge:{}", orderId, winner, juedge)

		secret, err := t.getFragment(orderId, winner)
		if err != nil {
			log.Println("Error, getFragment error: %v", err)
		}

		decrypt, err := t.decrypt(secret)
		if err != nil {
			log.Println("Error, Decrypt error: %v", err)
		}

		t.secrets[orderId] = decrypt
	}else if topic == uploadSecretTopic { //upload secret event
		//do nothing
	}
}

func BytesToInt64(buf []byte) int64 {
	return int64(binary.BigEndian.Uint64(buf))
}

//从合约中获取碎片
func (t *TrusteeNode) getFragment(order int64, user *big.Int) (string, error){

	from := common.HexToAddress(trusteeAddr)

	opts := &bind.CallOpts{
		Pending: true,
		From: from,
		Context: t.contractClient.ctx,
	}

	str, err := t.contractClient.trader.GetSecret(opts, big.NewInt(order), from, user)
	if err != nil{
		return "", nil
	}
	return str, nil
}

//获取订单仲裁结果
func (t *TrusteeNode) getWinner(order int64) (*big.Int, error){

	from := common.HexToAddress(trusteeAddr)

	opts := &bind.CallOpts{
		Pending: false,
		From: from,
		Context: t.contractClient.ctx,
	}

	winner, err := t.contractClient.trader.GetWinner(opts, big.NewInt(order))
	if err != nil{
		return nil, err
	}

	return winner, nil
}

func getContractClient() (*ContractClient, error){
	rawClient, err := ethclient.Dial(rawurl)
	if err != nil {
		return nil, err
	}

	c := context.Background()

	addr := common.HexToAddress(ContractAddr)
	t, err := stub.NewTradeCaller(addr, rawClient)
	if err != nil {
		return nil, err
	}

	contractClient := &ContractClient{
		rawClient: rawClient,
		ctx: c,
		trader: t,
	}
	return contractClient, nil
}
