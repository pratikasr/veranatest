package types

const DefaultTrustDepositShareValue uint64 = 1
const DefaultTrustDepositValue uint64 = 0

// NewParams creates a new Params instance.
func NewParams(trustDepositShareValue uint64, trust_deposit_value uint64) Params {
	return Params{TrustDepositShareValue: trustDepositShareValue, TrustDepositValue: trust_deposit_value}
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return NewParams(DefaultTrustDepositShareValue, DefaultTrustDepositValue)
}

// Validate validates the set of params.
func (p Params) Validate() error {
	if err := validateTrustDepositShareValue(p.TrustDepositShareValue); err != nil {
		return err
	}
	if err := validateTrustDepositValue(p.TrustDepositValue); err != nil {
		return err
	}

	return nil
}
func validateTrustDepositShareValue(v uint64) error {

	return nil
}
func validateTrustDepositValue(v uint64) error {

	return nil
}
