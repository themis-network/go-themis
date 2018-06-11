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

type TrusteeNode struct{

	secrets map[int64]string //每个订单对应的解密后的密钥碎片

	arbitrateEvents chan ArbitrateEvent

	orderWinner map[int64]*big.Int

	stop chan struct{} //Channel to wait for termination notifications

	config Config

	privKey keystore.Key

	contractClient *ContractClient
}

type ArbitrateEvent struct{
	orderId string
	winner string
}

func New(c Config) (t *TrusteeNode){

	//var pass string = "123456"
	var pass string

	fmt.Printf("Enter masked password: ")
	maskedPassword, _ := gopass.GetPasswdMasked() // Masked
	pass = string(maskedPassword)

	blob1, err := ioutil.ReadFile(c.DataDir)
	if err != nil {
		log.Fatal("failed to read freshly persisted node key: ", err)
	}
	log.Println("loaded keystore file...")

	privKey, err:= keystore.DecryptKey(blob1, pass)
	if err != nil {
		log.Fatal("failed to DecryptKey: ", err)
	}

	contractClient, err := getContractClient()
	if err != nil {
		log.Fatal("failed to get contractClient")
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

//启动服务
func (t *TrusteeNode) Start(){

	t.startApiServer()
	t.monitor()

	t.wait()
}

func (t *TrusteeNode) wait(){
	<- t.stop
}

func (t *TrusteeNode) Stop(){

}


//解密托管密钥
func  (t *TrusteeNode) decrypt(secret string) (string, error){

	priv := ecies.ImportECDSA(t.privKey.PrivateKey)

	//字符转成字节数组
	bytesSercret, _ := hex.DecodeString(secret)
	rawMsg, err := priv.Decrypt(bytesSercret, nil, nil)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(rawMsg), nil
}
