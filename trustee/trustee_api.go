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
	vhosts = []string{"*"}
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

	rpc.StartHTTPEndpoint(t.config.Endpoint, apis, modules, cors, vhosts)
}

// NewPublicWeb3API creates a new Web3Service instance
func NewTrusteeAPI(t *TrusteeNode) *TrusteeAPI {
	return &TrusteeAPI{
		trusteeNode: t,
	}
}

/**
 GetDecryptSecret API, RPC "method":"trustee_getDecryptSecret"
1. verify the order's arbitrate result
2. try get decrypt fragment from map. if fail, get the fragment from contract, then decrypt it
 */
func (t *TrusteeAPI) GetDecryptSecret(orderId int64) (string, error){
	log.Println("Request trustee_getDecryptSecret, orderId:", orderId)

	if v, ok := t.trusteeNode.secrets[orderId]; ok {
		return v, nil
	}

	var winner *big.Int
	if v, ok := t.trusteeNode.orderWinner[orderId]; ok {
		winner = v
	}else {
		w, err := t.trusteeNode.getWinner(orderId)
		if err != nil {
			winner = nil
			log.Println("get winner error")
		}else {
			winner = w
		}
	}

	if winner == nil || winner.Int64() == 0{
		//no winner error
		log.Println("no winner error, orderid: ", orderId)
		return "", &noWinnerError{"no winner error"}
	}

	log.Println("winner is: ", winner.Int64())

	sectet, err := t.trusteeNode.getFragment(orderId, winner)
	log.Println("secret from contract: ", sectet)
	decSectet, err:= t.trusteeNode.decrypt(sectet)

	if err != nil {
		log.Println("decrypt secret error, ", err)
		return "", &decryptError{fmt.Sprintf("decrypt secret error, %v", err)}
	}
	return decSectet, nil
}


func errorJson(code int, message string) string{
	errorJson := &jsonError{Code: code, Message: message}
	jsons, errs := json.Marshal(errorJson)
	if errs != nil {
		fmt.Println(errs.Error())
		return ""
	}
	return string(jsons)
}