package types

const DefaultTrustDepositShareValue uint64 = 1

// NewParams creates a new Params instance.
func NewParams(trustDepositShareValue uint64) Params {
	return Params{TrustDepositShareValue: trustDepositShareValue}
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return NewParams(DefaultTrustDepositShareValue)
}

// Validate validates the set of params.
func (p Params) Validate() error {
	if err := validateTrustDepositShareValue(p.TrustDepositShareValue); err != nil {
		return err
	}

	return nil
}
func validateTrustDepositShareValue(v uint64) error {

	return nil
}
