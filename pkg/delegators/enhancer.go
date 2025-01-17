package delegators

import (
	"fmt"
	"strings"

	"github.com/teambenny/goetl"
	"github.com/teambenny/goetl/etldata"

	"cosmossdk.io/log"

	"github.com/cosmos/cosmos-sdk/types/bech32"
)

type addressEnhancer struct {
	prefixes []string
	logger   log.Logger
	keys     []string
	name     string
}

// NewAddressEnhancer returns a new processor that enrich the data with addresses with the given prefixes.
func NewAddressEnhancer(prefixes []string, logger log.Logger) (goetl.Processor, error) {
	keys := make([]string, len(prefixes))
	for i, prefix := range prefixes {
		keys[i] = fmt.Sprintf("delegator-%s-address", prefix)
	}

	return &addressEnhancer{
		prefixes: prefixes,
		logger:   logger,
		keys:     keys,
		name:     fmt.Sprintf("AddressEnhancer(%s)", strings.Join(prefixes, ",")),
	}, nil
}

func (p *addressEnhancer) ProcessData(payload etldata.Payload, outputChan chan etldata.Payload, killChan chan error) {
	delegation := Delegation{}
	if err := payload.Parse(&delegation); err != nil {
		p.logger.Error(err.Error())
		killChan <- err
		return
	}
	data, err := payload.Objects()
	if err != nil {
		p.logger.Error(err.Error())
		killChan <- err
		return
	}

	for _, datum := range data {
		for idx, prefix := range p.prefixes {
			addr, err := convertAndEncodeMust(prefix, delegation.DelegatorNativeAddr)
			if err != nil {
				p.logger.Error(err.Error())
				killChan <- err
				return
			}
			datum[p.keys[idx]] = addr
		}
	}

	json, err := etldata.NewJSON(data)
	if err != nil {
		p.logger.Error(err.Error())
		killChan <- err
		return
	}

	outputChan <- json
}

func (p *addressEnhancer) Finish(_ chan etldata.Payload, _ chan error) {
}

func (p *addressEnhancer) String() string {
	return p.name
}

func convertAndEncodeMust(hrp string, bech string) (string, error) {
	_, bytes, err := bech32.DecodeAndConvert(bech)
	if err != nil {
		return "", err
	}

	return bech32.ConvertAndEncode(hrp, bytes)
}
