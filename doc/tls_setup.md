# TLS Setup and Distribution

## Generate CA Private Key and Root Certificate

To use a self-built CA (Certificate Authority) to establish a CA and batch issue certificates, you can follow these steps:

1. Install OpenSSL: Ensure you have installed the OpenSSL tool. Download and install it from the [OpenSSL official website](https://www.openssl.org/).

2. Generate CA Private Key and Root Certificate: Generate the CA's private key and self-signed root certificate.

```sh
# Generate CA private key, Enter PEM pass phrase, record it yourself
openssl genpkey -algorithm RSA -out ca.key -aes256

# Generate self-signed root certificate. First enter the password, then enter the subject fields
openssl req -x509 -new -nodes -key ca.key -sha256 -days 3650 -out ca.crt
```

When using a self-built CA to issue certificates, the subject of the CA's root certificate and the subject of the CSR must meet the following requirements:

- Must be the same:
  - Country (C): Usually needs to be consistent.
  - State/Province (ST): Usually needs to be consistent.
  - City/Locality (L): Usually needs to be consistent.
  - Organization (O): Usually needs to be consistent.
- Must be different:
  - Common Name (CN): Must be different. The CN of the CA's root certificate is usually the name of the CA, while the CN of the CSR is usually the name of the specific server or client.

3. Create CA Configuration File: Create an OpenSSL configuration file (e.g., `openssl.cnf`) to define the CA's configuration.

4. Initialize CA Directory Structure: Initialize the CA's directory structure and necessary files.

```sh
# Initialize CA Directory Structure
mkdir -p demoCA/private
mkdir -p demoCA/newcerts
touch demoCA/index.txt
echo 1000 > demoCA/serial
cp ca.crt demoCA/cacert.pem
cp ca.key demoCA/private/cakey.pem
rm ca.crt
rm ca.key
```

## Generate Peer Private Key and CSR

Generate private keys and certificate signing requests (CSR) for each Peer.

```sh
# Generate Peer private key
openssl genpkey -algorithm RSA -out peer.key

# Generate CSR
openssl req -new -key peer.key -out peer.csr -config openssl.cnf
```

## Issue Certificates

Use the CA to issue certificates for the Peers.

```sh
openssl ca -config openssl.cnf -in peer.csr -out peer.crt -batch
```
