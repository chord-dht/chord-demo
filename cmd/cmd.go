package cmd

import (
	"bufio"
	"chord/node"
	"chord/tools"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	PRINTSTATE = "PRINTSTATE"
	LOOKUP     = "LOOKUP"
	STOREFILE  = "STOREFILE"
	STOREFILES = "STOREFILES"
	GETFILE    = "GETFILE"
	QUIT       = "QUIT"
	CLEAR      = "CLEAR"
)

// DownloadDir download directory
const DownloadDir = "download"

const UserInputSeparatorLine = "----------------------------------"

const DirPermission = 0755

const FilePermission = 0644

// LoopProcessUserCommand : Process the user input command.
func LoopProcessUserCommand(chordNode *node.Node) {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Enter command: ")
		if !scanner.Scan() { // Scan() returns false if an error occurs or EOF is reached
			break
		}
		// but will return true if a token (a line) is scanned
		command := scanner.Text()
		command = strings.TrimSpace(command)
		command = strings.ToUpper(command)
		handleUserInput(command, chordNode, scanner)
	}
	if err := scanner.Err(); err != nil {
		_, err := fmt.Fprintln(os.Stderr, "reading standard input:", err)
		if err != nil {
			return
		}
	}
}

func handleUserInput(command string, chordNode *node.Node, scanner *bufio.Scanner) {
	switch command {
	case PRINTSTATE:
		handlePrintState(chordNode)
	case LOOKUP:
		handleLookup(chordNode, scanner)
	case STOREFILE:
		handleStoreFile(chordNode, scanner)
	case STOREFILES:
		handleStoreFiles(chordNode, scanner)
	case GETFILE:
		handleGetFile(chordNode, scanner)
	case QUIT:
		handleQuit(chordNode)
	case CLEAR:
		handleClear()
	default:
		handleInvalidCommand()
	}

	if err := scanner.Err(); err != nil {
		_, err := fmt.Fprintln(os.Stderr, "reading standard input:", err)
		if err != nil {
			return
		}
	}
}

func handlePrintState(chordNode *node.Node) {
	fmt.Println(UserInputSeparatorLine)
	fmt.Printf("Command: %s\n", PRINTSTATE)
	CmdPrintState(chordNode)
	fmt.Println(UserInputSeparatorLine)
}

func handleLookup(chordNode *node.Node, scanner *bufio.Scanner) {
	fmt.Print("Enter key to lookup: ")
	if scanner.Scan() {
		filename := scanner.Text()
		fmt.Println(UserInputSeparatorLine)
		fmt.Printf("Command: %s %s\n", LOOKUP, filename)

		identifier := tools.GenerateIdentifier(filename)
		fmt.Printf("Identifier of %s: %s\n", filename, identifier)

		successor, err := CmdLookUp(chordNode.GetInfo(), filename)
		if err != nil {
			fmt.Printf("Lookup %s failed: %v\n", filename, err)
		} else {
			fmt.Printf("Lookup %s success: ", filename)
			successor.PrintInfo()
		}

		fmt.Println(UserInputSeparatorLine)
	}
}

func handleStoreFile(chordNode *node.Node, scanner *bufio.Scanner) {
	fmt.Print("Enter the file location: ")
	if scanner.Scan() {
		location := scanner.Text()
		fmt.Println(UserInputSeparatorLine)
		fmt.Printf("Command: %s %s\n", STOREFILE, location)

		targetNode, err := CmdStoreFile(chordNode.GetInfo(), location)
		if err != nil {
			fmt.Printf("Storing file %s failed: %v\n", location, err)
		} else {
			fmt.Printf("Storing file %s success, target node: ", location)
			targetNode.PrintInfo()
		}

		fmt.Println(UserInputSeparatorLine)
	}
}

func handleStoreFiles(chordNode *node.Node, scanner *bufio.Scanner) {
	fmt.Print("Enter the directory location: ")
	if scanner.Scan() {
		dirLocation := scanner.Text()
		fmt.Println(UserInputSeparatorLine)
		fmt.Printf("Command: %s %s\n", STOREFILES, dirLocation)

		err := getAndStoreFilesInDirectory(dirLocation, chordNode)
		if err != nil {
			fmt.Printf("Storing files in directory %s failed: %v\n", dirLocation, err)
		}

		fmt.Println(UserInputSeparatorLine)
	}
}

func handleGetFile(chordNode *node.Node, scanner *bufio.Scanner) {
	fmt.Print("Enter the file name: ")
	if scanner.Scan() {
		filename := scanner.Text()
		fmt.Println(UserInputSeparatorLine)
		fmt.Printf("Command: %s %s\n", GETFILE, filename)

		targetNode, fileContent, err := CmdGetFile(chordNode.GetInfo(), filename)
		if err != nil {
			fmt.Printf("Getting file %s failed: %v\n", filename, err)
		} else {
			fmt.Printf("Successfully Getting file %s from node: ", filename)
			targetNode.PrintInfo()           // print the node info that stores the file
			PrintFirstNLines(fileContent, 3) // print the first 3 lines of the file
			filePath := filepath.Join(DownloadDir, filename)
			if err := SaveFile(filePath, fileContent); err != nil {
				fmt.Printf("Can't save file %s\n", filename)
			} else {
				fmt.Printf("Successfully save file %s\n", filename)
			}
		}
		fmt.Println(UserInputSeparatorLine)
	}
}

func handleQuit(chordNode *node.Node) {
	fmt.Println(UserInputSeparatorLine)
	fmt.Printf("Command: %s\n", QUIT)
	CmdQuit(chordNode)
	fmt.Println(UserInputSeparatorLine)
	os.Exit(0)
}

func handleClear() {
	fmt.Print("\033[H\033[2J") // clear screen
}

func handleInvalidCommand() {
	fmt.Println(UserInputSeparatorLine)
	fmt.Println("Invalid command!")
	fmt.Println(UserInputSeparatorLine)
}

/*                             Helper function                             */

func getAndStoreFilesInDirectory(dirLocation string, chordNode *node.Node) error {
	return filepath.Walk(dirLocation, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			targetNode, err := CmdStoreFile(chordNode.GetInfo(), path)
			if err != nil {
				fmt.Printf("Storing file %s failed: %v\n", path, err)
			} else {
				fmt.Printf("Storing file %s success, target node: ", path)
				targetNode.PrintInfo()
			}
		}
		return nil
	})
}

// PrintFirstNLines prints the first N lines from a byte slice.
func PrintFirstNLines(fileContent []byte, n int) {
	lines := strings.SplitN(string(fileContent), "\n", n+1)
	for i, line := range lines {
		if i < n {
			fmt.Println(line)
		}
	}
}

// SaveFile Save the file to the given filePath directory
func SaveFile(filePath string, fileContent []byte) error {
	// Extract the directory path from the file path
	dirPath := filepath.Dir(filePath)

	// Check if the directory exists, if not, create it
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		err = os.MkdirAll(dirPath, DirPermission)
		if err != nil {
			return err
		}
	}

	err := os.WriteFile(filePath, fileContent, FilePermission)
	return err
}

/*                             Helper function                             */
