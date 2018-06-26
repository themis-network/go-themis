package main

import (
	"log"
	"gopkg.in/urfave/cli.v1"
	"os"
	"fmt"
	"github.com/themis-network/go-themis/escrow"
)

var(
	app = cli.NewApp()

	dataDirFlag = cli.StringFlag{
		Name:  "datadir",
		Usage: "directory for keystore",
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
	app.Name = "Escrow service"
	app.Usage = "Escrow [options]"
	app.Copyright = "Copyright 2017-2018 The go-themis Authors"
	app.Version = "0.5.1"
	app.Action = escrow_start
	flags := []cli.Flag{
		dataDirFlag,
		endpointFlag,
		nodesFlag,
	}
	app.Flags = append(app.Flags, flags...)
}

func escrow_start(ctx *cli.Context){

	if !ctx.GlobalIsSet(endpointFlag.Name){
		log.Fatal("Error, need --endpoint ip:port")
	}

	config := escrow.Config{
		DataDir: ctx.GlobalString(dataDirFlag.Name),
		Endpoint: ctx.GlobalString(endpointFlag.Name),
		Nodes: ctx.GlobalString(nodesFlag.Name),
	}

	var escrowServer = escrow.New(config)
	escrowServer.Start()
}

func main(){
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
