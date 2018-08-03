package escrow

import (
	"sync"
	"github.com/themis-network/go-themis/accounts/keystore"
	"github.com/themis-network/go-themis/crypto/ecies"
	"encoding/hex"
	"io/ioutil"
	"fmt"
	"github.com/themis-network/go-themis/escrow/gopass"
	"github.com/themis-network/go-themis/common"
)

var(
	wg sync.WaitGroup
)

/**
	Escrow node service, response for secret decrypt request
 */
type EscrowNode struct{

	secrets map[int64]string //order:secret

	arbitrateEvents chan ArbitrateEvent

	orderWinner map[int64]uint32

	stop chan struct{} //Channel to wait for termination notifications

	config Config //EscrowNode config

	privKey keystore.Key //EscrowNode's private key

	contractClient *ContractClient

	escrowAddr common.Address
}

type ArbitrateEvent struct{
	orderId string
	winner string
}

//return new EscrowNode instance
func New(c Config) (t *EscrowNode){
	//var pass string = "123456"
	var pass string

	fmt.Printf("Enter keystore password: ")
	maskedPassword, err := gopass.GetPasswdMasked() // Masked
	if err != nil {
		logger.Fatal("readPassword error: ", err)
	}
	pass = string(maskedPassword)

	keyStoreBlob, err := ioutil.ReadFile(c.DataDir)
	if err != nil {
		logger.Fatal("failed to read keystore file: ", err)
	}
	logger.Println("loaded keystore file...")

	privKey, err:= keystore.DecryptKey(keyStoreBlob, pass)
	if err != nil {
		logger.Fatal("failed to DecryptKey: ", err)
	}

	contractClient, err := getContractClient(c.Nodes)
	if err != nil {
		logger.Fatal("failed to get contractClient: ", err)
	}

	var escrow = &EscrowNode{
		secrets : make(map[int64]string),
		orderWinner : make(map[int64]uint32),
		config : c,
		stop: make(chan struct{}),
		privKey: *privKey,
		contractClient: contractClient,
		escrowAddr: privKey.Address,
		}
	return escrow
}

//start escrow service
func (t *EscrowNode) Start(){

	t.startApiServer()
	t.monitor()

	t.wait()
}

//wait for escrow service stop
func (t *EscrowNode) wait(){
	<- t.stop
}

//stop escrow service
func (t *EscrowNode) Stop(){

}


//decrypt secret hold by escrow
func  (t *EscrowNode) decrypt(secret string) (string, error){

	priv := ecies.ImportECDSA(t.privKey.PrivateKey)

	//字符转成字节数组
	bytesSecret, _ := hex.DecodeString(secret)
	rawMsg, err := priv.Decrypt(bytesSecret, nil, nil)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(rawMsg), nil
}
