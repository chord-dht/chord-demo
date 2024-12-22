#!/bin/bash

# Check if CA_PASS and NUM_CERTS are provided as arguments
if [ -z "$1" ] || [ -z "$2" ]; then
  echo "Usage: $0 <CA_PASS> <NUM_CERTS>"
  exit 1
fi

# CA private key password
CA_PASS=$1

# Number of certificates to issue in batch
NUM_CERTS=$2

# CA configuration file path
CA_CONFIG="openssl.cnf"

# CA private key and certificate path
CA_KEY="private/cakey.pem"
CA_CERT="cacert.pem"

# Clear CA database file
> index.txt

# Issue certificates
for i in $(seq 1 $NUM_CERTS); do
  # Read current serial number
  if [ -f serial ]; then
    SERIAL=$(cat serial)
  else
    SERIAL="01"
  fi

  # Ensure the serial number is hexadecimal
  if ! [[ "$SERIAL" =~ ^[0-9A-Fa-f]+$ ]]; then
    echo $SERIAL
    echo "Error: Serial number contains non-hexadecimal characters"
    exit 1
  fi

  # Generate new hexadecimal serial number
  NEW_SERIAL=$(printf "%X" $((0x$SERIAL + 1)))
  echo $NEW_SERIAL > serial

  # Generate certificate
  if openssl ca -config $CA_CONFIG -in peer_$i.csr -out peer_$i.crt -batch -passin pass:$CA_PASS; then
    echo "Certificate issued: peer_$i.crt"
  else
    echo "Failed to issue certificate for peer_$i"
  fi
done