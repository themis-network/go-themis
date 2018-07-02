## Go Themis

Official golang implementation of the Themis protocol.

Automated builds are available for stable releases and the unstable master branch.

## Building the source

Building gthemis requires both a Go (version 1.7 or later) and a C compiler.
You can install them using your favourite package manager.
Once the dependencies are installed, run

    make gthemis

or, to build the full suite of utilities:

    make all
    
## Executables

The go-themis project comes with several wrappers/executables found in the `cmd` directory.

| Command    | Description |
|:----------:|-------------|
| **`gthemis`** | Our main Themis CLI client. It is the entry point into the Themis network (main-, test- or private net), capable of running as a full node (default) archive node (retaining all historical state) or a light node (retrieving data live). It can be used by other processes as a gateway into the Themis network via JSON RPC endpoints exposed on top of HTTP, WebSocket and/or IPC transports. `gthemis --help` for command line options. |
| `abigen` | Source code generator to convert Themis contract definitions into easy to use, compile-time type-safe Go packages. It operates on plain with expanded functionality if the contract bytecode is also available. However it also accepts Solidity source files, making development much more streamlined.
| `bootnode` | Stripped down version of our Themis client implementation that only takes part in the network node discovery protocol, but does not run any of the higher level application protocols. It can be used as a lightweight bootstrap node to aid in finding peers in private networks. |
| `evm` | Developer utility version of the EVM (Ethereum Virtual Machine) that is capable of running bytecode snippets within a configurable environment and execution mode. Its purpose is to allow isolated, fine-grained debugging of EVM opcodes (e.g. `evm --code 60ff60ff --debug`). |
| `rlpdump` | Developer utility tool to convert binary RLP dumps (data encoding used by the Themis protocol both network as well as consensus wise) to user friendlier hierarchical representation (e.g. `rlpdump --hex CE0183FFFFFFC4C304050583616263`). |
| `swarm`    | swarm daemon and tools. This is the entrypoint for the swarm network. `swarm --help` for command line options and subcommands. See https://swarm-guide.readthedocs.io for swarm documentation. |
| `puppeth`    | a CLI wizard that aids in creating a new Themis network. |
| `gescrow`   | a CLI that listen for escrow related events and decrypt/send user's shard to service automatically. [read more](https://github.com/themis-network/go-themis/blob/master/doc/trustee.md)

## Running gthemis

Up to now, main themis net is not available, but you can run test themis net by
"--themisTestnet" parameter

### Full node on the test Themis network

By far the most common scenario is people wanting to simply interact with the test themis network:
create accounts; transfer funds; deploy and interact with contracts. For this particular use-case
the user doesn't care about years-old historical data, so we can fast-sync quickly to the current
state of the network. To do so:

```
$ gthemis --themisTestnet console
```

This command will:

 * Start Themis test net in fast sync mode (default, can be changed with the `--syncmode` flag), causing it to
   download more data in exchange for avoiding processing the entire history of the test Themis network.
 * Start up Gthemis's built-in interactive JavaScript console(via the trailing `console` subcommand). 
   This too is optional and if you leave it out you can always attach to an already running Gthemis instance
   with `gthemis attach`.
   
### Full miner node on the test Themis network

You have to run a node in full sync mode to be a miner/signer. We will init some signers
in genesis block, and any signer can propose adding a new signer or removing an old signer.
The propose will be accepted if N/2 + 1 signers vote for it.(N is the total amount of signers)
to start a signer/node:

```
$ gthemis --themisTestnet --syncmode "full" --unlock yourSignerAddress --password yourPWDFile --mine --minerthreads 1 --targetgaslimit 8000000 console
```

In order to unlock signer account, you have to supply password by manually enter it if you don't 
want to clear it into a file, otherwise using `--password` flag as above.

This command will:

* Start Themis test net in full mode.
* Unlock your signer address with the password file(password should right in first line) to sign/mine a block.
* Start 1 thread to do miner work.(you can special other numbers, but it is not too much difference between them
  cause it's a poa network)
* Start up Gthemis's built-in interactive JavaScript console(via the trailing `console` subcommand).
* Set the target gas limit of block to 100000000(same with gas limit in genesis block).

*Note: Although there are some internal protective measures to prevent transactions from crossing
over between the main network and test network, you should make sure to always use separate accounts
for play-money and real-money. Unless you manually move accounts, Gthemis will by default correctly
separate the two networks and will not make any accounts available between them.* 

**Note: please understand your signer account is always unlocked while signing/mining and don't open 
any rpc/ws based service on it. Hackers on the internet can steal your money through
this un security service very simply!**

Any signer can propose adding/removing a signer(include his/her self), to do so, you have to start Gthemis's
built-in interactive JavaScript console:

Propose adding a new signer:
```
$ clique.propose(newSigner, true)
```

Propose removing a old signer:
```
$ clique.propose(oldSigner, false)
```

Get signers at a certain block(latest block number will be used if number is omitted):
```
$ clique.getSigners(blockNumber) // blockNumber should be a hex string(e.g. "0x0")
```

Discard a propose(Any invalid propose will be discarded automatically, as well as those have been signed into a block.
So this will only be used before a valid propose been signed)
```
$ clique.discard(signer)
```

Get current vote details(latest block number will be used if number is omitted):
```
$ clique.getSnapshot(blockNumber) // blockNumber should be a hex string(e.g. "0x0")
```

For more details, you can deep into the javascript console!

*Note: Adding an existing signer or removing a non-existing signer is invalid, thus will be rejected.
More, there is a super signer hard coded in this themis test net, who can add/remove signer separately
(vote is not necessary). Any propose about super signer is regarded as invalid.*

### Static nodes
As an alternative to reconnect to themis test net nodes after disconnection, you can add a static-nodes.json file on your data 
directory. For example：

```go
[
    "enode://ff2f149dfd070c194ae427af55df7b3ee3c3c719c9d800391c90d2622811ea3d28b0718ddd9a88be4e69fa232ee9f2f0cca6cf7542b31dd7fcc7d5d1e8e602e9@47.93.163.113:30303",
    "enode://0d1e7b39d6b8a917f66eee24c682e3a67837aac404a0df20a4ec2beb461d53c9f12c0c67f433b78795bc351ed19d49455caaa8fedc328f09d390163542220194@45.249.245.140:30303",
    "enode://f2d7bc51877064f148560aee1a02bdaadfa6299a99c79633181597375d6cd3911570937445984c75ee458fd9d06d7db2c06b2a91ae4f49359d98d0c7a366f1fb@103.14.34.124:30303",
    "enode://d6d7fa6655bcacfb28dd2838a84d8261fb438949845497878d9baa0bccbcee78b74272aa7e747e0ca13488372f2f03bcda0614a356cc70de0f37326ef5731797@54.206.18.140:30303",
    "enode://55b60bc3df4d336b10b625e9e2ed5f86217fd6f8208cf98e13c6d96db142f9773167e2c4410561171fe5fc93783dd82f18c7e3359ecaced2e3ced69fbd49ea02@104.199.188.174:30303"
]
```

### Configuration

As an alternative to passing the numerous flags to the `gthemis` binary, you can also pass a configuration file via:

```
$ gthemis --config /path/to/your_config.toml
```

To get an idea how the file should look like you can use the `dumpconfig` subcommand to export your existing configuration:

```
$ gthemis --your-favourite-flags dumpconfig
```

### Programatically interfacing Gthemis nodes

As a developer, sooner rather than later you'll want to start interacting with Gthemis and the Themis
network via your own programs and not manually through the console. To aid this, Gthemis has built-in
support for a JSON-RPC based APIs. These can be exposed via HTTP, WebSockets and IPC (unix sockets on
unix based platforms, and named pipes on Windows).

The IPC interface is enabled by default and exposes all the APIs supported by Gthemis, whereas the HTTP
and WS interfaces need to manually be enabled and only expose a subset of APIs due to security reasons.
These can be turned on/off and configured as you'd expect.

HTTP based JSON-RPC API options:

  * `--rpc` Enable the HTTP-RPC server
  * `--rpcaddr` HTTP-RPC server listening interface (default: "localhost")
  * `--rpcport` HTTP-RPC server listening port (default: 8545)
  * `--rpcapi` API's offered over the HTTP-RPC interface (default: "eth,net,web3")
  * `--rpccorsdomain` Comma separated list of domains from which to accept cross origin requests (browser enforced)
  * `--ws` Enable the WS-RPC server
  * `--wsaddr` WS-RPC server listening interface (default: "localhost")
  * `--wsport` WS-RPC server listening port (default: 8546)
  * `--wsapi` API's offered over the WS-RPC interface (default: "eth,net,web3")
  * `--wsorigins` Origins from which to accept websockets requests
  * `--ipcdisable` Disable the IPC-RPC server
  * `--ipcapi` API's offered over the IPC-RPC interface (default: "admin,debug,eth,miner,net,personal,shh,txpool,web3")
  * `--ipcpath` Filename for IPC socket/pipe within the datadir (explicit paths escape it)

You'll need to use your own programming environments' capabilities (libraries, tools, etc) to connect
via HTTP, WS or IPC to a Gthemis node configured with the above flags and you'll need to speak [JSON-RPC](http://www.jsonrpc.org/specification)
on all transports. You can reuse the same connection for multiple requests!

**Note: Please understand the security implications of opening up an HTTP/WS based transport before
doing so! Hackers on the internet are actively trying to subvert Themis nodes with exposed APIs!
Further, all browser tabs can access locally running webservers, so malicious webpages could try to
subvert locally available APIs!**

## Contribution

Thank you for considering to help out with the source code! We welcome contributions from
anyone on the internet, and are grateful for even the smallest of fixes!

If you'd like to contribute to go-themis, please fork, fix, commit and send a pull request
for the maintainers to review and merge into the main code base.

Please make sure your contributions adhere to our coding guidelines:

 * Code must adhere to the official Go [formatting](https://golang.org/doc/effective_go.html#formatting) guidelines (i.e. uses [gofmt](https://golang.org/cmd/gofmt/)).
 * Code must be documented adhering to the official Go [commentary](https://golang.org/doc/effective_go.html#commentary) guidelines.
 * Pull requests need to be based on and opened against the `master` branch.
 * Commit messages should be prefixed with the package(s) they modify.
   * E.g. "eth, rpc: make trace configs optional"

## License

The go-themis library (i.e. all code outside of the `cmd` directory) is licensed under the
[GNU Lesser General Public License v3.0](https://www.gnu.org/licenses/lgpl-3.0.en.html), also
included in our repository in the `COPYING.LESSER` file.

The go-themis binaries (i.e. all code inside of the `cmd` directory) is licensed under the
[GNU General Public License v3.0](https://www.gnu.org/licenses/gpl-3.0.en.html), also included
in our repository in the `COPYING` file.