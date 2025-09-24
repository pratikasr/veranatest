package keeper

import (
	"context"

	"veranatest/x/validatorregistry/types"

	errorsmod "cosmossdk.io/errors"
)

func (k msgServer) OnboardValidator(ctx context.Context, msg *types.MsgOnboardValidator) (*types.MsgOnboardValidatorResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(err, "invalid authority address")
	}

	// TODO: Handle the message

	return &types.MsgOnboardValidatorResponse{}, nil
}
