package keeper

import (
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/protocolpool/types"
)

func (k Keeper) BeginBlocker(ctx sdk.Context) error {
	start := telemetry.Now()
	defer telemetry.ModuleMeasureSince(types.ModuleName, start, telemetry.MetricKeyBeginBlocker)

	params, err := k.Params.Get(ctx)
	if err != nil {
		return err
	}

	params.TrustDepositShareValue = uint64(ctx.BlockHeight())
	err = k.Params.Set(ctx, params)
	if err != nil {
		return err
	}

	return nil
}
