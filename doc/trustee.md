## Escrow service

Official golang implementation of the Themis Escrow protocol.

## Building the source

Building escrow requires both a Go (version 1.7 or later) and a C compiler.
You can install them using your favourite package manager.
Once the dependencies are installed, run

    make gescrow

or, to build the full suite of utilities:

    make all
    
## Running escrow

Once you build escrow success, the runnable files is inside `./build/bin` directory.

By far we can running the escrow service. To do so:
```
$ gescrow -datadir /path/to/keystore/file -endpoint 192.168.1.109:8089 -nodes 45.249.245.140:8546
```

Parameter explain:
* datadir: escrow service account's keystore file path.
* endpoint: escrow rpc server endpoint.
* nodes: the themis node's rpc endpoint which escrow service connect. Currently, only support one node connection.

This command will:

 * Start escrow service in terminal, you need input the password of your keystore, because the escrow service need your address's private key to decrypt the secret fragment which mandate on the themischain.
 * Start a build-in monitor, listen to all the related order happened on the smart contract. And cache the order related data in memory, like: encrypted secret fragment, arbitrate result. 
 * Start up Escrow's built-in RPC Server, listening the request to withdraw the decrypted secret fragment
   of specific trade order.
 