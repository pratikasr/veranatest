package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) BeginBlocker(ctx sdk.Context) error {
	// for sending fund from verana-pool to the trust deposit module
	// and sending excess funds back to the protocolpool i.e community pool

	return nil
}
func (k Keeper) SendFundsFromVeranaPool(ctx sdk.Context) error {
	return nil
}

func (k Keeper) SendFundsBackToCommunityPool(ctx sdk.Context) error {
	return nil
}
