package main

import (
	"flag"
	"github.com/themis-network/go-themis/trustee"
	"log"
)

var(
	datadir = flag.String("datadir", "", `--datadir path_to_data`)
)

func main(){

	flag.Parse()
	log.Println(*datadir)

	config := trustee.Config{
		DataDir: *datadir,
	}

	var trusteeServer = trustee.New(config)
	trusteeServer.Start()

	//trustee.GetContractData()
}
