package infos

import (
	"path"

	"github.com/axone-protocol/cosmos-extractor/pkg/processors"
	"github.com/teambenny/goetl"

	"cosmossdk.io/log"
)

const (
	InfosCSVFilename = "infos.csv"
)

func Pipeline(chainName, src, dst string, logger log.Logger) (goetl.PipelineIface, error) {
	read, err := NewInfoReader(chainName, src, logger)
	if err != nil {
		return nil, err
	}
	write, err := processors.NewCSVWriter(path.Join(dst, InfosCSVFilename))
	if err != nil {
		return nil, err
	}

	pipeline := goetl.NewPipeline(read, write)
	pipeline.Name = "Chain"
	return pipeline, nil
}
