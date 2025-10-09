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

	// Generate index from member_id (use member_id as unique identifier)
	index := msg.MemberId

	// Check if validator with this index already exists
	exists, err := k.Validator.Has(ctx, index)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to check validator existence")
	}
	if exists {
		return nil, errorsmod.Wrapf(types.ErrInvalidValidator, "validator with index %s already exists", index)
	}

	// WORKAROUND: Since operator_address field is missing from proto,
	// we're temporarily using the 'endpoints' field to pass the operator address.
	// In the future, update the proto to add operator_address as a proper field.
	operatorAddress := msg.Endpoints

	// Validate operator address format (should be cosmosvaloper...)
	if operatorAddress == "" {
		return nil, errorsmod.Wrap(types.ErrInvalidValidator, "operator address (passed via endpoints) cannot be empty")
	}

	// Optionally validate it's a valid bech32 address with valoper prefix
	_, err = sdk.ValAddressFromBech32(operatorAddress)
	if err != nil {
		return nil, errorsmod.Wrapf(types.ErrInvalidValidator, "invalid operator address format: %v", err)
	}

	// Create the Validator object
	validator := types.Validator{
		Index:           index,
		MemberId:        msg.MemberId,
		OperatorAddress: operatorAddress,
		ConsensusPubkey: msg.NodePubkey,
		Status:          "active", // Hardcoded for now, should be added to proto
		TermEnd:         msg.TermEnd,
	}

	// Store the validator in the KV store
	if err := k.Validator.Set(ctx, index, validator); err != nil {
		return nil, errorsmod.Wrap(err, "failed to store validator")
	}

	// Emit event for validator onboarding
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			"validator_onboarded",
			sdk.NewAttribute("index", index),
			sdk.NewAttribute("member_id", msg.MemberId),
			sdk.NewAttribute("operator_address", operatorAddress),
		),
	)

	return &types.MsgOnboardValidatorResponse{}, nil
}
