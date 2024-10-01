# mpk

mpk is a multi party threshold signature CLI.
You can use it for parity management, keygen, signing, and signature verification.

## Installation

`go install github.com/rickliujh/mpk`

Or build from source code.

```shell
git clone github.com/rickliujh/mpk

cd mpk

make

cd ./build

./mpk help
```

## Usage

First we need create a peer group to manage the peers.

```
mpk group create mygroup
```

Then, create peers under group.

```
mpk peer create peer-a peer-b --group mygroup --threshold 1
```

You can also list the peers by.

```
mpk peer list
```

Now, use keygen to generate the keys for peers, and set the timeout for 1 minute

```
mpk keygen -g mygroup --timeout 1
```

And, we can now signing the message by these keys.

(Observe that threshold signing on a single machine defeats the purpose of threshold signing, but it only is a demonstration of the CLI)

```
mpk sign -g mygroup -f ./sig.json "hello world"
```

A signature file will be generated in current folder named `sig.json` specified by `-f` flag.

Finally, verify the Signature.

```
mpk verify -g mygroup -f ./sig.json "hello world"
```

You could try verify by using a wrong message, and see if it actually verifies the signature.

```
mpk verify -g mygroup -f ./sig.json "hello world!"
```

You can always use -h for any command and sub command for available actions.

```
Available Commands:
  completion  Generate the autocompletion script for the specified shell
  group       A brief description of your command
  help        Help about any command
  keygen      Generating the private keys for peer group
  peer        Peer is the party participating threshold signing
  sign        Signing a message as peer group
  verify      Verifying the signature

Flags:
  -h, --help     help for mpk
  -t, --toggle   Help message for toggle
```

### Vault path

Keys are stored at `~/.confg/mpk`

> [!WARNING]
> Encryption of keys is not supported at this point.

### Underlay
[bnb-chain/tss-lib Threshold Signature Scheme, for ECDSA and EDDSA](https://github.com/bnb-chain/tss-lib/tree/master)
