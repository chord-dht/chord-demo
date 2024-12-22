package node

import (
	"chord/log"
	"fmt"
	"os"
	"time"
)

// Initialize begins the node, create or join
func (node *Node) Initialize(mode, joinAddress, joinPort string) {
	defer log.LogFunction()()

	switch mode {
	case "create":
		node.create()
		fmt.Println("Create Chord Ring success")
	case "join":
		node.joinRing(joinAddress, joinPort)
		fmt.Println("Join Chord Ring success")
	}

	// register it in rpc and start the server
	node.startServer()

	// start the periodic tasks
	node.StartPeriodicTasks()
}

// Create a new ring.
func (node *Node) create() {
	defer log.LogFunction()()

	// predecessor = nil
	// successor = node itself

	node.SetFirstSuccessor(&node.info)
	log.Info("node.Successors[0]: %v", node.info)
}

func (node *Node) joinRing(joinAddress, joinPort string) {
	// get full Info of join node
	joinNode := NewNodeInfoWithAddress(joinAddress, joinPort)
	joinNode, err := joinNode.GetNodeInfo()
	if err != nil {
		log.Error("Try to get join node Info failed, error: %v", err)
		fmt.Printf("Try to get join node Info failed, error: %v\n", err)
		os.Exit(1)
	}

	// They should have the same IdentifierLength and SuccessorsLength
	// Otherwise, the join operation will fail
	reply, err := joinNode.GetLength()
	if err != nil {
		log.Error("Try to get join node length failed, error: %v", err)
		fmt.Printf("Try to get join node length failed, error: %v\n", err)
		os.Exit(1)
	}
	if reply.IdentifierLength != node.identifierLength || reply.SuccessorsLength != node.successorsLength {
		log.Error("The join node has different IdentifierLength or SuccessorsLength")
		fmt.Printf("The join node has different IdentifierLength or SuccessorsLength\n")
		os.Exit(1)
	}

	// join the chord ring
	if err := node.join(joinNode); err != nil {
		log.Error("Join Chord Ring failed, error: %v", err)
		fmt.Printf("Join Chord Ring failed, error: %v\n", err)
		os.Exit(1)
	}
}

// Join an existing Chord ring containing node n' (joinNode).
func (node *Node) join(joinNode *NodeInfo) error {
	defer log.LogFunction()()

	log.Info("%v.join(%v)", node.info, joinNode)

	// predecessor = nil
	// successor = n'.find_successor(n)
	nodeInfo, err := joinNode.FindSuccessorIter(node.info.Identifier)
	if err != nil {
		log.Info("%v.find_successor(%v) failed, error: %v", joinNode, node.info, err)
		return fmt.Errorf("%v.find_successor(%v) failed, error: %v", joinNode, node.info, err)
	}
	if err := nodeInfo.LiveCheck(); err != nil {
		log.Info("%v.find_successor(%v) has bad result: %v", joinNode, node.info, err)
		return fmt.Errorf("%v.find_successor(%v) has bad result: %v", joinNode, node.info, err)
	}

	node.SetFirstSuccessor(nodeInfo)
	log.Info("Successfully join! Its successor is %v", nodeInfo)
	return nil
}

func (node *Node) StartPeriodicTasks() {
	go node.periodicStabilize(node.stabilizeTime)
	go node.periodicFixFingers(node.fixFingersTime)
	go node.periodicCheckPredecessor(node.checkPredecessorTime)

	fmt.Println("Waiting for periodic tasks to stabilize...")
	// Sleep for a duration to allow periodic tasks to stabilize
	time.Sleep(5 * time.Second) // Adjust the duration as needed
}

func (node *Node) periodicStabilize(stabilizeTime time.Duration) {
	ticker := time.NewTicker(stabilizeTime * time.Millisecond)
	for {
		select {
		case <-ticker.C:
			node.stabilize()
		case <-node.shutdownCh:
			ticker.Stop()
			return
		}
	}
}

func (node *Node) periodicFixFingers(fixFingersTime time.Duration) {
	ticker := time.NewTicker(fixFingersTime * time.Millisecond)
	for {
		select {
		case <-ticker.C:
			node.fixFingers()
		case <-node.shutdownCh:
			ticker.Stop()
			return
		}
	}
}

func (node *Node) periodicCheckPredecessor(checkPredecessorTime time.Duration) {
	ticker := time.NewTicker(checkPredecessorTime * time.Millisecond)
	for {
		select {
		case <-ticker.C:
			node.checkPredecessor()
		case <-node.shutdownCh:
			ticker.Stop()
			return
		}
	}
}
