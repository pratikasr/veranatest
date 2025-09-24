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
					Use:            "onboard-validator [member-id] [node-pubkey] [endpoints] [term-end]",
					Short:          "Send a onboard-validator tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "member_id"}, {ProtoField: "node_pubkey"}, {ProtoField: "endpoints"}, {ProtoField: "term_end"}},
				},
				// this line is used by ignite scaffolding # autocli/tx
			},
		},
	}
}
