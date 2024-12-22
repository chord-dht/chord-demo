#!/bin/bash

# Check if NUM_CERTS is provided as an argument
if [ -z "$1" ]; then
  echo "Usage: $0 <NUM_CERTS>"
  exit 1
fi

# Number of private keys and CSRs to generate in batch
NUM_CERTS=$1

# Generate private keys and CSRs
for i in $(seq 1 $NUM_CERTS); do
  # Generate Peer private key
  if openssl genpkey -algorithm RSA -out peer_$i.key; then
    # Construct the subject string
    subj="/C=SE"
    subj="$subj/ST=Vastra_Gotaland"
    subj="$subj/L=Goteborg"
    subj="$subj/O=DS16"
    subj="$subj/OU=Chalmers"
    subj="$subj/CN=peer_$i.example.com"

    # Generate CSR
    if openssl req -new -key peer_$i.key -out peer_$i.csr -subj "$subj"; then
      echo "Private key and CSR generated: peer_$i.key and peer_$i.csr"
    else
      echo "Failed to generate CSR for peer_$i"
    fi
  else
    echo "Failed to generate private key for peer_$i"
  fi
done