package ante

import (
	validatorregistrykeeper "veranatest/x/validatorregistry/keeper"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// ValidatorWhitelistDecorator checks if a validator address is whitelisted before allowing validator creation
type ValidatorWhitelistDecorator struct {
	validatorRegistryKeeper validatorregistrykeeper.Keeper
}

// NewValidatorWhitelistDecorator creates a new ValidatorWhitelistDecorator
func NewValidatorWhitelistDecorator(validatorRegistryKeeper validatorregistrykeeper.Keeper) ValidatorWhitelistDecorator {
	return ValidatorWhitelistDecorator{
		validatorRegistryKeeper: validatorRegistryKeeper,
	}
}

// AnteHandle checks if the validator creating a validator is whitelisted
// This check runs at ALL block heights including genesis (block height 0)
// because validatorregistry.InitGenesis runs BEFORE genutil.InitGenesis
func (vwd ValidatorWhitelistDecorator) AnteHandle(
	ctx sdk.Context,
	tx sdk.Tx,
	simulate bool,
	next sdk.AnteHandler,
) (sdk.Context, error) {
	for _, msg := range tx.GetMsgs() {
		if createValMsg, ok := msg.(*stakingtypes.MsgCreateValidator); ok {
			// Check if the validator address is whitelisted in the validatorregistry module
			if !vwd.validatorRegistryKeeper.IsValidatorWhitelisted(ctx, createValMsg.ValidatorAddress) {
				return ctx, errors.Wrapf(
					sdkerrors.ErrUnauthorized,
					"validator address %s is not whitelisted. Only whitelisted validators can create validators",
					createValMsg.ValidatorAddress,
				)
			}
		}
	}

	return next(ctx, tx, simulate)
}
