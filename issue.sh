#!/bin/bash

# Define colors for better output readability
GREEN="\033[0;32m"
YELLOW="\033[1;33m"
RED="\033[0;31m"
CYAN="\033[0;36m"
NC="\033[0m" # No Color

# Print a separator line
separator() {
  echo -e "${CYAN}----------------------------------------------------------${NC}"
}

# Check if PASSWORD and NUM_CERTS are provided as arguments
if [ -z "$1" ] || [ -z "$2" ]; then
  echo -e "${RED}Usage: $0 <PASSWORD> <NUM_CERTS>${NC}"
  exit 1
fi

# CA private key password
PASSWORD=$1

# Number of certificates to issue in batch
NUM_CERTS=$2

separator
echo -e "${YELLOW}Step 1: Generating CA private key and root certificate...${NC}"
# Generate CA private key
openssl genpkey -algorithm RSA -out ca.key -aes256 -pass pass:$PASSWORD > /dev/null 2>&1

# Generate self-signed root certificate
openssl req -x509 -new -key ca.key -sha256 -days 3650 -out ca.crt \
  -passin pass:$PASSWORD \
  -subj "/C=SE/ST=Vastra_Gotaland/L=Goteborg/O=DS16/OU=Chalmers/CN=ca.example.com" > /dev/null 2>&1

if [ $? -eq 0 ]; then
  echo -e "${GREEN}CA private key and root certificate generated successfully!${NC}"
else
  echo -e "${RED}Failed to generate CA private key or root certificate.${NC}"
  exit 1
fi

separator
echo -e "${YELLOW}Step 2: Initializing CA directory structure...${NC}"
# Initialize CA Directory Structure
mkdir -p demoCA/private demoCA/newcerts
touch demoCA/index.txt
echo 1000 > demoCA/serial
cp ca.crt demoCA/cacert.pem
cp ca.key demoCA/private/cakey.pem
rm ca.crt ca.key

if [ $? -eq 0 ]; then
  echo -e "${GREEN}CA directory structure initialized successfully!${NC}"
else
  echo -e "${RED}Failed to initialize CA directory structure.${NC}"
  exit 1
fi

cd demoCA
chmod +x generate_keys_and_csrs.sh gen_cert.sh

separator
echo -e "${YELLOW}Step 3: Generating private keys and CSRs for peers...${NC}"
# Generate keys and csrs
./generate_keys_and_csrs.sh $NUM_CERTS > /dev/null 2>&1

if [ $? -eq 0 ]; then
  echo -e "${GREEN}Keys and CSRs generated successfully for $NUM_CERTS peers!${NC}"
else
  echo -e "${RED}Failed to generate keys and CSRs.${NC}"
  exit 1
fi

separator
echo -e "${YELLOW}Step 4: Issuing certificates for peers...${NC}"
# Generate (Issue) the cert for peers
./gen_cert.sh $PASSWORD $NUM_CERTS > /dev/null 2>&1

if [ $? -eq 0 ]; then
  echo -e "${GREEN}Certificates issued successfully for $NUM_CERTS peers!${NC}"
else
  echo -e "${RED}Failed to issue certificates.${NC}"
  exit 1
fi

separator
echo -e "${YELLOW}Step 5: Organizing output files...${NC}"
# Move the cacert.pem, key and crt out to `crt_key` dir, delete csr files
mkdir -p ../crt_key
cp cacert.pem ../crt_key/
mv peer_*.crt ../crt_key/
mv peer_*.key ../crt_key/
rm peer_*.csr

if [ $? -eq 0 ]; then
  echo -e "${GREEN}All files organized successfully in the 'crt_key' directory!${NC}"
else
  echo -e "${RED}Failed to organize output files.${NC}"
  exit 1
fi

cd ../
separator
echo -e "${GREEN}Script completed successfully!${NC}"