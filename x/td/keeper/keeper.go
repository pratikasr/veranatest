package keeper

import (
	"context"
	"cosmossdk.io/math"
	"errors"
	"fmt"

	"cosmossdk.io/collections"
	"cosmossdk.io/core/address"
	corestore "cosmossdk.io/core/store"
	"github.com/cosmos/cosmos-sdk/codec"

	"veranatest/x/td/types"
)

type Keeper struct {
	storeService corestore.KVStoreService
	cdc          codec.Codec
	addressCodec address.Codec
	// Address capable of executing a MsgUpdateParams message.
	// Typically, this should be the x/gov module account.
	authority []byte

	Schema        collections.Schema
	Params        collections.Item[types.Params]
	DustAmount    collections.Item[types.DustAmount]
	bankKeeper    types.BankKeeper
	accountKeeper types.AccountKeeper
}

func NewKeeper(
	storeService corestore.KVStoreService,
	cdc codec.Codec,
	addressCodec address.Codec,
	authority []byte,
	bankKeeper types.BankKeeper,
	accountKeeper types.AccountKeeper,

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

		Params:        collections.NewItem(sb, types.ParamsKey, "params", codec.CollValue[types.Params](cdc)),
		bankKeeper:    bankKeeper,
		DustAmount:    collections.NewItem(sb, types.DustAmountKey, "dust_amount", codec.CollValue[types.DustAmount](cdc)),
		accountKeeper: accountKeeper,
	}

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

// GetDustAmount retrieves the current dust amount as LegacyDec
func (k Keeper) GetDustAmount(ctx context.Context) (math.LegacyDec, error) {
	dustAmount, err := k.DustAmount.Get(ctx)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			// Return zero if not found (first time)
			return math.LegacyZeroDec(), nil
		}
		return math.LegacyDec{}, err
	}

	return dustAmount.Dust, nil
}

// SetDustAmount sets the dust amount from LegacyDec
func (k Keeper) SetDustAmount(ctx context.Context, amount math.LegacyDec) error {
	return k.DustAmount.Set(ctx, types.DustAmount{Dust: amount})
}

// AddToDustAmount adds an amount to the existing dust amount
func (k Keeper) AddToDustAmount(ctx context.Context, amount math.LegacyDec) error {
	currentDust, err := k.GetDustAmount(ctx)
	if err != nil {
		return err
	}

	newDust := currentDust.Add(amount)
	return k.SetDustAmount(ctx, newDust)
}

// ResetDustAmount resets the dust amount to zero
func (k Keeper) ResetDustAmount(ctx context.Context) error {
	return k.SetDustAmount(ctx, math.LegacyZeroDec())
}
