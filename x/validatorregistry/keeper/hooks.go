package keeper

import (
	"context"
	"cosmossdk.io/math"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// Hooks implements the staking hooks
type Hooks struct {
	k Keeper
}

// Create hooks
func (k Keeper) Hooks() Hooks {
	return Hooks{k}
}

// BeforeValidatorModified - HARDCODED ADDRESS TEST
func (h Hooks) BeforeValidatorModified(ctx context.Context, valAddr sdk.ValAddress) error {
	// HARDCODE YOUR TEST ADDRESS HERE
	allowedAddress := "veranaXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"

	// Convert validator address to account address
	operatorAddr := sdk.AccAddress(valAddr)

	// Simple check
	if operatorAddr.String() != allowedAddress {
		h.k.Logger().Error(
			"HOOK BLOCKED: Unauthorized validator creation attempt",
			"attempted_by", operatorAddr.String(),
			"allowed_only", allowedAddress,
		)
		return fmt.Errorf("validator creation not authorized for address %s. Only %s is allowed",
			operatorAddr.String(), allowedAddress)
	}

	h.k.Logger().Info(
		"HOOK ALLOWED: Authorized validator creation",
		"operator", operatorAddr.String(),
	)

	return nil
}

// AfterValidatorCreated - Just log
func (h Hooks) AfterValidatorCreated(ctx context.Context, valAddr sdk.ValAddress) error {
	operatorAddr := sdk.AccAddress(valAddr)

	h.k.Logger().Info(
		"HOOK SUCCESS: Validator successfully created",
		"operator", operatorAddr.String(),
	)

	return nil
}

// ------------------- Empty Required Hooks -------------------
func (h Hooks) AfterValidatorBonded(ctx context.Context, consAddr sdk.ConsAddress, valAddr sdk.ValAddress) error {
	return nil
}
func (h Hooks) AfterValidatorBeginUnbonding(ctx context.Context, consAddr sdk.ConsAddress, valAddr sdk.ValAddress) error {
	return nil
}
func (h Hooks) AfterValidatorRemoved(ctx context.Context, consAddr sdk.ConsAddress, valAddr sdk.ValAddress) error {
	return nil
}
func (h Hooks) BeforeDelegationCreated(ctx context.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) error {
	return nil
}
func (h Hooks) BeforeDelegationSharesModified(ctx context.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) error {
	return nil
}
func (h Hooks) BeforeDelegationRemoved(ctx context.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) error {
	return nil
}
func (h Hooks) AfterDelegationModified(ctx context.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) error {
	return nil
}
func (h Hooks) AfterUnbondingInitiated(ctx context.Context, id uint64) error {
	return nil
}
func (h Hooks) BeforeValidatorSlashed(ctx context.Context, valAddr sdk.ValAddress, fraction math.LegacyDec) error {
	return nil
}

// Ensure interface compliance
var _ stakingtypes.StakingHooks = Hooks{}
