# Validator Whitelist Setup Guide

## How It Works

### Module Initialization Order (CRITICAL!)

```
InitChain Process (Block Height 0):
1. validatorregistry.InitGenesis() → Loads whitelist into KV store ✅
2. genutil.InitGenesis() → Processes gentxs (MsgCreateValidator) → Whitelist check runs ✅
```

**Key Point**: `validatorregistry` initializes **BEFORE** `genutil`, so the whitelist is available when gentxs are processed.

## Setup Process

### Step 1: Update Genesis File with Whitelisted Validators

**BEFORE running `setup_validator.sh`**, you need to know which validator addresses to whitelist.

#### Get Your Validator Operator Address

```bash
# First, create the key temporarily to get the address
echo "pink glory help gown abstract eight nice crazy forward ketchup skill cheese" | \
  veranatestd keys add cooluser --recover --keyring-backend test

# Get the operator address
veranatestd keys show cooluser --bech32=val --keyring-backend test -a
```

Output example:
```
cosmosvaloper16mzeyu9l6kua2cdg9x0jk5g6e7h0kk8q0qpggj
```

**Save this address!** You'll need it for the genesis file.

### Step 2: Run Setup Script

The script does NOT automatically whitelist validators. It just creates the chain.

```bash
rm -rf ~/.veranatest  # Clean old data
./setup_validator.sh
```

**This will FAIL** with error:
```
validator address cosmosvaloper16mzeyu9l6kua2cdg9x0jk5g6e7h0kk8q0qpggj is not whitelisted
```

This is **CORRECT behavior**! The whitelist is enforced.

### Step 3: Add Validator to Genesis Whitelist

After the script creates the genesis file but before collecting gentxs, you need to add validators to the whitelist.

#### Option A: Modify the Script (Recommended)

Update `setup_validator.sh` to add whitelist BEFORE `collect-gentxs`:

```bash
# After line 74 (after gentx creation), add:

# Get validator operator address
VALIDATOR_OPERATOR_ADDR=$($BINARY keys show $VALIDATOR_NAME --bech32=val --keyring-backend test -a)
log "Validator operator address: $VALIDATOR_OPERATOR_ADDR"

# Update genesis with whitelisted validator
log "Adding validator to whitelist..."
if command -v jq &> /dev/null; then
    # Create temporary whitelist entry
    cat > /tmp/validator_whitelist.json <<EOF
{
  "index": "validator1",
  "member_id": "member001",
  "operator_address": "$VALIDATOR_OPERATOR_ADDR",
  "consensus_pubkey": "",
  "status": "active",
  "term_end": 0
}
EOF

    # Add to genesis
    jq --slurpfile val /tmp/validator_whitelist.json \
       '.app_state.validatorregistry.validator_map += $val' \
       "$GENESIS_JSON_PATH" > /tmp/genesis_temp.json && \
       mv /tmp/genesis_temp.json "$GENESIS_JSON_PATH"
    
    rm /tmp/validator_whitelist.json
    log "Validator added to whitelist"
else
    log "ERROR: jq is required. Install it: brew install jq"
    exit 1
fi

# Then continue with collect-gentxs (line 110)
```

#### Option B: Manual Genesis Edit

After the script fails, manually edit `~/.veranatest/config/genesis.json`:

Find the `validatorregistry` section and add your validator:

```json
{
  "app_state": {
    "validatorregistry": {
      "params": {},
      "validator_map": [
        {
          "index": "validator1",
          "member_id": "member001",
          "operator_address": "cosmosvaloper16mzeyu9l6kua2cdg9x0jk5g6e7h0kk8q0qpggj",
          "consensus_pubkey": "",
          "status": "active",
          "term_end": 0
        }
      ]
    }
  }
}
```

Then manually run:
```bash
veranatestd genesis collect-gentxs
veranatestd genesis validate
veranatestd start
```

## Correct Workflow Summary

### ❌ WRONG (Current script - won't work):
```
1. Create gentx
2. Collect gentxs → Chain starts → Validator not whitelisted → FAIL ❌
```

### ✅ CORRECT:
```
1. Create gentx
2. Get validator operator address
3. Add to genesis whitelist (validatorregistry.validator_map)
4. Collect gentxs → Chain starts → Whitelist check passes → SUCCESS ✅
```

## Testing the Whitelist

### Test 1: Whitelisted Validator at Genesis (Should Succeed)

1. Add validator to genesis whitelist
2. Run setup script
3. Chain should start successfully ✅

### Test 2: Non-Whitelisted Validator at Genesis (Should Fail)

1. Don't add validator to genesis whitelist
2. Run setup script
3. Chain should fail with "not whitelisted" error ✅

### Test 3: New Validator After Genesis (Should Fail)

After chain is running:

```bash
# Create new validator key
veranatestd keys add unauthorized --keyring-backend test

# Fund the account
veranatestd tx bank send cooluser $(veranatestd keys show unauthorized -a --keyring-backend test) 1000000000uvna --from cooluser --fees 500000uvna --keyring-backend test --chain-id vna-testnet-1 -y

# Try to create validator (should FAIL)
veranatestd tx staking create-validator \
  --amount=1000000uvna \
  --pubkey=$(veranatestd tendermint show-validator) \
  --moniker="unauthorized" \
  --commission-rate="0.10" \
  --commission-max-rate="0.20" \
  --commission-max-change-rate="0.01" \
  --min-self-delegation="1" \
  --from=unauthorized \
  --chain-id=vna-testnet-1 \
  --keyring-backend=test \
  --fees=500000uvna
```

**Expected**: ❌ Error: "validator address is not whitelisted"

### Test 4: Add Validator to Whitelist and Create (Should Succeed)

```bash
# Add to whitelist via transaction
veranatestd tx validatorregistry onboard-validator \
  validator2 \
  member002 \
  $(veranatestd keys show unauthorized --bech32=val --keyring-backend test -a) \
  "" \
  active \
  0 \
  --from=cooluser \
  --chain-id=vna-testnet-1 \
  --fees=500000uvna \
  --keyring-backend=test -y

# Now create validator (should SUCCEED)
veranatestd tx staking create-validator \
  --amount=1000000uvna \
  --pubkey=$(veranatestd tendermint show-validator) \
  --moniker="now-authorized" \
  --commission-rate="0.10" \
  --commission-max-rate="0.20" \
  --commission-max-change-rate="0.01" \
  --min-self-delegation="1" \
  --from=unauthorized \
  --chain-id=vna-testnet-1 \
  --keyring-backend=test \
  --fees=500000uvna -y
```

## Updated Setup Script

Here's the complete updated script that adds validators to the whitelist automatically:

Save as `setup_validator_with_whitelist.sh`:

```bash
#!/bin/bash
set -e

# ... (keep all the existing variables and functions)

# After gentx creation (after line 74), add:

# Get validator operator address
VALIDATOR_OPERATOR_ADDR=$($BINARY keys show $VALIDATOR_NAME --bech32=val --keyring-backend test -a)
log "Validator operator address: $VALIDATOR_OPERATOR_ADDR"

# Update genesis with whitelisted validator BEFORE collecting gentxs
log "Adding validator to whitelist..."
if command -v jq &> /dev/null; then
    jq --arg addr "$VALIDATOR_OPERATOR_ADDR" \
       '.app_state.validatorregistry.validator_map = [{
         "index": "validator1",
         "member_id": "member001",
         "operator_address": $addr,
         "consensus_pubkey": "",
         "status": "active",
         "term_end": 0
       }]' \
       "$GENESIS_JSON_PATH" > /tmp/genesis_temp.json && \
       mv /tmp/genesis_temp.json "$GENESIS_JSON_PATH"
    
    log "Validator added to whitelist"
else
    log "ERROR: jq is required. Install it with: brew install jq (macOS) or apt install jq (Linux)"
    exit 1
fi

# Then continue with the rest of the script (collect-gentxs, etc.)
```

## Key Points

1. ✅ **Whitelist is enforced at ALL block heights** (including genesis)
2. ✅ **validatorregistry initializes BEFORE genutil** (fixed in app_config.go)
3. ✅ **Validators must be in whitelist BEFORE gentxs are collected**
4. ✅ **Script must be updated to add validators to whitelist automatically**
5. ❌ **Current script will FAIL because it doesn't whitelist the validator**

## Next Steps

1. Install `jq`: `brew install jq` (macOS) or `apt install jq` (Linux)
2. Update your `setup_validator.sh` script with the whitelist addition code
3. Clean and restart: `rm -rf ~/.veranatest && ./setup_validator.sh`
4. Chain should start successfully with whitelisted validator!

