package config

import (
	"chord/aes"
	"crypto/tls"
	"flag"
	"fmt"
	"net"
	"os"
	"regexp"
	"strconv"
)

const Unspecified = "Unspecified"

type Config struct {
	IpAddress            string
	Port                 string
	JoinAddress          string
	JoinPort             string
	StabilizeTime        int
	FixFingersTime       int
	CheckPredecessorTime int
	Successors           int
	Identifier           string

	Mode string // "create" or "join"

	AESBool    bool   // turn on/off AES encryption when storing files
	AESKeyPath string // path to the AES key file
	AESKey     []byte // AES key

	TLSBool         bool // turn on/off TLS connection
	CaCert          string
	ServerCert      string
	ServerKey       string
	ServerTLSConfig *tls.Config
	ClientTLSConfig *tls.Config
}

var NodeConfig *Config

// Parse the command line arguments.
func parseFlags() *Config {
	cfg := &Config{}

	flag.StringVar(&cfg.IpAddress, "a", Unspecified, "The IP address that the Chord client will bind to, as well as advertise to other nodes. Must be specified.")
	flag.StringVar(&cfg.Port, "p", Unspecified, "The port that the Chord client will bind to and listen on. Must be specified.")
	flag.StringVar(&cfg.JoinAddress, "ja", Unspecified, "The IP address of the machine running a Chord node. Must be specified if --jp is specified.")
	flag.StringVar(&cfg.JoinPort, "jp", Unspecified, "The port that an existing Chord node is bound to and listening on. Must be specified if --ja is specified.")
	flag.IntVar(&cfg.StabilizeTime, "ts", 0, "The time in milliseconds between invocations of 'stabilize'. Must be specified, with a value in the range of [1,60000].")
	flag.IntVar(&cfg.FixFingersTime, "tff", 0, "The time in milliseconds between invocations of 'fix fingers'. Must be specified, with a value in the range of [1,60000].")
	flag.IntVar(&cfg.CheckPredecessorTime, "tcp", 0, "The time in milliseconds between invocations of 'check predecessor'. Must be specified, with a value in the range of [1,60000].")
	flag.IntVar(&cfg.Successors, "r", 0, "The number of successors maintained by the Chord client. Must be specified, with a value in the range of [1,32].")
	flag.StringVar(&cfg.Identifier, "i", Unspecified, "The Identifier (ID) assigned to the Chord client which will override the ID computed by the SHA1 sum of the client's IP address and port number. Represented as a string of 40 characters matching [0-9a-fA-F]. Optional parameter.")
	flag.BoolVar(&cfg.AESBool, "aes", false, "Enable AES encryption. Optional parameter.")
	flag.StringVar(&cfg.AESKeyPath, "aeskey", "", "The path to the AES key file. Must be specified if --aes is specified.")
	flag.BoolVar(&cfg.TLSBool, "tls", false, "Enable TLS connection. Optional parameter.")
	flag.StringVar(&cfg.CaCert, "cacert", "", "The path to the CA certificate file. Must be specified if --tls is specified.")
	flag.StringVar(&cfg.ServerCert, "servercert", "", "The path to the server certificate file. Must be specified if --tls is specified.")
	flag.StringVar(&cfg.ServerKey, "serverkey", "", "The path to the server key file. Must be specified if --tls is specified.")

	flag.Parse()

	return cfg
}

func ReadConfig() *Config {
	cfg := parseFlags()
	if err := validateConfig(cfg); err != nil {
		fmt.Println("Failed to validate config:", err)
		os.Exit(1)
	}

	determineMode(cfg)

	if err := determineAES(cfg); err != nil {
		fmt.Println("Failed to determine AES:", err)
		os.Exit(1)
	}

	if err := determineTLS(cfg); err != nil {
		fmt.Println("Failed to determine TLS:", err)
		os.Exit(1)
	}
	return cfg
}

// Validate the configuration, returning an error if the configuration is invalid.
func validateConfig(cfg *Config) error {
	if cfg.IpAddress == Unspecified {
		return fmt.Errorf("ip address must be specified")
	}

	if net.ParseIP(cfg.IpAddress) == nil {
		return fmt.Errorf("invalid ip address format")
	}

	port, err := strconv.Atoi(cfg.Port)
	if err != nil || port <= 1024 || port > 65535 {
		return fmt.Errorf("port must be in the range of (1024,65535]")
	}

	if cfg.StabilizeTime < 1 || cfg.StabilizeTime > 60000 {
		return fmt.Errorf("stabilize time must be in the range of [1,60000] milliseconds")
	}

	if cfg.FixFingersTime < 1 || cfg.FixFingersTime > 60000 {
		return fmt.Errorf("fix fingers time must be in the range of [1,60000] milliseconds")
	}

	if cfg.CheckPredecessorTime < 1 || cfg.CheckPredecessorTime > 60000 {
		return fmt.Errorf("check predecessor time must be in the range of [1,60000] milliseconds")
	}

	if cfg.Successors < 1 || cfg.Successors > 32 {
		return fmt.Errorf("number of successors must be in the range of [1,32]")
	}

	if (cfg.JoinAddress != Unspecified && cfg.JoinPort == Unspecified) || (cfg.JoinAddress == Unspecified && cfg.JoinPort != Unspecified) {
		return fmt.Errorf("both --ja and --jp must be specified together")
	}

	if cfg.JoinAddress != Unspecified && cfg.JoinPort != Unspecified {
		if net.ParseIP(cfg.JoinAddress) == nil {
			return fmt.Errorf("invalid join address format")
		}

		joinPort, err := strconv.Atoi(cfg.JoinPort)
		if err != nil || joinPort <= 1024 || joinPort > 65535 {
			return fmt.Errorf("join port must be in the range of (1024,65535]")
		}
	}

	if cfg.Identifier != Unspecified {
		matched, err := regexp.MatchString("^[0-9a-fA-F]{40}$", cfg.Identifier)
		if err != nil || !matched {
			return fmt.Errorf("invalid Identifier format")
		}
	}

	if cfg.AESBool {
		if cfg.AESKeyPath == "" {
			return fmt.Errorf("AES key path must be specified if --aes is specified")
		}
		if _, err := os.Stat(cfg.AESKeyPath); os.IsNotExist(err) {
			return fmt.Errorf("AES key file does not exist at specified path")
		}
	}

	if cfg.TLSBool {
		if cfg.CaCert == "" {
			return fmt.Errorf("CA certificate path must be specified if --tls is specified")
		} else if _, err := os.Stat(cfg.CaCert); os.IsNotExist(err) {
			return fmt.Errorf("CA certificate file does not exist at specified path")
		}
		if cfg.ServerCert == "" {
			return fmt.Errorf("server certificate path must be specified if --tls is specified")
		} else if _, err := os.Stat(cfg.ServerCert); os.IsNotExist(err) {
			return fmt.Errorf("server certificate file does not exist at specified path")
		}
		if cfg.ServerKey == "" {
			return fmt.Errorf("server key path must be specified if --tls is specified")
		} else if _, err := os.Stat(cfg.ServerKey); os.IsNotExist(err) {
			return fmt.Errorf("server key file does not exist at specified path")
		}
	}

	return nil
}

// Determine the mode of the Chord client.
// If the join address and join port are both specified, the mode is "join".
// Otherwise, the mode is "create".
func determineMode(cfg *Config) {
	if cfg.JoinAddress != Unspecified && cfg.JoinPort != Unspecified {
		cfg.Mode = "join"
	} else {
		cfg.Mode = "create"
	}
}

func determineAES(cfg *Config) error {
	var err error = nil
	if cfg.AESBool {
		cfg.AESKey, err = aes.LoadKey(cfg.AESKeyPath)
		return err
	}
	return nil
}

func determineTLS(cfg *Config) error {
	if cfg.TLSBool {
		var err error = nil
		cfg.ServerTLSConfig, cfg.ClientTLSConfig, err = SetupTLS(cfg.CaCert, cfg.ServerCert, cfg.ServerKey)
		return err
	}
	return nil
}
