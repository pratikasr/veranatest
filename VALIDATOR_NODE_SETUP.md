# Validator Node Setup Guide

Complete guide to set up and run a new validator node after being whitelisted through council governance.

## Table of Contents
- [Prerequisites](#prerequisites)
- [Quick Start](#quick-start)
- [Detailed Setup Steps](#detailed-setup-steps)
- [Troubleshooting](#troubleshooting)

---

## Prerequisites

### 1. Verify Whitelisting Status

Before setting up a node, ensure your validator is whitelisted:

```bash
# Check if your validator is in the whitelist
veranatestd query validatorregistry list-validator

# You should see your validator in the list:
# - index: validator2
#   member_id: member002
#   operator_address: cosmosvaloper1wh9djh6cfyncqzs4dp6g9ksmr63ekugvlkuk0n
#   status: active
```

### 2. Requirements

- `veranatestd` binary installed
- Validator whitelisted via council governance
- Access to existing node for genesis file and persistent peers

---

## Quick Start

```bash
# 1. Set up variables
export VALIDATOR_NAME="validator2"
export VALIDATOR_KEY="council-member-2"  # Key name in keyring
export NODE_HOME="$HOME/.veranatest-validator2"
export CHAIN_ID="vna-testnet-1"

# 2. Initialize new node
veranatestd init $VALIDATOR_NAME --home $NODE_HOME --chain-id $CHAIN_ID

# 3. Copy genesis from existing node
cp ~/.veranatest/config/genesis.json $NODE_HOME/config/genesis.json

# 4. Get persistent peers from existing node
PERSISTENT_PEERS=$(veranatestd status --node http://localhost:26657 2>&1 | jq -r '.node_info.id')@localhost:26656
sed -i '' "s/^persistent_peers *=.*/persistent_peers = \"$PERSISTENT_PEERS\"/" $NODE_HOME/config/config.toml

# 5. Configure different ports (to avoid conflicts) - CRITICAL!
sed -i '' 's/proxy_app = "tcp:\/\/127.0.0.1:26658"/proxy_app = "tcp:\/\/127.0.0.1:26668"/' $NODE_HOME/config/config.toml
sed -i '' 's/laddr = "tcp:\/\/127.0.0.1:26657"/laddr = "tcp:\/\/127.0.0.1:26667"/' $NODE_HOME/config/config.toml
sed -i '' 's/laddr = "tcp:\/\/0.0.0.0:26656"/laddr = "tcp:\/\/0.0.0.0:26666"/' $NODE_HOME/config/config.toml
sed -i '' 's/pprof_laddr = "localhost:6060"/pprof_laddr = "localhost:6061"/' $NODE_HOME/config/config.toml
sed -i '' 's/address = "localhost:9090"/address = "localhost:9092"/' $NODE_HOME/config/app.toml
sed -i '' 's/address = "localhost:9091"/address = "localhost:9093"/' $NODE_HOME/config/app.toml
sed -i '' 's/address = "tcp:\/\/localhost:1317"/address = "tcp:\/\/localhost:1327"/' $NODE_HOME/config/app.toml
sed -i '' 's/minimum-gas-prices = ""/minimum-gas-prices = "0.25uvna"/' $NODE_HOME/config/app.toml

# 6. Import or create validator key (if not already in keyring)
# Option A: If key already exists in keyring (from council setup)
veranatestd keys show $VALIDATOR_KEY --keyring-backend test

# Option B: If need to import from mnemonic
# veranatestd keys add $VALIDATOR_KEY --recover --keyring-backend test --home $NODE_HOME

# 7. Start the node (sync with network) - Open NEW terminal!
# In a separate terminal window, run:
veranatestd start --home $NODE_HOME

# Keep that terminal open to monitor logs
# In ANOTHER terminal, wait for sync then create validator
# See "Create Validator Transaction" section below
```

---

## Detailed Setup Steps

### Step 1: Initialize New Node

Create a new home directory for your validator node:

```bash
# Set up environment variables
export VALIDATOR_NAME="validator2"
export NODE_HOME="$HOME/.veranatest-validator2"
export CHAIN_ID="vna-testnet-1"

# Initialize the node
veranatestd init $VALIDATOR_NAME --home $NODE_HOME --chain-id $CHAIN_ID

# Output:
# {"app_message":{"auth":{"accounts":[],...}}}
```

This creates:
- `$NODE_HOME/config/` - Configuration files
- `$NODE_HOME/data/` - Blockchain data
- `$NODE_HOME/config/genesis.json` - Default genesis (will be replaced)
- `$NODE_HOME/config/config.toml` - Node configuration
- `$NODE_HOME/config/app.toml` - App configuration

---

### Step 2: Configure Genesis File

Copy the genesis file from the existing running node:

```bash
# Copy genesis from primary node
cp ~/.veranatest/config/genesis.json $NODE_HOME/config/genesis.json

# Verify genesis hash matches
veranatestd genesis validate-genesis --home $NODE_HOME

# Compare with primary node
veranatestd genesis validate-genesis --home ~/.veranatest
```

**Important:** Both nodes must have the **exact same genesis file** to join the network.

---

### Step 3: Configure Persistent Peers

Get the node ID from the existing validator:

```bash
# Get primary validator node ID
PRIMARY_NODE_ID=$(veranatestd status --node http://localhost:26657 2>&1 | jq -r '.node_info.id')
echo "Primary Node ID: $PRIMARY_NODE_ID"

# Output example: 7f3473f3c8f1c8e0c8f1c8e0c8f1c8e0c8f1c8e0

# Set persistent peers
PERSISTENT_PEERS="${PRIMARY_NODE_ID}@localhost:26656"
echo "Persistent Peers: $PERSISTENT_PEERS"

# Update config.toml
sed -i '' "s/^persistent_peers *=.*/persistent_peers = \"$PERSISTENT_PEERS\"/" $NODE_HOME/config/config.toml

# Verify
grep "persistent_peers" $NODE_HOME/config/config.toml
```

**For Production (Different Machines):**
```bash
# If primary node is on different machine
PRIMARY_IP="192.168.1.100"  # Replace with actual IP
PERSISTENT_PEERS="${PRIMARY_NODE_ID}@${PRIMARY_IP}:26656"
```

---

### Step 4: Configure Node Ports (Avoid Conflicts)

Since you're running multiple nodes on the same machine, configure different ports:

```bash
# ABCI application port (26658 -> 26668)
sed -i '' 's/proxy_app = "tcp:\/\/127.0.0.1:26658"/proxy_app = "tcp:\/\/127.0.0.1:26668"/' $NODE_HOME/config/config.toml

# RPC port (26657 -> 26667)
sed -i '' 's/laddr = "tcp:\/\/127.0.0.1:26657"/laddr = "tcp:\/\/127.0.0.1:26667"/' $NODE_HOME/config/config.toml

# P2P port (26656 -> 26666)
sed -i '' 's/laddr = "tcp:\/\/0.0.0.0:26656"/laddr = "tcp:\/\/0.0.0.0:26666"/' $NODE_HOME/config/config.toml

# pprof port (6060 -> 6061) ⚠️ IMPORTANT: This causes conflicts!
sed -i '' 's/pprof_laddr = "localhost:6060"/pprof_laddr = "localhost:6061"/' $NODE_HOME/config/config.toml

# gRPC port (9090 -> 9092)
sed -i '' 's/address = "localhost:9090"/address = "localhost:9092"/' $NODE_HOME/config/app.toml

# gRPC-web port (9091 -> 9093)
sed -i '' 's/address = "localhost:9091"/address = "localhost:9093"/' $NODE_HOME/config/app.toml

# API port (1317 -> 1327)
sed -i '' 's/address = "tcp:\/\/localhost:1317"/address = "tcp:\/\/localhost:1327"/' $NODE_HOME/config/app.toml

# Set minimum gas prices (required!)
sed -i '' 's/minimum-gas-prices = ""/minimum-gas-prices = "0.25uvna"/' $NODE_HOME/config/app.toml

# Verify ports
echo "=== New Node Ports ==="
grep "proxy_app" $NODE_HOME/config/config.toml | grep -v "#"
grep "laddr" $NODE_HOME/config/config.toml | grep "26" | grep -v "#"
grep "pprof_laddr" $NODE_HOME/config/config.toml | grep -v "#"
grep "address.*9092" $NODE_HOME/config/app.toml
grep "address.*1327" $NODE_HOME/config/app.toml
grep "minimum-gas-prices" $NODE_HOME/config/app.toml | grep -v "#"
```

**Port Summary:**
```
validator1 (primary):          validator2 (new):
- ABCI:      26658             - ABCI:      26668
- RPC:       26657             - RPC:       26667
- P2P:       26656             - P2P:       26666
- pprof:     6060              - pprof:     6061  ⚠️
- gRPC:      9090              - gRPC:      9092
- gRPC-web:  9091              - gRPC-web:  9093
- API:       1317              - API:       1327
```

⚠️ **CRITICAL**: The pprof port (6060) MUST be different or validator1 will halt!

---

### Step 5: Set Up Validator Key

Your validator key should already exist from the council setup:

```bash
export VALIDATOR_KEY="council-member-2"

# Check if key exists
veranatestd keys show $VALIDATOR_KEY --keyring-backend test

# Output:
# - address: cosmos1wh9djh6cfyncqzs4dp6g9ksmr63ekugv56t2xd
#   name: council-member-2
#   pubkey: '{"@type":"/cosmos.crypto.secp256k1.PubKey","key":"..."}'
#   type: local

# Get validator operator address
OPERATOR_ADDR=$(veranatestd keys show $VALIDATOR_KEY --bech val --keyring-backend test -a)
echo "Operator Address: $OPERATOR_ADDR"
# Should match: cosmosvaloper1wh9djh6cfyncqzs4dp6g9ksmr63ekugvlkuk0n
```

**If Key Doesn't Exist:**
```bash
# Option A: Create new key
veranatestd keys add $VALIDATOR_KEY --keyring-backend test --home $NODE_HOME

# Option B: Import from mnemonic
veranatestd keys add $VALIDATOR_KEY --recover --keyring-backend test --home $NODE_HOME
# Enter your 24-word mnemonic when prompted
```

---

### Step 6: Start the Node

**⚠️ IMPORTANT**: Start the node in a **separate terminal** so you can monitor logs in real-time.

**Option A: Foreground (Recommended)**
```bash
# Open a NEW terminal and run:
veranatestd start --home $NODE_HOME

# Keep this terminal open to see live logs
# Use Ctrl+C to stop the node
```

**Option B: Background (Advanced)**
```bash
# Only use if you understand process management
veranatestd start --home $NODE_HOME > $NODE_HOME/node.log 2>&1 &

# Save PID for later
echo $! > $NODE_HOME/validator.pid

# View logs:
tail -f $NODE_HOME/node.log

# Stop later:
kill $(cat $NODE_HOME/validator.pid)
```

**For this guide, we'll use Option A (foreground) for easier debugging.**

**Monitor Sync Status:**
```bash
# Check sync status
veranatestd status --node http://localhost:26667 2>&1 | jq '.sync_info'

# Output:
# {
#   "latest_block_height": "1234",
#   "catching_up": false  # Wait until this is false
# }

# Watch sync progress
watch -n 5 'veranatestd status --node http://localhost:26667 2>&1 | jq ".sync_info.latest_block_height, .sync_info.catching_up"'
```

**Wait for `catching_up: false` before creating validator!**

---

### Step 7: Fund the Validator Account

Ensure your validator account has funds for the transaction:

```bash
# Check balance
veranatestd query bank balances $(veranatestd keys show $VALIDATOR_KEY --keyring-backend test -a) --node http://localhost:26657

# If balance is zero, send funds from primary validator
veranatestd tx bank send validator1 $(veranatestd keys show $VALIDATOR_KEY --keyring-backend test -a) 1000000000uvna \
  --from validator1 \
  --chain-id $CHAIN_ID \
  --keyring-backend test \
  --fees 500000uvna \
  --node http://localhost:26657 \
  -y

# Wait for transaction, then verify balance
veranatestd query bank balances $(veranatestd keys show $VALIDATOR_KEY --keyring-backend test -a) --node http://localhost:26657
```

---

### Step 8: Create Validator Transaction

Once synced and funded, create the validator:

```bash
# Get consensus public key from new node
CONS_PUBKEY=$(veranatestd comet show-validator --home $NODE_HOME)
echo "Consensus PubKey: $CONS_PUBKEY"

# Create validator
veranatestd tx staking create-validator \
  --amount=1000000000uvna \
  --pubkey="$CONS_PUBKEY" \
  --moniker="$VALIDATOR_NAME" \
  --chain-id=$CHAIN_ID \
  --commission-rate="0.10" \
  --commission-max-rate="0.20" \
  --commission-max-change-rate="0.01" \
  --min-self-delegation="1" \
  --from=$VALIDATOR_KEY \
  --keyring-backend=test \
  --fees=500000uvna \
  --node=http://localhost:26657 \
  -y

# Output:
# txhash: ABC123...
```

**Important Parameters:**
- `--amount`: Amount to self-delegate (must have this much in account)
- `--pubkey`: Consensus public key from your node
- `--moniker`: Validator name (visible in explorer)
- `--from`: Your validator key name
- `--node`: Point to PRIMARY node's RPC (since you're submitting tx to network)

---

### Step 9: Verify Validator is Active

```bash
# Check validator list
veranatestd query staking validators --node http://localhost:26657

# Should see both validators:
# - validator1 (cosmosvaloper16mzeyu9l6kua2cdg9x0jk5g6e7h0kk8q0qpggj)
# - validator2 (cosmosvaloper1wh9djh6cfyncqzs4dp6g9ksmr63ekugvlkuk0n)

# Check specific validator
veranatestd query staking validator $OPERATOR_ADDR --node http://localhost:26657

# Check validator status from new node's perspective
veranatestd query staking validator $OPERATOR_ADDR --node http://localhost:26667

# Check validator is signing blocks
veranatestd query slashing signing-info $(veranatestd comet show-validator --home $NODE_HOME | jq -r '.key') --node http://localhost:26657
```

✅ **Success!** Your validator is now active and participating in consensus.

---

## Node Management

### Stop Node
```bash
# If started in background
kill $(cat $NODE_HOME/validator.pid)

# Or find process
ps aux | grep veranatest-validator2
kill <PID>
```

### Start Node
```bash
veranatestd start --home $NODE_HOME > $NODE_HOME/node.log 2>&1 &
echo $! > $NODE_HOME/validator.pid
```

### View Logs
```bash
# Real-time logs
tail -f $NODE_HOME/node.log

# Search for errors
grep -i error $NODE_HOME/node.log

# Check consensus
grep "consensus" $NODE_HOME/node.log | tail -20
```

### Query Node Status
```bash
# Node info
veranatestd status --node http://localhost:26667 2>&1 | jq '.'

# Validator info
veranatestd query staking validator $OPERATOR_ADDR --node http://localhost:26667

# Network info
veranatestd status --node http://localhost:26667 2>&1 | jq '.node_info.network'
```

---

## Troubleshooting

### Node Won't Start

**Error: `address already in use`**
```bash
# Check what's using the port
lsof -i :26666  # or whichever port is conflicting

# Kill the process or use different ports
```

**Error: `genesis file not found`**
```bash
# Ensure genesis exists
ls -la $NODE_HOME/config/genesis.json

# Copy from primary node
cp ~/.veranatest/config/genesis.json $NODE_HOME/config/genesis.json
```

### Node Not Syncing

**Error: `no peers available`**
```bash
# Check persistent peers
grep "persistent_peers" $NODE_HOME/config/config.toml

# Verify primary node is reachable
curl http://localhost:26657/status

# Check P2P connection
netstat -an | grep 26656
```

### Create Validator Fails

**Error: `validator not whitelisted`**
```bash
# Verify your validator is whitelisted
veranatestd query validatorregistry list-validator --node http://localhost:26657

# Check operator address matches
echo "Expected: cosmosvaloper1wh9djh6cfyncqzs4dp6g9ksmr63ekugvlkuk0n"
echo "Your address: $(veranatestd keys show $VALIDATOR_KEY --bech val --keyring-backend test -a)"
```

**Error: `insufficient funds`**
```bash
# Check balance
veranatestd query bank balances $(veranatestd keys show $VALIDATOR_KEY --keyring-backend test -a) --node http://localhost:26657

# Send funds from primary validator
veranatestd tx bank send validator1 $(veranatestd keys show $VALIDATOR_KEY --keyring-backend test -a) 2000000000uvna \
  --from validator1 \
  --chain-id $CHAIN_ID \
  --keyring-backend test \
  --fees 500000uvna \
  --node http://localhost:26657 \
  -y
```

**Error: `pubkey already exists`**
```bash
# This means a validator with this consensus pubkey already exists
# You need to use a different home directory / consensus key
```

### Validator Not Signing Blocks

```bash
# Check node is synced
veranatestd status --node http://localhost:26667 2>&1 | jq '.sync_info.catching_up'
# Must be: false

# Check validator status
veranatestd query staking validator $OPERATOR_ADDR --node http://localhost:26657 | jq '.status'
# Must be: BOND_STATUS_BONDED

# Check signing info
veranatestd query slashing signing-info $(veranatestd comet show-validator --home $NODE_HOME | jq -r '.key') --node http://localhost:26657

# Restart node if needed
kill $(cat $NODE_HOME/validator.pid)
veranatestd start --home $NODE_HOME > $NODE_HOME/node.log 2>&1 &
echo $! > $NODE_HOME/validator.pid
```

---

## Complete Example: Setting Up validator2

Here's the complete flow for validator2 (council-member-2):

```bash
#!/bin/bash

# Variables
export VALIDATOR_NAME="validator2"
export VALIDATOR_KEY="council-member-2"
export NODE_HOME="$HOME/.veranatest-validator2"
export CHAIN_ID="vna-testnet-1"

echo "=== Step 1: Initialize Node ==="
veranatestd init $VALIDATOR_NAME --home $NODE_HOME --chain-id $CHAIN_ID

echo "=== Step 2: Copy Genesis ==="
cp ~/.veranatest/config/genesis.json $NODE_HOME/config/genesis.json

echo "=== Step 3: Configure Persistent Peers ==="
PRIMARY_NODE_ID=$(veranatestd status --node http://localhost:26657 2>&1 | jq -r '.node_info.id')
PERSISTENT_PEERS="${PRIMARY_NODE_ID}@localhost:26656"
sed -i '' "s/^persistent_peers *=.*/persistent_peers = \"$PERSISTENT_PEERS\"/" $NODE_HOME/config/config.toml

echo "=== Step 4: Configure Ports ==="
sed -i '' 's/proxy_app = "tcp:\/\/127.0.0.1:26658"/proxy_app = "tcp:\/\/127.0.0.1:26668"/' $NODE_HOME/config/config.toml
sed -i '' 's/laddr = "tcp:\/\/127.0.0.1:26657"/laddr = "tcp:\/\/127.0.0.1:26667"/' $NODE_HOME/config/config.toml
sed -i '' 's/laddr = "tcp:\/\/0.0.0.0:26656"/laddr = "tcp:\/\/0.0.0.0:26666"/' $NODE_HOME/config/config.toml
sed -i '' 's/address = "localhost:9090"/address = "localhost:9092"/' $NODE_HOME/config/app.toml
sed -i '' 's/address = "localhost:9091"/address = "localhost:9093"/' $NODE_HOME/config/app.toml
sed -i '' 's/address = "tcp:\/\/localhost:1317"/address = "tcp:\/\/localhost:1327"/' $NODE_HOME/config/app.toml

echo "=== Step 5: Verify Key Exists ==="
veranatestd keys show $VALIDATOR_KEY --keyring-backend test
OPERATOR_ADDR=$(veranatestd keys show $VALIDATOR_KEY --bech val --keyring-backend test -a)
echo "Operator Address: $OPERATOR_ADDR"

echo "=== Step 6: Fund Account ==="
VALIDATOR_ADDR=$(veranatestd keys show $VALIDATOR_KEY --keyring-backend test -a)
veranatestd tx bank send validator1 $VALIDATOR_ADDR 2000000000uvna \
  --from validator1 \
  --chain-id $CHAIN_ID \
  --keyring-backend test \
  --fees 500000uvna \
  --node http://localhost:26657 \
  -y

echo "Waiting 10 seconds for transaction..."
sleep 10

echo "=== Step 7: Start Node ==="
echo "Starting node in background..."
veranatestd start --home $NODE_HOME > $NODE_HOME/node.log 2>&1 &
echo $! > $NODE_HOME/validator.pid

echo "=== Waiting for node to sync (60 seconds) ==="
sleep 60

echo "=== Step 8: Create Validator ==="
CONS_PUBKEY=$(veranatestd comet show-validator --home $NODE_HOME)
veranatestd tx staking create-validator \
  --amount=1000000000uvna \
  --pubkey="$CONS_PUBKEY" \
  --moniker="$VALIDATOR_NAME" \
  --chain-id=$CHAIN_ID \
  --commission-rate="0.10" \
  --commission-max-rate="0.20" \
  --commission-max-change-rate="0.01" \
  --min-self-delegation="1" \
  --from=$VALIDATOR_KEY \
  --keyring-backend=test \
  --fees=500000uvna \
  --node=http://localhost:26657 \
  -y

echo "=== Setup Complete! ==="
echo "Node PID: $(cat $NODE_HOME/validator.pid)"
echo "Operator Address: $OPERATOR_ADDR"
echo ""
echo "Verify with:"
echo "  veranatestd query staking validators --node http://localhost:26657"
echo "  veranatestd status --node http://localhost:26667"
```

Save as `setup_validator2.sh` and run:
```bash
chmod +x setup_validator2.sh
./setup_validator2.sh
```

---

## References

- **Validator Whitelisting**: [VALIDATOR_WHITELIST.md](./VALIDATOR_WHITELIST.md)
- **Council Governance**: [COUNCIL_GOVERNANCE.md](./COUNCIL_GOVERNANCE.md)
- **Main README**: [README.md](./readme.md)

---

## Summary

✅ **Validator Node Setup Checklist:**

1. ✅ Verify validator is whitelisted
2. ✅ Initialize new node with unique home directory
3. ✅ Copy genesis from existing node
4. ✅ Configure persistent peers
5. ✅ Set up non-conflicting ports
6. ✅ Verify validator key exists
7. ✅ Fund validator account
8. ✅ Start node and wait for sync
9. ✅ Create validator with consensus pubkey
10. ✅ Verify validator is active and signing blocks

**Your validator should now be active and participating in consensus!** 🎉

