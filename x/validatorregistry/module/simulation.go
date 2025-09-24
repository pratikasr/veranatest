package validatorregistry

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	validatorregistrysimulation "veranatest/x/validatorregistry/simulation"
	"veranatest/x/validatorregistry/types"
)

// GenerateGenesisState creates a randomized GenState of the module.
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	validatorregistryGenesis := types.GenesisState{
		Params: types.DefaultParams(),
	}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&validatorregistryGenesis)
}

// RegisterStoreDecoder registers a decoder.
func (am AppModule) RegisterStoreDecoder(_ simtypes.StoreDecoderRegistry) {}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	operations := make([]simtypes.WeightedOperation, 0)
	const (
		opWeightMsgOnboardValidator          = "op_weight_msg_validatorregistry"
		defaultWeightMsgOnboardValidator int = 100
	)

	var weightMsgOnboardValidator int
	simState.AppParams.GetOrGenerate(opWeightMsgOnboardValidator, &weightMsgOnboardValidator, nil,
		func(_ *rand.Rand) {
			weightMsgOnboardValidator = defaultWeightMsgOnboardValidator
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgOnboardValidator,
		validatorregistrysimulation.SimulateMsgOnboardValidator(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))

	return operations
}

// ProposalMsgs returns msgs used for governance proposals for simulations.
func (am AppModule) ProposalMsgs(simState module.SimulationState) []simtypes.WeightedProposalMsg {
	return []simtypes.WeightedProposalMsg{}
}
