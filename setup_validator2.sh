#!/bin/bash

# Validator 2 Setup Script
# This script sets up a second validator node after being whitelisted via council governance
# See VALIDATOR_NODE_SETUP.md for detailed explanation

set -e  # Exit on error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Variables
export VALIDATOR_NAME="validator2"
export VALIDATOR_KEY="council-member-2"
export NODE_HOME="$HOME/.veranatest-validator2"
export CHAIN_ID="vna-testnet-1"
export PRIMARY_RPC="http://localhost:26657"
export NEW_NODE_RPC="http://localhost:26667"

echo -e "${GREEN}=== Validator 2 Setup Script ===${NC}"
echo "Validator Name: $VALIDATOR_NAME"
echo "Validator Key: $VALIDATOR_KEY"
echo "Node Home: $NODE_HOME"
echo "Chain ID: $CHAIN_ID"
echo ""

# Function to check if command succeeded
check_status() {
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✓ $1${NC}"
    else
        echo -e "${RED}✗ $1 failed${NC}"
        exit 1
    fi
}

# Step 1: Verify whitelisting
echo -e "${YELLOW}=== Step 1: Verify Validator is Whitelisted ===${NC}"
OPERATOR_ADDR=$(veranatestd keys show $VALIDATOR_KEY --bech val --keyring-backend test -a 2>/dev/null || echo "")
if [ -z "$OPERATOR_ADDR" ]; then
    echo -e "${RED}✗ Key $VALIDATOR_KEY not found in keyring${NC}"
    echo "Please create the key first or check the key name"
    exit 1
fi
echo "Operator Address: $OPERATOR_ADDR"

# Check if whitelisted
WHITELIST_CHECK=$(veranatestd query validatorregistry list-validator --node $PRIMARY_RPC 2>&1 | grep "$OPERATOR_ADDR" || echo "")
if [ -z "$WHITELIST_CHECK" ]; then
    echo -e "${RED}✗ Validator $OPERATOR_ADDR is not whitelisted${NC}"
    echo "Please whitelist via council governance first"
    echo "See COUNCIL_GOVERNANCE.md for instructions"
    exit 1
fi
echo -e "${GREEN}✓ Validator is whitelisted${NC}"
echo ""

# Step 2: Check if node already exists
if [ -d "$NODE_HOME" ]; then
    echo -e "${YELLOW}Node directory $NODE_HOME already exists${NC}"
    read -p "Do you want to remove it and start fresh? (yes/no): " REMOVE_CHOICE
    if [ "$REMOVE_CHOICE" = "yes" ]; then
        echo "Removing existing node directory..."
        rm -rf "$NODE_HOME"
        check_status "Removed existing directory"
    else
        echo "Please backup or remove $NODE_HOME manually"
        exit 1
    fi
fi
echo ""

# Step 3: Initialize node
echo -e "${YELLOW}=== Step 2: Initialize Node ===${NC}"
veranatestd init $VALIDATOR_NAME --home $NODE_HOME --chain-id $CHAIN_ID > /dev/null 2>&1
check_status "Node initialized"
echo ""

# Step 4: Copy genesis
echo -e "${YELLOW}=== Step 3: Copy Genesis File ===${NC}"
if [ ! -f "$HOME/.veranatest/config/genesis.json" ]; then
    echo -e "${RED}✗ Primary node genesis not found${NC}"
    echo "Please ensure primary node is set up at ~/.veranatest"
    exit 1
fi
cp $HOME/.veranatest/config/genesis.json $NODE_HOME/config/genesis.json
check_status "Genesis file copied"
echo ""

# Step 5: Configure persistent peers
echo -e "${YELLOW}=== Step 4: Configure Persistent Peers ===${NC}"
PRIMARY_NODE_ID=$(veranatestd status --node $PRIMARY_RPC 2>&1 | jq -r '.node_info.id')
if [ -z "$PRIMARY_NODE_ID" ] || [ "$PRIMARY_NODE_ID" = "null" ]; then
    echo -e "${RED}✗ Could not get primary node ID${NC}"
    echo "Is the primary node running?"
    exit 1
fi
PERSISTENT_PEERS="${PRIMARY_NODE_ID}@localhost:26656"
echo "Persistent Peers: $PERSISTENT_PEERS"
sed -i '' "s/^persistent_peers *=.*/persistent_peers = \"$PERSISTENT_PEERS\"/" $NODE_HOME/config/config.toml
check_status "Persistent peers configured"
echo ""

# Step 6: Configure ports to avoid conflicts
echo -e "${YELLOW}=== Step 5: Configure Ports (Avoid Conflicts) ===${NC}"
# ABCI: 26658 -> 26668
sed -i '' 's/proxy_app = "tcp:\/\/127.0.0.1:26658"/proxy_app = "tcp:\/\/127.0.0.1:26668"/' $NODE_HOME/config/config.toml
# RPC: 26657 -> 26667
sed -i '' 's/laddr = "tcp:\/\/127.0.0.1:26657"/laddr = "tcp:\/\/127.0.0.1:26667"/' $NODE_HOME/config/config.toml
# P2P: 26656 -> 26666
sed -i '' 's/laddr = "tcp:\/\/0.0.0.0:26656"/laddr = "tcp:\/\/0.0.0.0:26666"/' $NODE_HOME/config/config.toml
# gRPC: 9090 -> 9092
sed -i '' 's/address = "localhost:9090"/address = "localhost:9092"/' $NODE_HOME/config/app.toml
# gRPC-web: 9091 -> 9093
sed -i '' 's/address = "localhost:9091"/address = "localhost:9093"/' $NODE_HOME/config/app.toml
# API: 1317 -> 1327
sed -i '' 's/address = "tcp:\/\/localhost:1317"/address = "tcp:\/\/localhost:1327"/' $NODE_HOME/config/app.toml
# pprof: 6060 -> 6061 (in config.toml)
sed -i '' 's/pprof_laddr = "localhost:6060"/pprof_laddr = "localhost:6061"/' $NODE_HOME/config/config.toml
# Set minimum gas prices
sed -i '' 's/minimum-gas-prices = ""/minimum-gas-prices = "0.25uvna"/' $NODE_HOME/config/app.toml
check_status "Ports configured (RPC: 26667, P2P: 26666, gRPC: 9092, API: 1327, pprof: 6061)"
echo ""

# Step 7: Check and fund account
echo -e "${YELLOW}=== Step 6: Fund Validator Account ===${NC}"
VALIDATOR_ADDR=$(veranatestd keys show $VALIDATOR_KEY --keyring-backend test -a)
echo "Validator Address: $VALIDATOR_ADDR"

BALANCE=$(veranatestd query bank balances $VALIDATOR_ADDR --node $PRIMARY_RPC -o json 2>&1 | jq -r '.balances[] | select(.denom=="uvna") | .amount' || echo "0")
echo "Current Balance: ${BALANCE} uvna"

REQUIRED_AMOUNT=2000000000
if [ "$BALANCE" -lt "$REQUIRED_AMOUNT" ]; then
    echo "Insufficient funds. Sending ${REQUIRED_AMOUNT} uvna from cooluser..."
    veranatestd tx bank send cooluser $VALIDATOR_ADDR ${REQUIRED_AMOUNT}uvna \
      --from cooluser \
      --chain-id $CHAIN_ID \
      --keyring-backend test \
      --fees 500000uvna \
      --node $PRIMARY_RPC \
      -y > /dev/null 2>&1
    
    echo "Waiting 10 seconds for transaction to confirm..."
    sleep 10
    
    NEW_BALANCE=$(veranatestd query bank balances $VALIDATOR_ADDR --node $PRIMARY_RPC -o json 2>&1 | jq -r '.balances[] | select(.denom=="uvna") | .amount')
    echo "New Balance: ${NEW_BALANCE} uvna"
    check_status "Account funded"
else
    echo -e "${GREEN}✓ Account has sufficient funds${NC}"
fi
echo ""

# Step 8: Instructions to start node manually
echo -e "${YELLOW}=== Step 7: Start Node Manually ===${NC}"
echo -e "${GREEN}✓ Configuration complete!${NC}"
echo ""
echo "Node is ready to start. Open a NEW TERMINAL and run:"
echo ""
echo -e "${GREEN}  veranatestd start --home $NODE_HOME${NC}"
echo ""
echo "Leave that terminal running to see node logs."
echo ""
echo -e "${YELLOW}Press ENTER when the node is running and synced (catching_up: false)...${NC}"
read -r

# Step 9: Wait for node to sync
echo -e "${YELLOW}=== Step 8: Wait for Node to Sync ===${NC}"
echo "Waiting for node to start up (30 seconds)..."
sleep 30

MAX_WAIT=120
ELAPSED=0
while [ $ELAPSED -lt $MAX_WAIT ]; do
    CATCHING_UP=$(veranatestd status --node $NEW_NODE_RPC 2>&1 | jq -r '.sync_info.catching_up' 2>/dev/null || echo "true")
    LATEST_HEIGHT=$(veranatestd status --node $NEW_NODE_RPC 2>&1 | jq -r '.sync_info.latest_block_height' 2>/dev/null || echo "0")
    
    if [ "$CATCHING_UP" = "false" ] && [ "$LATEST_HEIGHT" != "0" ]; then
        echo -e "${GREEN}✓ Node is synced (Height: $LATEST_HEIGHT)${NC}"
        break
    fi
    
    echo "Syncing... (Height: $LATEST_HEIGHT, Catching up: $CATCHING_UP)"
    sleep 10
    ELAPSED=$((ELAPSED + 10))
done

if [ "$CATCHING_UP" = "true" ]; then
    echo -e "${YELLOW}⚠ Node is still catching up. You may need to wait longer before creating validator.${NC}"
    echo "Monitor sync status: veranatestd status --node $NEW_NODE_RPC 2>&1 | jq '.sync_info'"
fi
echo ""

# Step 10: Create validator
echo -e "${YELLOW}=== Step 9: Create Validator ===${NC}"
CONS_PUBKEY=$(veranatestd comet show-validator --home $NODE_HOME)
echo "Consensus PubKey: $CONS_PUBKEY"

# Create validator JSON file
cat > $NODE_HOME/validator_create.json << EOF
{
  "pubkey": $CONS_PUBKEY,
  "amount": "1000000000uvna",
  "moniker": "$VALIDATOR_NAME",
  "identity": "",
  "website": "",
  "security": "",
  "details": "Validator created via council governance",
  "commission-rate": "0.10",
  "commission-max-rate": "0.20",
  "commission-max-change-rate": "0.01",
  "min-self-delegation": "1"
}
EOF

echo "Creating validator..."
TX_RESULT=$(veranatestd tx staking create-validator $NODE_HOME/validator_create.json \
  --from=$VALIDATOR_KEY \
  --chain-id=$CHAIN_ID \
  --keyring-backend=test \
  --fees=500000uvna \
  --node=$PRIMARY_RPC \
  -y 2>&1)

TXHASH=$(echo "$TX_RESULT" | grep -o 'txhash: [A-Z0-9]*' | cut -d' ' -f2 || echo "")
if [ -z "$TXHASH" ]; then
    echo -e "${RED}✗ Failed to create validator${NC}"
    echo "$TX_RESULT"
    echo ""
    echo "Check if:"
    echo "  1. Validator is whitelisted"
    echo "  2. Account has sufficient funds"
    echo "  3. Node is synced"
    exit 1
fi

echo "Transaction Hash: $TXHASH"
check_status "Create validator transaction submitted"
echo ""

# Step 11: Wait and verify
echo -e "${YELLOW}=== Step 10: Verify Validator ===${NC}"
echo "Waiting 15 seconds for transaction to be included in block..."
sleep 15

echo "Checking validator status..."
VALIDATOR_INFO=$(veranatestd query staking validator $OPERATOR_ADDR --node $PRIMARY_RPC -o json 2>&1)
VALIDATOR_STATUS=$(echo "$VALIDATOR_INFO" | jq -r '.status' 2>/dev/null || echo "unknown")

if [ "$VALIDATOR_STATUS" = "BOND_STATUS_BONDED" ] || [ "$VALIDATOR_STATUS" = "BOND_STATUS_UNBONDING" ] || [ "$VALIDATOR_STATUS" = "BOND_STATUS_UNBONDED" ]; then
    echo -e "${GREEN}✓ Validator created successfully!${NC}"
    echo ""
    echo "Validator Details:"
    echo "$VALIDATOR_INFO" | jq '{moniker:.description.moniker, operator_address:.operator_address, status:.status, tokens:.tokens}'
else
    echo -e "${YELLOW}⚠ Could not verify validator status${NC}"
    echo "Check manually: veranatestd query staking validator $OPERATOR_ADDR --node $PRIMARY_RPC"
fi
echo ""

# Final summary
echo -e "${GREEN}╔════════════════════════════════════════════════════════════╗${NC}"
echo -e "${GREEN}║          Validator 2 Setup Complete! 🎉                   ║${NC}"
echo -e "${GREEN}╚════════════════════════════════════════════════════════════╝${NC}"
echo ""
echo "Node Information:"
echo "  - Name: $VALIDATOR_NAME"
echo "  - Home: $NODE_HOME"
echo "  - RPC: $NEW_NODE_RPC"
echo "  - Operator Address: $OPERATOR_ADDR"
echo ""
echo "⚠️  IMPORTANT: Keep the node terminal running!"
echo ""
echo "Useful Commands:"
echo "  # Check validator status"
echo "  veranatestd query staking validator $OPERATOR_ADDR --node $PRIMARY_RPC"
echo ""
echo "  # Check node sync status (in another terminal)"
echo "  veranatestd status --node $NEW_NODE_RPC 2>&1 | jq '.sync_info'"
echo ""
echo "  # To run node in background (if needed):"
echo "  veranatestd start --home $NODE_HOME > $NODE_HOME/node.log 2>&1 &"
echo "  echo \\\$! > $NODE_HOME/validator.pid"
echo ""
echo "  # Stop background node:"
echo "  kill \\\$(cat $NODE_HOME/validator.pid)"
echo ""
echo "  # List all validators"
echo "  veranatestd query staking validators --node $PRIMARY_RPC"
echo ""
echo "For detailed documentation, see: VALIDATOR_NODE_SETUP.md"

