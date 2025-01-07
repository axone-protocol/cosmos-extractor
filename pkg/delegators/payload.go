package delegators

type Delegation struct {
	ChainName           string `json:"chain_name"`
	DelegatorNativeAddr string `json:"delegator_native_addr"`
	DelegatorCosmosAddr string `json:"delegator_cosmos_addr"`
	DelegatorAxoneAddr  string `json:"delegator_axone_addr"`

	ValidatorAddr string `json:"validator_addr"`
	Shares        string `json:"shares"`
}
