package main

import (
	"github.com/themis-network/go-themis/trustee"
	"log"
	"gopkg.in/urfave/cli.v1"
	"os"
	"fmt"
)

var(
	app = cli.NewApp()

	dataDirFlag = cli.StringFlag{
		Name:  "datadir",
		Usage: "Data directory for keystore",
	}
	endpointFlag = cli.StringFlag{
		Name:  "endpoint",
		Usage: "ip:port, eg. 192.168.1.102:8090",
	}
	nodesFlag = cli.StringFlag{
		Name:  "nodes",
		Usage: "full node ws endpoint, ip:port, eg. 192.168.1.102:8090",
	}
)

func init(){
	app.Name = "Trustee service"
	app.Usage = "gtrustee [options]"
	app.Copyright = "Copyright 2017-2018 The go-themis Authors"
	app.Version = "0.5.1"
	app.Action = gtrustee
	flags := []cli.Flag{
		dataDirFlag,
		endpointFlag,
		nodesFlag,
	}
	app.Flags = append(app.Flags, flags...)
}

func gtrustee(ctx *cli.Context){

	if !ctx.GlobalIsSet(endpointFlag.Name){
		log.Fatal("Error, need --endpoint ip:port")
	}

	config := trustee.Config{
		DataDir: ctx.GlobalString(dataDirFlag.Name),
		Endpoint: ctx.GlobalString(endpointFlag.Name),
		Nodes: ctx.GlobalString(nodesFlag.Name),
	}

	var trusteeServer = trustee.New(config)
	trusteeServer.Start()
}

func main(){
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
