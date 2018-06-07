package main

import (
	"flag"
	"log"
	"github.com/themis-network/go-themis/trustee"
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

	var trustee = trustee.New(config)
	trustee.Start()
	//trustee.GetContractData()
}
