# Validator Whitelist

## Overview

The Verana Network implements a council-based governance model for validators. Only whitelisted validators can create validator nodes on the network. This is enforced through an ante handler that checks the `validatorregistry` module before allowing `MsgCreateValidator` transactions.

**Key Features:**
- ✅ Dynamic whitelist stored in KV store (no hardcoded addresses)
- ✅ Follows standard Cosmos SDK patterns (keeper injection)
- ✅ Configurable via genesis or council governance
- ✅ Enforced at all block heights (including genesis)
- ✅ Council-based voting for validator onboarding (2/3 threshold)

> **📖 For Council Governance:** See [COUNCIL_GOVERNANCE.md](./COUNCIL_GOVERNANCE.md) for the complete guide on:
> - Setting up the council
> - Adding validators through governance proposals
> - Voting and execution process with tested examples

## Architecture

```
Transaction (MsgCreateValidator)
    ↓
Ante Handler
    ↓
ValidatorWhitelistDecorator
    ↓
validatorregistryKeeper.IsValidatorWhitelisted()
    ↓
Checks KV Store (validator_map)
    ↓
✅ Allow if whitelisted
❌ Reject if not whitelisted
```

### Implementation Details

The implementation follows the same pattern used by `auth` and `bank` modules in Cosmos SDK:

1. **KV Store** - `validatorregistry` module stores whitelisted validators
2. **Keeper** - Provides `IsValidatorWhitelisted()` method
3. **Ante Decorator** - Receives keeper via dependency injection
4. **App Setup** - Injects keeper into ante handler

**Files Modified:**
- `x/validatorregistry/keeper/keeper.go` - Added `IsValidatorWhitelisted()` method
- `ante/validator_whitelist.go` - Uses keeper instead of hardcoded address
- `ante/ante.go` - Accepts keeper as parameter
- `app/app.go` - Passes keeper to ante handler

## Quick Start

### 1. Get Your Validator Address

```bash
veranatestd keys show cooluser --bech32=val --keyring-backend test -a
```

Output example:
```
cosmosvaloper1rkz2eeu3rveg7u6srnsdkcjqmwc32kyl9565pm
```

### 2. Setup and Whitelist

The `setup_validator.sh` script automatically whitelists your validator during chain initialization:

```bash
rm -rf ~/.veranatest  # Clean old data
./setup_validator.sh
```

The script:
1. Creates your validator key
2. Gets the operator address
3. Adds it to genesis whitelist (`validatorregistry.validator_map`)
4. Collects gentxs and starts the chain

### 3. Verify Whitelist

```bash
# List all whitelisted validators
veranatestd query validatorregistry list-validator

# Check specific validator
veranatestd query validatorregistry show-validator validator1
```

## Managing the Whitelist

### Adding Validators via Genesis

Edit `~/.veranatest/config/genesis.json` before starting the chain:

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
        },
        {
          "index": "validator2",
          "member_id": "member002",
          "operator_address": "cosmosvaloper1abc...",
          "consensus_pubkey": "",
          "status": "active",
          "term_end": 0
        }
      ]
    }
  }
}
```

### Adding Validators via Council Governance

After the chain is running, validators are added through the **council governance process**:

1. **Council submits proposal** to onboard validator
2. **Council votes** (requires 2/3 approval)
3. **Execute proposal** to add validator to whitelist

**Complete Guide:** See [COUNCIL_GOVERNANCE.md](./COUNCIL_GOVERNANCE.md#how-to-add-validators-to-whitelist) for:
- Full council setup instructions
- Step-by-step proposal creation
- Voting and execution process
- Tested examples with transaction hashes

**Quick Command Reference:**
```bash
# Submit proposal (creates JSON first)
veranatestd tx group submit-proposal add_validator_proposal.json --from council-member-1 ...

# Vote on proposal
veranatestd tx group vote <proposal-id> <voter-address> VOTE_OPTION_YES "..." --from ...

# Execute after 2/3 approval
veranatestd tx group exec <proposal-id> --from ...
```

**Note:** The `validatorregistry` module doesn't have queryable params. State is in the validator map.

## Testing

### Test 1: Whitelisted Validator (Should Succeed)

```bash
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

**Expected:** ✅ Transaction succeeds

### Test 2: Non-Whitelisted Validator (Should Fail)

```bash
# Create unauthorized key
veranatestd keys add unauthorized --keyring-backend test

# Fund the account
veranatestd tx bank send cooluser \
  $(veranatestd keys show unauthorized -a --keyring-backend test) \
  1000000000uvna \
  --from cooluser \
  --fees 500000uvna \
  --keyring-backend test \
  --chain-id vna-testnet-1 -y

# Try to create validator
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

**Expected:** ❌ Error: "validator address is not whitelisted"

### Test 3: Add to Whitelist Then Create (Should Succeed)

```bash
# Get unauthorized validator's operator address
UNAUTHORIZED_ADDR=$(veranatestd keys show unauthorized --bech32=val --keyring-backend test -a)

# Add to whitelist
veranatestd tx validatorregistry onboard-validator \
  validator2 \
  member002 \
  $UNAUTHORIZED_ADDR \
  "" \
  active \
  0 \
  --from=cooluser \
  --chain-id=vna-testnet-1 \
  --fees=500000uvna \
  --keyring-backend=test -y

# Now create validator
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

**Expected:** ✅ Transaction succeeds

## Code Implementation

### Keeper Method

```go
// x/validatorregistry/keeper/keeper.go
func (k Keeper) IsValidatorWhitelisted(ctx context.Context, operatorAddress string) bool {
    var found bool
    _ = k.Validator.Walk(ctx, nil, func(key string, val types.Validator) (stop bool, err error) {
        if val.OperatorAddress == operatorAddress {
            found = true
            return true, nil // stop iteration when found
        }
        return false, nil // continue iteration
    })
    return found
}
```

### Ante Decorator

```go
// ante/validator_whitelist.go
type ValidatorWhitelistDecorator struct {
    validatorRegistryKeeper validatorregistrykeeper.Keeper
}

func (vwd ValidatorWhitelistDecorator) AnteHandle(
    ctx sdk.Context, 
    tx sdk.Tx, 
    simulate bool, 
    next sdk.AnteHandler,
) (sdk.Context, error) {
    for _, msg := range tx.GetMsgs() {
        if createValMsg, ok := msg.(*stakingtypes.MsgCreateValidator); ok {
            if !vwd.validatorRegistryKeeper.IsValidatorWhitelisted(ctx, createValMsg.ValidatorAddress) {
                return ctx, errors.Wrapf(
                    sdkerrors.ErrUnauthorized, 
                    "validator address %s is not whitelisted",
                    createValMsg.ValidatorAddress,
                )
            }
        }
    }
    return next(ctx, tx, simulate)
}
```

### Dependency Injection

```go
// app/app.go
anteHandler, err := appante.NewAnteHandler(
    app.AuthKeeper,
    app.BankKeeper,
    app.txConfig.SignModeHandler(),
    ante.DefaultSigVerificationGasConsumer,
    app.ValidatorregistryKeeper, // ← Injected keeper
)
```

## Troubleshooting

### Error: "validator address is not whitelisted"

**Solution:** Verify the validator is in the whitelist:

```bash
# List all whitelisted validators
veranatestd query validatorregistry list-validator

# Get your validator operator address
veranatestd keys show <key-name> --bech32=val --keyring-backend test -a

# Compare addresses - they must match exactly
```

If not whitelisted, add via transaction or genesis.

### Error: Chain fails to start after setup

**Cause:** Module initialization order issue.

**Solution:** Verify `validatorregistry` initializes before `genutil` in `app/app_config.go`:

```go
{
    ModuleName: validatorregistrymoduletypes.ModuleName,
},
{
    ModuleName: genutiltypes.ModuleName,
},
```

### Error: Build errors after implementation

**Solution:**

```bash
go mod tidy
go build ./...
```

## Security Considerations

1. **Read-Only Access** - Ante handler only reads from KV store, never writes
2. **No State Modification** - Checking whitelist doesn't modify blockchain state
3. **Gas Efficiency** - Walk operation stops early when match is found
4. **Authority Control** - Only authorized accounts can modify whitelist

## Future Enhancements

1. **Caching** - Cache whitelist in memory for better performance
2. **Indexing** - Use indexed map for O(1) lookup instead of O(n) walk
3. **Term Management** - Implement automatic validator expiration using `term_end` field
4. **Council Voting** - Integrate with `x/group` for council-based onboarding decisions

## FAQ

**Q: Can I use my existing validator?**  
A: Yes! Just add its operator address to the genesis whitelist.

**Q: How do I get the operator address?**  
A: `veranatestd keys show <key-name> --bech32=val --keyring-backend test -a`

**Q: Can I add validators after the chain starts?**  
A: Yes, use the `onboard-validator` transaction (requires proper authority).

**Q: What happens if I try to create a non-whitelisted validator?**  
A: The transaction fails with error: "validator address is not whitelisted"

**Q: Is this the standard Cosmos SDK way?**  
A: Yes! This follows the same keeper injection pattern used by `auth` and `bank` modules.

## Summary

✅ Validator whitelist enforced at all block heights  
✅ Dynamic configuration via KV store  
✅ Follows Cosmos SDK best practices  
✅ Configurable via genesis or transactions  
✅ No hardcoded addresses  
✅ Supports council-based governance model  

The implementation is complete and production-ready! 🎉

