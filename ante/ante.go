package ante

import (
	"fmt"
	"log"

	txsigning "cosmossdk.io/x/tx/signing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
)

// NewAnteHandler returns an AnteHandler with debugging
func NewAnteHandler(
	accountKeeper authkeeper.AccountKeeper,
	bankKeeper bankkeeper.Keeper,
	signModeHandler *txsigning.HandlerMap,
	sigGasConsumer ante.SignatureVerificationGasConsumer,
) (sdk.AnteHandler, error) {

	if bankKeeper == nil {
		return nil, fmt.Errorf("bank keeper is required")
	}
	if signModeHandler == nil {
		return nil, fmt.Errorf("sign mode handler is required")
	}

	// Start with just basic decorators
	anteDecorators := []sdk.AnteDecorator{
		ante.NewSetUpContextDecorator(),
		ante.NewValidateBasicDecorator(),

		// ADD OUR CUSTOM VALIDATOR WHITELIST DECORATOR HERE
		NewValidatorWhitelistDecorator(),

		// Try with minimal decorators first
		ante.NewDeductFeeDecorator(accountKeeper, bankKeeper, nil, nil),
		ante.NewIncrementSequenceDecorator(accountKeeper),
	}

	log.Printf("DEBUG: Successfully created ante decorators")
	return sdk.ChainAnteDecorators(anteDecorators...), nil
}
