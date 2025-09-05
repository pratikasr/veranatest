package types

import "cosmossdk.io/math"

const DefaultTrustDepositValue uint64 = 0

const (
	DefaultTrustDepositShareValue = "1.0"
	DefaultTrustDepositYieldRate  = "0.15"
)

// NewParams creates a new Params instance.
func NewParams(trustDepositShareValue math.LegacyDec, trust_deposit_value uint64, trust_deposit_yield_rate math.LegacyDec) Params {
	return Params{TrustDepositShareValue: trustDepositShareValue, TrustDepositValue: trust_deposit_value, TrustDepositYieldRate: trust_deposit_yield_rate}
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	TrustDepositShareValue, _ := math.LegacyNewDecFromStr(DefaultTrustDepositShareValue)
	TrustDepositYieldRate, _ := math.LegacyNewDecFromStr(DefaultTrustDepositYieldRate)
	return NewParams(TrustDepositShareValue, DefaultTrustDepositValue, TrustDepositYieldRate)
}

// Validate validates the set of params.
func (p Params) Validate() error {
	if err := validateTrustDepositShareValue(p.TrustDepositShareValue); err != nil {
		return err
	}
	if err := validateTrustDepositValue(p.TrustDepositValue); err != nil {
		return err
	}
	if err := validateTrustDepositYieldRate(p.TrustDepositYieldRate); err != nil {
		return err
	}

	return nil
}
func validateTrustDepositShareValue(v math.LegacyDec) error {

	return nil
}
func validateTrustDepositValue(v uint64) error {

	return nil
}
func validateTrustDepositYieldRate(v math.LegacyDec) error {

	return nil
}
