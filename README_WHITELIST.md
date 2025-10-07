# Validator Whitelist - Quick Start

## What Changed?

Your chain now uses a **proper KV store** to manage whitelisted validators, following standard Cosmos SDK patterns. The ante handler accesses the `validatorregistry` module keeper to check if validators are whitelisted.

## Quick Start

### Step 1: Add Validators to Whitelist via Genesis

Before starting your chain, add whitelisted validators to the genesis file.

#### Get Your Validator Address
```bash
veranatestd keys show cooluser --bech32=val --keyring-backend test -a
```

This will output something like:
```
cosmosvaloper1rkz2eeu3rveg7u6srnsdkcjqmwc32kyl9565pm
```

#### Edit Genesis File

Edit `~/.veranatest/config/genesis.json` and find the `validatorregistry` section under `app_state`:

```json
{
  "app_state": {
    "validatorregistry": {
      "params": {},
      "validator_map": [
        {
          "index": "validator1",
          "member_id": "member001",
          "operator_address": "cosmosvaloper1rkz2eeu3rveg7u6srnsdkcjqmwc32kyl9565pm",
          "consensus_pubkey": "",
          "status": "active",
          "term_end": 0
        }
      ]
    }
  }
}
```

### Step 2: Start Your Chain

```bash
veranatestd start
```

### Step 3: Verify Whitelist

```bash
# List all whitelisted validators
veranatestd query validatorregistry list-validator

# Expected output:
# validator_map:
# - index: validator1
#   member_id: member001
#   operator_address: cosmosvaloper1rkz2eeu3rveg7u6srnsdkcjqmwc32kyl9565pm
#   ...
```

### Step 4: Test Creating a Validator

```bash
# This should succeed (if your address is whitelisted)
veranatestd tx staking create-validator \
  --amount=1000000uvna \
  --pubkey=$(veranatestd tendermint show-validator) \
  --moniker="my-validator" \
  --commission-rate="0.10" \
  --commission-max-rate="0.20" \
  --commission-max-change-rate="0.01" \
  --min-self-delegation="1" \
  --from=cooluser \
  --chain-id=vna-testnet-1 \
  --keyring-backend=test \
  --fees=500000uvna
```

## How It Works

```
Transaction (MsgCreateValidator)
    ↓
Ante Handler
    ↓
ValidatorWhitelistDecorator
    ↓
validatorregistryKeeper.IsValidatorWhitelisted()
    ↓
Checks KV Store
    ↓
✅ Allow if whitelisted
❌ Reject if not whitelisted
```

## Adding More Validators

### Option 1: Before Chain Start (Genesis)

Add to `validator_map` in genesis.json as shown above.

### Option 2: After Chain Start (Transaction)

```bash
veranatestd tx validatorregistry onboard-validator \
  <index> \
  <member-id> \
  <operator-address> \
  "" \
  active \
  0 \
  --from=<authority> \
  --chain-id=vna-testnet-1
```

**Note**: Check who has authority to add validators:
```bash
veranatestd query validatorregistry params
```

## Files Modified

1. ✅ `x/validatorregistry/keeper/keeper.go` - Added `IsValidatorWhitelisted()` method
2. ✅ `ante/validator_whitelist.go` - Uses keeper instead of hardcoded address
3. ✅ `ante/ante.go` - Accepts keeper as parameter
4. ✅ `app/app.go` - Passes keeper to ante handler

## Common Questions

### Q: Can I still use my existing validator?
**A:** Yes! Just add its operator address to the genesis whitelist.

### Q: How do I get the operator address?
**A:** `veranatestd keys show <key-name> --bech32=val --keyring-backend test -a`

### Q: Can I add validators after the chain starts?
**A:** Yes, use the `onboard-validator` transaction (if you have authority).

### Q: What if I try to create a validator that's not whitelisted?
**A:** The transaction will fail with error: "validator address is not whitelisted"

### Q: Is this the standard Cosmos SDK way?
**A:** Yes! This follows the same pattern as how `auth` and `bank` modules are used in ante handlers.

## For More Details

See [VALIDATOR_WHITELIST_IMPLEMENTATION.md](./VALIDATOR_WHITELIST_IMPLEMENTATION.md) for complete technical documentation.

## Summary

✅ **Done**: Validator whitelist now uses KV store  
✅ **Done**: Ante handler accesses store via keeper  
✅ **Done**: Follows standard Cosmos SDK patterns  
✅ **Done**: Configurable via genesis file  
✅ **Done**: No hardcoded addresses  

Your implementation is complete and ready to use! 🎉

