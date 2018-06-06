.PHONY: build cluster test

default: build

build:
	go build -o bin/swarm    github.com/themis-network/go-themis/cmd/swarm
	go build -o bin/geth     github.com/themis-network/go-themis/cmd/geth
	go build -o bin/bootnode github.com/themis-network/go-themis/cmd/bootnode

cluster: build
	scripts/boot-cluster.sh

test:
	go test -v github.com/themis-network/go-themis/swarm/...
