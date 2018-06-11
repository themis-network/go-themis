package main

import (
	"flag"
	"github.com/themis-network/go-themis/trustee"
)

var(
	datadir = flag.String("datadir", "", `--datadir path_to_data`)
)

func main(){

	flag.Parse()

	config := trustee.Config{
		DataDir: *datadir,
	}

	var trusteeServer = trustee.New(config)
	trusteeServer.Start()

	//trustee.GetContractData()
}
