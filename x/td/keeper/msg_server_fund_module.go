package keeper

import (
	"context"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"veranatest/x/td/types"

	errorsmod "cosmossdk.io/errors"
)

func (k msgServer) FundModule(ctx context.Context, msg *types.MsgFundModule) (*types.MsgFundModuleResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(err, "invalid authority address")
	}
	senderAcc, _ := sdk.AccAddressFromBech32(msg.Creator)

	err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, senderAcc, msg.Module, sdk.NewCoins(sdk.NewCoin("uvna", sdkmath.NewInt(msg.Amount))))
	if err != nil {
		return nil, err
	}
	if msg.Module == types.ModuleName {

		params, _ := k.Params.Get(ctx)
		params.TrustDepositValue = params.TrustDepositValue + uint64(msg.Amount)
		err = k.Params.Set(ctx, params)
		if err != nil {
			return nil, err
		}
	}

	return &types.MsgFundModuleResponse{}, nil
}
