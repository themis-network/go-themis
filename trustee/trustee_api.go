package trustee

import (
	"fmt"
	"github.com/themis-network/go-themis/rpc"
	"log"
	"encoding/json"
	"math/big"
)

var(
	modules = []string{"trustee"}
	cors = []string{"*"}
	vhosts = []string{"localhost"}
	endpoint = "127.0.0.1:8089"
)

// PublicWeb3API offers helper utils
type TrusteeAPI struct {
	trusteeNode *TrusteeNode
}

type jsonError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func (t *TrusteeNode) startApiServer(){

	log.Println("Start Trustee api service...")

	apis := []rpc.API{
		{
			Namespace: "trustee",
			Version:   "1.0",
			Service:   NewTrusteeAPI(t),
			Public:    true,
		},
	}

	rpc.StartHTTPEndpoint(endpoint, apis, modules, cors, vhosts)
}

// NewPublicWeb3API creates a new Web3Service instance
func NewTrusteeAPI(t *TrusteeNode) *TrusteeAPI {
	return &TrusteeAPI{
		trusteeNode: t,
	}
}

/**
 获取解密密钥碎片接口, RPC "method":"trustee_getDecryptSecret"
1. verify the order's arbitrate result
2. try get decrypt fragment from map. if fail, get the fragment from contract, then decrypt it
 */
func (t *TrusteeAPI) GetDecryptSecret(orderId int64) string {
	log.Println("GetDecryptSecret, orderId:", orderId)

	if v, ok := t.trusteeNode.secrets[orderId]; ok {
		return v
	}

	var winner *big.Int
	if v, ok := t.trusteeNode.orderWinner[orderId]; ok {
		winner = v
	}else {
		w, err := t.trusteeNode.getWinner(orderId)
		if err != nil {
			winner = nil
		}else {
			winner = w
		}
	}

	if winner == nil{
		//没有仲裁的winner，返回JSON
		return errorJson(-1, "This order has no winner") //byte[]转换成string 输出
	}

	log.Println("winner is: ", winner.Int64())

	sectet, err := t.trusteeNode.getFragment(orderId, winner)
	log.Println("secret from contract: ", sectet)
	decSectet, err:= t.trusteeNode.decrypt(sectet)

	if err != nil {
		log.Println("decrypt secret error: ", err)
		//todo return error
	}
	return decSectet
}

func errorJson(code int, message string) string{
	errorJson := &jsonError{Code: code, Message: message}
	jsons, errs := json.Marshal(errorJson) //转换成JSON返回的是byte[]
	if errs != nil {
		fmt.Println(errs.Error())
		return ""
	}
	return string(jsons) //byte[]转换成string 输出
}