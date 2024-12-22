package cmd

import (
	"chord/aes"
	"chord/config"
	"chord/node"
	"chord/tools"
	"fmt"
	"os"
	"path/filepath"
)

/*                             Directly operating on local node                             */

func CmdPrintState(chordNode *node.Node) {
	chordNode.PrintState()
}

func CmdQuit(chordNode *node.Node) {
	chordNode.Quit()
}

/*                             Directly operating on local node                             */

/*                             Operating through Node address (nodeInfo)                            */

// the start point of the function is the local node
// but throught the start node, we can find the target node
// then we directly communicate with the target node!

// lookup the successor node of the key in the chord ring
func CmdLookUp(startNode *node.NodeInfo, filename string) (*node.NodeInfo, error) {
	// step 1: generate the identifier of the filename
	identifier := tools.GenerateIdentifier(filename)
	fmt.Println("The identifier of the filename is", identifier)
	// step 2: find the successor node of the (filename) identifier
	targetNode, err := startNode.FindSuccessorIter(identifier)
	return targetNode, err
}

// store the file in the chord ring
func CmdStoreFile(startNode *node.NodeInfo, location string) (*node.NodeInfo, error) {
	// Step 1: Validate and normalize the file path
	absPath, err := filepath.Abs(location)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path: %v", err)
	}

	// Step 2: Extract the file name from the path
	filename := filepath.Base(absPath)

	// Step 3: Perform a "LookUp" to findSuccessorIter the correct node to store the file
	targetNode, err := CmdLookUp(startNode, filename)
	if err != nil {
		return nil, fmt.Errorf("failed to lookup the target node: %v", err)
	}

	// Step 4: Read the file content
	fileContent, err := os.ReadFile(absPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read the file content: %v", err)
	}

	// Step 5: Encrypt the file content if AESBool is true
	if config.NodeConfig.AESBool {
		fileContent, err = aes.EncryptAES(fileContent, config.NodeConfig.AESKey)
		if err != nil {
			return nil, fmt.Errorf("failed to encrypt the file content: %v", err)
		}
	}

	// Step 6: Store the file content in the target node's storage
	reply, err := targetNode.StoreFile(filename, fileContent)
	if err != nil {
		return nil, fmt.Errorf("failed to get the reply from node %s: %v", targetNode.Identifier.String(), err)
	}
	if !reply.Success {
		return nil, fmt.Errorf("node %s reply: it can't store the file: %v", targetNode.Identifier.String(), err)
	}
	return targetNode, nil
}

// get the file content from the chord ring, also return the target node information
func CmdGetFile(startNode *node.NodeInfo, filename string) (*node.NodeInfo, []byte, error) {
	// step 1: find the successor node (targetNode) of the key (filename)
	targetNode, err := CmdLookUp(startNode, filename)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to lookup the target node: %v", err)
	}

	// step 2: get the GetFile reply from the target node (successful flag and file content)
	reply, err := targetNode.GetFile(filename)
	// if error occurs, it means RPC call failed
	if err != nil {
		return targetNode, nil, fmt.Errorf("failed to get the reply from node %s: %v", targetNode.Identifier.String(), err)
	}
	// if successful flag is false, it means the file doesn't exist on the target node
	if !reply.Success {
		return targetNode, nil, fmt.Errorf("node %s reply: it doesn't have the file", targetNode.Identifier.String())
	}
	// Now we have the file content!
	fileContent := reply.FileContent

	// step 4: Decrypt the file content if AESBool is true
	if config.NodeConfig.AESBool {
		fileContent, err = aes.DecryptAES(fileContent, config.NodeConfig.AESKey)
		if err != nil {
			return targetNode, nil, fmt.Errorf("failed to decrypt the file content: %v", err)
		}

		// step 4.5: check the decrypted file content's entropy
		fileEntropy := aes.CalculateEntropy(fileContent)
		if fileEntropy > aes.FileEntropyThreshold {
			fmt.Printf(
				"The entropy of the file content is %f > FileEntropyThreshold %f, "+
					"you may not have the right key to decrypt the file.\n",
				fileEntropy, aes.FileEntropyThreshold,
			)
		} else {
			fmt.Printf("The entropy of the file content is %f < FileEntropyThreshold %f, "+
				"you may have the right key to decrypt the file.\n",
				fileEntropy, aes.FileEntropyThreshold,
			)
		}
	}
	return targetNode, fileContent, nil
}

/*                             Operating through Node address (nodeInfo)                            */
