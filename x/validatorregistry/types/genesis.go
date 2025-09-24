package types

import "fmt"

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params:       DefaultParams(),
		ValidatorMap: []Validator{}}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	validatorIndexMap := make(map[string]struct{})

	for _, elem := range gs.ValidatorMap {
		index := fmt.Sprint(elem.Index)
		if _, ok := validatorIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for validator")
		}
		validatorIndexMap[index] = struct{}{}
	}

	return gs.Params.Validate()
}
