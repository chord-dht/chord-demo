package main

import (
	cfs "chord/cachefilesystem"
	"chord/cmd"
	"chord/config"
	"chord/node"
	st "chord/storage"
	"chord/tools"
	"fmt"
	"math/big"
	"time"
)

func main() {
	// stage 1: parse the command line arguments, validate them and determine the mode and tls settings
	config.NodeConfig = config.ReadConfig()
	config.NodeConfig.Print()

	// stage 2: create a new chordNode
	chordNode, err := NewNodeWithConfig(config.NodeConfig, cfs.CacheStorageFactory)
	if err != nil {
		panic(err)
	}

	// stage 3: initialize the node
	// including create or join the ring
	// and start the server
	// and start the periodic tasks
	chordNode.Initialize(config.NodeConfig.Mode, config.NodeConfig.JoinAddress, config.NodeConfig.JoinPort)

	// stage 4: start the command line interface
	// read commands from stdin and execute them
	cmd.LoopProcessUserCommand(chordNode)
}

// NewNodeWithConfig uses the configuration to create a new node.
func NewNodeWithConfig(
	cfg *config.Config,
	storageFactory func(string) (st.Storage, error),
) (*node.Node, error) {
	// first set the IdentifierLength, you have to set it first
	identifierLength := 10 // identifier length (m)

	// then the path of the storage
	storageDir := "storage" // storage directory
	backupDir := "backup"   // backup directory

	// set up the identifier
	var identifier *big.Int
	if cfg.Identifier == config.Unspecified {
		networkAddress := cfg.IpAddress + ":" + cfg.Port
		identifier = tools.GenerateHash(networkAddress)
	} else {
		var err error
		identifier, err = tools.HexStringToBigInt(cfg.Identifier) // already checked before
		if err != nil {
			return nil, fmt.Errorf("failed to convert string to *big.Int: %w", err)
		}
	}
	identifier = identifier.And(identifier, tools.TwoMMinusOne) // need to mod 2^m anyway

	chordNode, err := node.NewNode(
		identifierLength,
		cfg.Successors,
		cfg.IpAddress,
		cfg.Port,
		identifier,
		storageFactory,
		storageDir,
		backupDir,
		time.Duration(cfg.StabilizeTime),
		time.Duration(cfg.FixFingersTime),
		time.Duration(cfg.CheckPredecessorTime),
		cfg.TLSBool,
		cfg.ServerTLSConfig,
		cfg.ClientTLSConfig,
	)
	if err != nil {
		return nil, fmt.Errorf("error creating node: %w", err)
	}

	return chordNode, nil
}
