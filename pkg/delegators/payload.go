package delegators

type Delegation struct {
	ChainName           string `json:"chain_name"`
	DelegatorNativeAddr string `json:"delegator_native_addr"`
	ValidatorAddr       string `json:"validator_addr"`
	Shares              string `json:"shares"`
}
