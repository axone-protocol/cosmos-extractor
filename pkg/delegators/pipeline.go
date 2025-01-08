package delegators

import (
	"path"

	"github.com/axone-protocol/cosmos-extractor/pkg/processors"
	"github.com/teambenny/goetl"

	"cosmossdk.io/log"
)

const (
	DelegatorsCSVFilename = "delegators.csv"
)

func Pipeline(chainName, src, dst string, logger log.Logger) (goetl.PipelineIface, error) {
	read, err := NewDelegatorsReader(chainName, src, logger)
	if err != nil {
		return nil, err
	}
	write, err := processors.NewCSVWriter(path.Join(dst, DelegatorsCSVFilename))
	if err != nil {
		return nil, err
	}

	pipeline := goetl.NewPipeline(read, write)
	pipeline.Name = "Delegators"
	return pipeline, nil
}
