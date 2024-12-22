# Lab 3: Chord

## Doc

- [Task](doc/task.md)
- [Partial chord paper](doc/chord_paper.md)
- [successors and backups](doc/successors_backups.md)
- [TLS Principles](doc/tls.md)
- [TLS setup](doc/tls_setup.md)
- [Tricks and Issues](doc/tricks_issue.md)

## Base structure

![Base structure](doc/pic/basic_structure.png)

1. Node: Responsible for implementing the core functionality of the Chord protocol and including a CLI Interface. It enables distributed key lookups, node maintenance, and routing within the system. Through fileSystem interface, it can interact with the file system to perform operations such as reading and writing files.
2. NodeInfo: Contains the node’s unique identifier, IP address, and port. This information uniquely identifies a node and enables it to be called via RPC.
3. CLI Interface: Provides a command-line interface for users to interact with the node and perform operations like lookups or debugging.
4. FileSystem Interface: We define a StorageSystem which implements the FileSystem Interface, enabling caching and saving file operations on the local disk and memory.

## Build and run on local

```shell
chmod +x issue.sh
chmod +x create_pack.sh
chmod +x run_all.sh
```

```shell
./create_pack.sh -n 8
```

An example usage to start a new Chord ring is:

```shell
ADDRESS="127.0.0.1"
PORT="4171"

BIN_DIR="pack/peer_1"

cd $BIN_DIR
VERBOSE=1 ./chord -a $ADDRESS -p $PORT --ts 3000 --tff 1000 --tcp 3000 -r 4
```

An example usage to join an existing Chord ring is:

```shell
ADDRESS="127.0.0.1"
PORT="4170"

BIN_DIR="pack/peer_2"

cd $BIN_DIR
VERBOSE=1 ./chord -a $ADDRESS -p $PORT --ja 127.0.0.1 --jp 4171 --ts 3000 --tff 1000 --tcp 3000 -r 4
```

```shell
ADDRESS="127.0.0.1"
PORT="4172"

BIN_DIR="pack/peer_3"

cd $BIN_DIR
VERBOSE=1 ./chord -a $ADDRESS -p $PORT --ja 127.0.0.1 --jp 4171 --ts 3000 --tff 1000 --tcp 3000 -r 4
```

## test on local

```shell
cd test
go build -o chord ../
rm -rf node_*/
VERBOSE=1 python main.py
```

## Pack binary file, crts and key (bonus)

Add more flags:

- `-aes` and `-aeskey`
- `-tls`, `-cacert`, `-servercert` and `-serverkey`

```shell
NUM_CLIENTS=8
PASSWORD="XXXX"

# After issuing, the crts and keys will be located in crt_key/ directory.
./issue.sh $PASSWORD $NUM_CLIENTS

# Pack all file into corresponding directory.
./create_pack.sh -m -n $NUM_CLIENTS
```

```shell
ADDRESS="127.0.0.1"
PORT="4170"

BIN_DIR="pack/peer_1"

cd $BIN_DIR
VERBOSE=1 ./chord -a $ADDRESS -p $PORT \
        --ts 3000 --tff 1000 --tcp 3000 \
        -r 4 \
        -aes -aeskey "aes_key.txt" \
        -tls -cacert "cacert.pem" -servercert "peer.crt" -serverkey "peer.key"
```

An example usage to join an existing Chord ring is:

```shell
ADDRESS="127.0.0.1"
PORT="4171"

BIN_DIR="pack/peer_2"

cd $BIN_DIR
VERBOSE=1 ./chord -a $ADDRESS -p $PORT \
        --ja 127.0.0.1 --jp 4170 \
        --ts 3000 --tff 1000 --tcp 3000 \
        -r 4 \
        -aes -aeskey "aes_key.txt" \
        -tls -cacert "cacert.pem" -servercert "peer.crt" -serverkey "peer.key"
```

```shell
ADDRESS="127.0.0.1"
PORT="4172"

BIN_DIR="pack/peer_3"

cd $BIN_DIR
VERBOSE=1 ./chord -a $ADDRESS -p $PORT \
        --ja 127.0.0.1 --jp 4170 \
        --ts 3000 --tff 1000 --tcp 3000 \
        -r 4 \
        -aes -aeskey "aes_key.txt" \
        -tls -cacert "cacert.pem" -servercert "peer.crt" -serverkey "peer.key"
```

## Run 8 instances on local using tmux

```shell
./run_all.sh
```
