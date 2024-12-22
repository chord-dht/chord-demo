#!/bin/bash

# Define colors for better output readability
GREEN="\033[0;32m"
YELLOW="\033[1;33m"
RED="\033[0;31m"
CYAN="\033[0;36m"
NC="\033[0m" # No Color

separator() {
  echo -e "${CYAN}----------------------------------------------------------${NC}"
}

MOVE_KEYS_AND_CERTS=false
N=0
CRT_KEY_DIR="crt_key"
PACK_DIR="pack"

# Remove the previous pack directory
separator
echo -e "${YELLOW}Cleaning up previous pack directory...${NC}"
rm -rf $PACK_DIR

if [ $? -eq 0 ]; then
  echo -e "${GREEN}Previous pack directory cleaned successfully!${NC}"
else
  echo -e "${RED}Failed to clean previous pack directory.${NC}"
  exit 1
fi

# Parse command line arguments
while getopts "mn:" opt; do
  case $opt in
    m)
      MOVE_KEYS_AND_CERTS=true
      ;;
    n)
      N=$OPTARG
      ;;
    \?)
      echo -e "${RED}Invalid option: -$OPTARG${NC}" >&2
      exit 1
      ;;
  esac
done

# Determine the number of peers
if [ $N -eq 0 ]; then
  if [ -d $CRT_KEY_DIR ]; then
    N=$(ls $CRT_KEY_DIR/*.crt 2>/dev/null | wc -l)
  fi

  if [ $N -eq 0 ]; then
    echo -e "${RED}No .crt files found in $CRT_KEY_DIR. Please check the directory.${NC}"
    exit 1
  fi
fi

separator
echo -e "${YELLOW}Building the chord executable...${NC}"
go build -o chord > /dev/null 2>&1

if [ $? -eq 0 ]; then
  echo -e "${GREEN}Chord executable built successfully!${NC}"
else
  echo -e "${RED}Failed to build the chord executable.${NC}"
  exit 1
fi

separator
echo -e "${YELLOW}Creating pack directory with $N peers...${NC}"

# Create the pack directory and peer subdirectories
for ((i=1; i<=N; i++)); do
  PEER_DIR="$PACK_DIR/peer_$i"

  mkdir -p $PEER_DIR

  if [ "$MOVE_KEYS_AND_CERTS" = true ]; then
    # Copy the cacert.pem file to the crt directory
    cp $CRT_KEY_DIR/cacert.pem $PEER_DIR/ > /dev/null 2>&1
    # Copy the corresponding peer_i.crt and peer_i.key files to the key directory
    cp $CRT_KEY_DIR/peer_$i.crt $PEER_DIR/peer.crt > /dev/null 2>&1
    cp $CRT_KEY_DIR/peer_$i.key $PEER_DIR/peer.key > /dev/null 2>&1
  fi

  # Copy the chord executable to the peer directory
  cp chord $PEER_DIR/ > /dev/null 2>&1

  # generate aes key to each peer directory
  openssl rand -hex 32 > $PEER_DIR/aes_key.txt

  echo -e "${GREEN}Peer $i directory created successfully.${NC}"
done

rm chord

separator
echo -e "${GREEN}Pack directory created successfully with $N peers!${NC}"