package ante

import (
	"fmt"
	"log"

	txsigning "cosmossdk.io/x/tx/signing"

	validatorregistrykeeper "veranatest/x/validatorregistry/keeper"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
)

// NewAnteHandler returns an AnteHandler that includes validator whitelist checking
// It follows the standard Cosmos SDK pattern of injecting module keepers as dependencies
func NewAnteHandler(
	accountKeeper authkeeper.AccountKeeper,
	bankKeeper bankkeeper.Keeper,
	signModeHandler *txsigning.HandlerMap,
	sigGasConsumer ante.SignatureVerificationGasConsumer,
	validatorRegistryKeeper validatorregistrykeeper.Keeper,
) (sdk.AnteHandler, error) {

	if bankKeeper == nil {
		return nil, fmt.Errorf("bank keeper is required")
	}
	if signModeHandler == nil {
		return nil, fmt.Errorf("sign mode handler is required")
	}

	// Chain decorators in order
	anteDecorators := []sdk.AnteDecorator{
		ante.NewSetUpContextDecorator(),
		ante.NewValidateBasicDecorator(),

		// Validator whitelist check - uses validatorregistry keeper to check KV store
		NewValidatorWhitelistDecorator(validatorRegistryKeeper),

		// Standard Cosmos SDK decorators
		ante.NewDeductFeeDecorator(accountKeeper, bankKeeper, nil, nil),
		ante.NewIncrementSequenceDecorator(accountKeeper),
	}

	log.Printf("DEBUG: Successfully created ante decorators with validator whitelist")
	return sdk.ChainAnteDecorators(anteDecorators...), nil
}
