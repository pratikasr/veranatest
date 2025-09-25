package keeper

//import (
//	"fmt"
//	sdk "github.com/cosmos/cosmos-sdk/types"
//)
//
//// Hooks implements the staking hooks
//type Hooks struct {
//	k Keeper
//}
//
//func (h Hooks) BeforeValidatorCreated(ctx sdk.Context, valAddr sdk.ValAddress) error {
//	allowedAddress := "cosmos1rkz2eeu3rveg7u6srnsdkcjqmwc32kylqqwpdg"
//   //veranatestd keys add val --keyring-backend test --recover
//	// seed = violin harvest scrap alter economy sheriff pen narrow mule gold gallery hollow dust dry near bullet volcano bean bar lab wagon scorpion antenna seven
//
//	// Convert validator address to account address
//	operatorAddr := sdk.AccAddress(valAddr)
//
//	// Simple check
//	if operatorAddr.String() != allowedAddress {
//		h.k.Logger().Error(
//			"HOOK BLOCKED: Unauthorized validator creation attempt",
//			"attempted_by", operatorAddr.String(),
//			"allowed_only", allowedAddress,
//		)
//		return fmt.Errorf("validator creation not authorized for address %s. Only %s is allowed",
//			operatorAddr.String(), allowedAddress)
//	}
//
//	h.k.Logger().Info(
//		"HOOK ALLOWED: Authorized validator creation",
//		"operator", operatorAddr.String(),
//	)
//
//	return nil
//}
