package config

import (
	"chord/log"
	"crypto/tls"
	"crypto/x509"
	"os"
)

func SetupTLS(caCrt, serverCrt, serverKey string) (serverTLSConfig, clientTLSConfig *tls.Config, err error) {
	cert, err := tls.LoadX509KeyPair(serverCrt, serverKey)
	if err != nil {
		log.Error("server: loadkeys: %v", err)
		return nil, nil, err
	}
	serverTLSConfig = &tls.Config{
		Certificates:       []tls.Certificate{cert},
		InsecureSkipVerify: true, // just for testing
	}

	caCert, err := os.ReadFile(caCrt)
	if err != nil {
		log.Error("client: read ca cert: %v", err)
		return nil, nil, err
	}
	certPool := x509.NewCertPool()
	certPool.AppendCertsFromPEM(caCert)
	clientTLSConfig = &tls.Config{
		RootCAs:            certPool,
		InsecureSkipVerify: true, // just for testing
	}

	return serverTLSConfig, clientTLSConfig, nil
}
