package delegators


type Delegations struct {
	ChainName           string `csv:"chain_name"`
	DelegatorNativeAddr string `csv:"delegator_native_addr"`
	DelegatorCosmosAddr string `csv:"delegator_cosmos_addr"`
	DelegatorAxoneAddr  string `csv:"delegator_axone_addr"`

	ValidatorAddr string `csv:"validator_addr"`
	Shares        string `csv:"shares"`
}
