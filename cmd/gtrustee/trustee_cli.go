package main

import (
	"flag"
	"github.com/themis-network/go-themis/trustee"
	"log"
)

var(
	datadir = flag.String("datadir", "", `--datadir path_to_data`)
	endpoint = flag.String("endpoint", "", `--endpoint ip:port`)
)

func main(){

	flag.Parse()

	if *endpoint == "" {
		log.Fatal("Error, need --endpoint ip:port")
	}

	config := trustee.Config{
		DataDir: *datadir,
		Endpoint: *endpoint,
	}

	var trusteeServer = trustee.New(config)
	trusteeServer.Start()

	//trustee.GetContractData()
}
