package validatorregistry

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"

	"veranatest/x/validatorregistry/types"
)

// AutoCLIOptions implements the autocli.HasAutoCLIConfig interface.
func (am AppModule) AutoCLIOptions() *autocliv1.ModuleOptions {
	return &autocliv1.ModuleOptions{
		Query: &autocliv1.ServiceCommandDescriptor{
			Service: types.Query_serviceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "Params",
					Use:       "params",
					Short:     "Shows the parameters of the module",
				},
				{
					RpcMethod: "ListValidator",
					Use:       "list-validator",
					Short:     "List all validator",
				},
				{
					RpcMethod:      "GetValidator",
					Use:            "get-validator [id]",
					Short:          "Gets a validator",
					Alias:          []string{"show-validator"},
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "index"}},
				},
				// this line is used by ignite scaffolding # autocli/query
			},
		},
		Tx: &autocliv1.ServiceCommandDescriptor{
			Service:              types.Msg_serviceDesc.ServiceName,
			EnhanceCustomCommand: true, // only required if you want to use the custom command
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "UpdateParams",
					Skip:      true, // skipped because authority gated
				},
				{
					RpcMethod:      "OnboardValidator",
					Use:            "onboard-validator [index] [member-id] [operator-address] [consensus-pubkey] [status] [term-end]",
					Short:          "Send a onboard-validator tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "index"}, {ProtoField: "member_id"}, {ProtoField: "operator_address"}, {ProtoField: "consensus_pubkey"}, {ProtoField: "status"}, {ProtoField: "term_end"}},
				},
				// this line is used by ignite scaffolding # autocli/tx
			},
		},
	}
}
