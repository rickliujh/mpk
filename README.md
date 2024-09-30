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
`mpk group create mygroup`
```
Usage:
  mpk group [flags]
  mpk group [command]

Available Commands:
  create      Creating peer groups
  list        List existing peer groups

```
Then, create peers under group.
`mpk peer create peer-a peer-b --group mygroup --threshold 1`
```
Usage:
  mpk group [flags]
  mpk group [command]

Available Commands:
  create      Creating peer groups
  list        List existing peer groups
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
mpk sign -g mygroup -f ./sig "hello world"
```

A signature file will be generated in current folder named `sig.json`.

Finally, verify the Signature.

```
mpk verify -g mygroup -f ./sig.json "hello world"
```

You could try verify by using a wrong message, and see if it actually verifies the signature.

```
mpk verify -g mygroup -f ./sig.json "hello world!"
```

You can always use -h for any command and sub command for available actions.

### Vault path

The keys are stored at `~./.confg/mpk`

> [!WARNING]
> The keys encryption is not supported at this point.

### Underlay
[bnb-chain/tss-lib Threshold Signature Scheme, for ECDSA and EDDSA](https://github.com/bnb-chain/tss-lib/tree/master)
