package keeper

import (
	"context"

	"fmt"
	"veranatest/x/validatorregistry/types"

	"cosmossdk.io/collections"
	"cosmossdk.io/core/address"
	corestore "cosmossdk.io/core/store"
	"cosmossdk.io/log"
	"github.com/cosmos/cosmos-sdk/codec"
)

type Keeper struct {
	storeService corestore.KVStoreService
	cdc          codec.Codec
	addressCodec address.Codec
	logger       log.Logger
	// Address capable of executing a MsgUpdateParams message.
	// Typically, this should be the x/gov module account.
	authority []byte

	Schema    collections.Schema
	Params    collections.Item[types.Params]
	Validator collections.Map[string, types.Validator]
}

func NewKeeper(
	storeService corestore.KVStoreService,
	cdc codec.Codec,
	addressCodec address.Codec,
	authority []byte,
	logger log.Logger,

) Keeper {
	if _, err := addressCodec.BytesToString(authority); err != nil {
		panic(fmt.Sprintf("invalid authority address %s: %s", authority, err))
	}

	sb := collections.NewSchemaBuilder(storeService)

	k := Keeper{
		storeService: storeService,
		cdc:          cdc,
		addressCodec: addressCodec,
		authority:    authority,
		logger:       logger,
		Params:       collections.NewItem(sb, types.ParamsKey, "params", codec.CollValue[types.Params](cdc)),
		Validator:    collections.NewMap(sb, types.ValidatorKey, "validator", collections.StringKey, codec.CollValue[types.Validator](cdc))}

	schema, err := sb.Build()
	if err != nil {
		panic(err)
	}
	k.Schema = schema

	return k
}

// GetAuthority returns the module's authority.
func (k Keeper) GetAuthority() []byte {
	return k.authority
}

func (k Keeper) Logger() log.Logger {
	return k.logger.With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// IsValidatorWhitelisted checks if a validator operator address is whitelisted
// This method is used by the ante handler to verify if a validator can create a validator
func (k Keeper) IsValidatorWhitelisted(ctx context.Context, operatorAddress string) bool {
	// Walk through all validators in the store and check if the operator address matches
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
