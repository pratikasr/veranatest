package keeper

import (
	"context"

	"veranatest/x/validatorregistry/types"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) OnboardValidator(ctx context.Context, msg *types.MsgOnboardValidator) (*types.MsgOnboardValidatorResponse, error) {
	// Validate creator address (should be the group policy or authorized address)
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(err, "invalid authority address")
	}

	// Validate required fields
	if msg.Index == "" {
		return nil, errorsmod.Wrap(types.ErrInvalidValidator, "index cannot be empty")
	}
	if msg.MemberId == "" {
		return nil, errorsmod.Wrap(types.ErrInvalidValidator, "member_id cannot be empty")
	}
	if msg.OperatorAddress == "" {
		return nil, errorsmod.Wrap(types.ErrInvalidValidator, "operator_address cannot be empty")
	}
	if msg.Status == "" {
		return nil, errorsmod.Wrap(types.ErrInvalidValidator, "status cannot be empty")
	}

	// Check if validator with this index already exists
	exists, err := k.Validator.Has(ctx, msg.Index)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to check validator existence")
	}
	if exists {
		return nil, errorsmod.Wrapf(types.ErrInvalidValidator, "validator with index %s already exists", msg.Index)
	}

	// Validate operator address format (should be cosmosvaloper...)
	_, err = sdk.ValAddressFromBech32(msg.OperatorAddress)
	if err != nil {
		return nil, errorsmod.Wrapf(types.ErrInvalidValidator, "invalid operator address format: %v", err)
	}

	// Create the Validator object
	validator := types.Validator{
		Index:           msg.Index,
		MemberId:        msg.MemberId,
		OperatorAddress: msg.OperatorAddress,
		ConsensusPubkey: msg.ConsensusPubkey,
		Status:          msg.Status,
		TermEnd:         msg.TermEnd,
	}

	// Store the validator in the KV store
	if err := k.Validator.Set(ctx, msg.Index, validator); err != nil {
		return nil, errorsmod.Wrap(err, "failed to store validator")
	}

	// Emit event for validator onboarding
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			"validator_onboarded",
			sdk.NewAttribute("index", msg.Index),
			sdk.NewAttribute("member_id", msg.MemberId),
			sdk.NewAttribute("operator_address", msg.OperatorAddress),
			sdk.NewAttribute("status", msg.Status),
		),
	)

	return &types.MsgOnboardValidatorResponse{}, nil
}
