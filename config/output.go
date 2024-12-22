package config

import (
	"chord/log"
	"fmt"
)

func (cfg *Config) printJoinInfo() {
	log.Logger.Print(log.CenterTitle("Join Information", "-"))
	log.PrintKeyValue("Join Address", cfg.JoinAddress)
	log.PrintKeyValue("Join Port", cfg.JoinPort)
}

func (cfg *Config) PrintTimingInfo() {
	log.Logger.Print(log.CenterTitle("Timing Information", "-"))
	log.PrintKeyValue("Stabilize Time", fmt.Sprintf("%d ms", cfg.StabilizeTime))
	log.PrintKeyValue("Fix Fingers Time", fmt.Sprintf("%d ms", cfg.FixFingersTime))
	log.PrintKeyValue("Check Predecessor Time", fmt.Sprintf("%d ms", cfg.CheckPredecessorTime))
}

func (cfg *Config) printSuccessors() {
	log.Logger.Print(log.CenterTitle("Successors", "-"))
	log.PrintKeyValue("Successors", fmt.Sprintf("%d", cfg.Successors))
}

func (cfg *Config) printIdentifier() {
	log.Logger.Print(log.CenterTitle("Identifier", "-"))
	log.PrintKeyValue("Identifier", cfg.Identifier)
}

func (cfg *Config) printMode() {
	log.Logger.Print(log.CenterTitle("Mode", "-"))
	log.PrintKeyValue("Mode", cfg.Mode)
}

func (cfg *Config) printAES() {
	log.Logger.Print(log.CenterTitle("AES", "-"))
	log.PrintKeyValue("AES", cfg.AESBool)
}

func (cfg *Config) printTLS() {
	log.Logger.Print(log.CenterTitle("TLS", "-"))
	log.PrintKeyValue("TLS", cfg.TLSBool)
}

// Print the configuration to the console.
func (cfg *Config) Print() {
	log.Logger.Print(log.CenterTitle("Configuration", "="))
	defer log.Logger.Print(log.CenterTitle("Configuration", "="))

	log.PrintKeyValue("IP Address", cfg.IpAddress)
	log.PrintKeyValue("Port", cfg.Port)

	if cfg.Mode == "join" {
		cfg.printJoinInfo()
	}

	cfg.PrintTimingInfo()

	cfg.printSuccessors()

	if cfg.Identifier != Unspecified {
		cfg.printIdentifier()
	}

	cfg.printMode()

	cfg.printAES()

	cfg.printTLS()
}
