package trustee

import (
	"sync"
	"github.com/themis-network/go-themis/accounts/keystore"
	"github.com/themis-network/go-themis/crypto/ecies"
	"log"
	"encoding/hex"
	"io/ioutil"
	"math/big"
	"fmt"
	"github.com/themis-network/go-themis/trustee/gopass"
)

var(
	wg sync.WaitGroup
)

/**
	Trustee node service, response for secret decrypt request
 */
type TrusteeNode struct{

	secrets map[int64]string //order:secret

	arbitrateEvents chan ArbitrateEvent

	orderWinner map[int64]*big.Int

	stop chan struct{} //Channel to wait for termination notifications

	config Config //TrusteeNode config

	privKey keystore.Key //trustee node's private key

	contractClient *ContractClient
}

type ArbitrateEvent struct{
	orderId string
	winner string
}

//return new trustee instance
func New(c Config) (t *TrusteeNode){
	//var pass string = "123456"
	var pass string

	fmt.Printf("Enter keystore password: ")
	maskedPassword, err := gopass.GetPasswdMasked() // Masked
	if err != nil {
		log.Fatal("readPassword error: ", err)
	}
	pass = string(maskedPassword)

	keyStoreBlob, err := ioutil.ReadFile(c.DataDir)
	if err != nil {
		log.Fatal("failed to read keystore file: ", err)
	}
	log.Println("loaded keystore file...")

	privKey, err:= keystore.DecryptKey(keyStoreBlob, pass)
	if err != nil {
		log.Fatal("failed to DecryptKey: ", err)
	}

	contractClient, err := getContractClient()
	if err != nil {
		log.Fatal("failed to get contractClient: ", err)
	}

	var trustee = &TrusteeNode{
		secrets : make(map[int64]string),
		orderWinner : make(map[int64]*big.Int),
		config : c,
		stop: make(chan struct{}),
		privKey: *privKey,
		contractClient: contractClient,
		}
	return trustee
}

//start TrusteeNode service
func (t *TrusteeNode) Start(){

	t.startApiServer()
	t.monitor()

	t.wait()
}

//wait for TrusteeNode service stop
func (t *TrusteeNode) wait(){
	<- t.stop
}

//stop TrusteeNode service
func (t *TrusteeNode) Stop(){

}


//decrypt secret hold by trustee
func  (t *TrusteeNode) decrypt(secret string) (string, error){

	priv := ecies.ImportECDSA(t.privKey.PrivateKey)

	//字符转成字节数组
	bytesSecret, _ := hex.DecodeString(secret)
	rawMsg, err := priv.Decrypt(bytesSecret, nil, nil)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(rawMsg), nil
}
