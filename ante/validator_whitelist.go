package ante

import (
	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

const ALLOWED_VALIDATOR_ADDRESS = "cosmosvaloper1rkz2eeu3rveg7u6srnsdkcjqmwc32kyl9565pm"

type ValidatorWhitelistDecorator struct{}

func NewValidatorWhitelistDecorator() ValidatorWhitelistDecorator {
	return ValidatorWhitelistDecorator{}
}

func (vwd ValidatorWhitelistDecorator) AnteHandle(
	ctx sdk.Context,
	tx sdk.Tx,
	simulate bool,
	next sdk.AnteHandler,
) (sdk.Context, error) {
	for _, msg := range tx.GetMsgs() {
		if createValMsg, ok := msg.(*stakingtypes.MsgCreateValidator); ok {
			if createValMsg.ValidatorAddress != ALLOWED_VALIDATOR_ADDRESS {
				return ctx, errors.Wrapf(
					sdkerrors.ErrUnauthorized,
					"validator address %s is not authorized. Only %s can create validators",
					createValMsg.ValidatorAddress,
					ALLOWED_VALIDATOR_ADDRESS,
				)
			}
		}
	}

	return next(ctx, tx, simulate)
}
