# Validator Whitelist Implementation

## Overview

This implementation uses the **standard Cosmos SDK pattern** of injecting module keepers into the ante handler to access the validator whitelist from the KV store.

## Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                         Application                              │
│  ┌────────────┐        ┌──────────────────┐                    │
│  │  App.go    │───────▶│ AnteHandler      │                    │
│  └────────────┘        └──────────────────┘                    │
│        │                        │                                │
│        │ Injects Keepers       │ Uses Keepers                  │
│        ▼                        ▼                                │
│  ┌─────────────────────────────────────────┐                   │
│  │  ValidatorWhitelistDecorator            │                   │
│  │  - Has: validatorRegistryKeeper         │                   │
│  │  - Checks: IsValidatorWhitelisted()     │                   │
│  └─────────────────────────────────────────┘                   │
│                        │                                         │
│                        │ Calls Method                           │
│                        ▼                                         │
│  ┌─────────────────────────────────────────┐                   │
│  │  ValidatorRegistry Keeper               │                   │
│  │  - Method: IsValidatorWhitelisted()     │                   │
│  │  - Accesses: KV Store (Validator Map)   │                   │
│  └─────────────────────────────────────────┘                   │
│                        │                                         │
│                        │ Reads From                             │
│                        ▼                                         │
│  ┌─────────────────────────────────────────┐                   │
│  │  KV Store                               │                   │
│  │  Key: "validator/value/{index}"         │                   │
│  │  Value: Validator{OperatorAddress, ...} │                   │
│  └─────────────────────────────────────────┘                   │
└─────────────────────────────────────────────────────────────────┘
```

## How It Works

### 1. **KV Store (Data Layer)**
- The `validatorregistry` module stores whitelisted validators in a KV store
- Each validator has an `operator_address` field (cosmosvaloper...)
- Stored using Cosmos SDK collections: `collections.Map[string, types.Validator]`

### 2. **Keeper (Business Logic Layer)**
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

### 3. **Ante Decorator (Authorization Layer)**
```go
// ante/validator_whitelist.go
type ValidatorWhitelistDecorator struct {
    validatorRegistryKeeper validatorregistrykeeper.Keeper
}

func (vwd ValidatorWhitelistDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
    for _, msg := range tx.GetMsgs() {
        if createValMsg, ok := msg.(*stakingtypes.MsgCreateValidator); ok {
            if !vwd.validatorRegistryKeeper.IsValidatorWhitelisted(ctx, createValMsg.ValidatorAddress) {
                return ctx, errors.Wrapf(sdkerrors.ErrUnauthorized, "validator address %s is not whitelisted")
            }
        }
    }
    return next(ctx, tx, simulate)
}
```

### 4. **Dependency Injection (Setup Layer)**
```go
// app/app.go
anteHandler, err := appante.NewAnteHandler(
    app.AuthKeeper,
    app.BankKeeper,
    app.txConfig.SignModeHandler(),
    ante.DefaultSigVerificationGasConsumer,
    app.ValidatorregistryKeeper, // ← Injected here
)
```

## Why This Pattern?

This follows the **same pattern** used throughout Cosmos SDK:

| Module | Keeper | Used In Ante For |
|--------|--------|------------------|
| `auth` | `accountKeeper` | Signature verification |
| `bank` | `bankKeeper` | Fee deduction |
| **`validatorregistry`** | **`validatorRegistryKeeper`** | **Whitelist checking** |

### Benefits:
1. ✅ **Modularity**: Keeper encapsulates KV store logic
2. ✅ **Testability**: Can mock keeper in tests
3. ✅ **Standard Pattern**: Follows Cosmos SDK conventions
4. ✅ **Type Safety**: Compile-time checks
5. ✅ **Maintainability**: Clear separation of concerns

## Managing the Whitelist

### Adding Validators to Whitelist

You can add validators in two ways:

#### Option 1: Via Genesis File

Edit your genesis file (`~/.veranatest/config/genesis.json`):

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

#### Option 2: Via Transaction (After Chain Start)

Use the `onboard-validator` message:

```bash
veranatestd tx validatorregistry onboard-validator \
  validator3 \
  member003 \
  cosmosvaloper1xyz... \
  "" \
  active \
  0 \
  --from=authority \
  --chain-id=vna-testnet-1
```

### Querying the Whitelist

List all whitelisted validators:
```bash
veranatestd query validatorregistry list-validator
```

Check specific validator:
```bash
veranatestd query validatorregistry show-validator validator1
```

## Testing the Whitelist

### Test 1: Whitelisted Validator (Should Succeed)

```bash
# Get your validator operator address
VALIDATOR_ADDR=$(veranatestd keys show cooluser --bech32=val --keyring-backend test -a)

# Create validator
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

**Expected Result**: ✅ Success (if address is in whitelist)

### Test 2: Non-Whitelisted Validator (Should Fail)

```bash
# Create new key
veranatestd keys add unauthorized --keyring-backend test

# Try to create validator
veranatestd tx staking create-validator \
  --amount=1000000uvna \
  --pubkey=$(veranatestd tendermint show-validator) \
  --moniker="unauthorized-validator" \
  --commission-rate="0.10" \
  --commission-max-rate="0.20" \
  --commission-max-change-rate="0.01" \
  --min-self-delegation="1" \
  --from=unauthorized \
  --chain-id=vna-testnet-1 \
  --keyring-backend=test \
  --fees=500000uvna
```

**Expected Result**: ❌ Error: "validator address is not whitelisted"

## Code Files Modified

### 1. `x/validatorregistry/keeper/keeper.go`
**Added:**
- `IsValidatorWhitelisted(ctx context.Context, operatorAddress string) bool`
- Imports: `context`

**Purpose**: Provides method to check if validator address exists in KV store

### 2. `ante/validator_whitelist.go`
**Changed:**
- Removed: Hardcoded `ALLOWED_VALIDATOR_ADDRESS` constant
- Added: `validatorRegistryKeeper` field to decorator struct
- Updated: Constructor to accept keeper
- Updated: `AnteHandle` to use keeper method

**Purpose**: Uses keeper to check whitelist instead of hardcoded value

### 3. `ante/ante.go`
**Changed:**
- Added: `validatorregistrykeeper` import
- Added: `validatorRegistryKeeper` parameter to `NewAnteHandler`
- Updated: Passes keeper to `NewValidatorWhitelistDecorator`
- Added: Documentation comments

**Purpose**: Injects keeper dependency into decorator

### 4. `app/app.go`
**Changed:**
- Updated: `NewAnteHandler` call to pass `app.ValidatorregistryKeeper`
- Added: Documentation comment

**Purpose**: Provides keeper instance to ante handler

## Migration from Hardcoded Address

If you were using the hardcoded address approach:

### Before:
```go
const ALLOWED_VALIDATOR_ADDRESS = "cosmosvaloper1rkz2eeu3rveg7u6srnsdkcjqmwc32kyl9565pm"
```

### After:
Add to genesis or via transaction:
```json
{
  "index": "validator1",
  "operator_address": "cosmosvaloper1rkz2eeu3rveg7u6srnsdkcjqmwc32kyl9565pm",
  ...
}
```

## Troubleshooting

### Issue: "validator address is not whitelisted"

**Check if validator is in whitelist:**
```bash
veranatestd query validatorregistry list-validator
```

**Verify your validator address:**
```bash
veranatestd keys show <key-name> --bech32=val --keyring-backend test -a
```

**Add to whitelist:**
Either update genesis file or send onboard transaction.

### Issue: Build errors

**Verify imports:**
```bash
go mod tidy
go build ./...
```

### Issue: Keeper not accessible in ante

**Verify keeper is exported in app:**
```go
// app/app.go
ValidatorregistryKeeper validatorregistrymodulekeeper.Keeper
```

## Security Considerations

1. **Read-Only Access**: Ante handler only reads from KV store, never writes
2. **No State Modification**: Checking whitelist doesn't modify blockchain state
3. **Gas Efficiency**: Walk operation stops early when match is found
4. **Authority Control**: Only governance or authorized accounts can modify whitelist

## Future Enhancements

1. **Caching**: Cache whitelist in memory for better performance
2. **Indexing**: Use indexed map for O(1) lookup instead of O(n) walk
3. **Governance Integration**: Allow governance proposals to modify whitelist
4. **Expiration**: Use `term_end` field for time-based expiration

## Summary

This implementation demonstrates the **correct Cosmos SDK pattern** for accessing module state from ante handlers:

1. ✅ Module stores data in KV store
2. ✅ Keeper provides interface methods
3. ✅ Keeper is injected into ante handler
4. ✅ Decorator uses keeper to access state
5. ✅ Follows same pattern as `auth` and `bank` modules

The validator whitelist is now **dynamic, configurable, and follows Cosmos SDK best practices**.

