package delegators

import (
	"path"

	"github.com/teambenny/goetl"

	"cosmossdk.io/log"
)

const (
	DelegatorsCSVFilename = "delegators.csv"
	ChainsCSVFilename     = "chains.csv"
)

type pipelines struct {
	pipelines []goetl.PipelineIface
}

func Pipeline(chainName, src, dst string, logger log.Logger) (goetl.PipelineIface, error) {
	readDelegators, err := NewDelegatorsReader(chainName, src, logger)
	if err != nil {
		return nil, err
	}
	writeDelegators, err := NewCSVWriter(path.Join(dst, DelegatorsCSVFilename))
	if err != nil {
		return nil, err
	}

	readChain, err := NewChainReader(chainName, src, logger)
	if err != nil {
		return nil, err
	}
	writeChain, err := NewCSVWriter(path.Join(dst, ChainsCSVFilename))
	if err != nil {
		return nil, err
	}

	return &pipelines{
		pipelines: []goetl.PipelineIface{
			func() goetl.PipelineIface {
				pipeline := goetl.NewPipeline(readChain, writeChain)
				pipeline.Name = "Chain"
				return pipeline
			}(),
			func() goetl.PipelineIface {
				pipeline := goetl.NewPipeline(readDelegators, writeDelegators)
				pipeline.Name = "Delegators"
				return pipeline
			}(),
		},
	}, nil
}

func (p *pipelines) Run() chan error {
	errChan := make(chan error)

	go func() {
		defer close(errChan)

		for _, pipeline := range p.pipelines {
			c := pipeline.Run()
			err := <-c
			if err != nil {
				errChan <- err
			}
		}
	}()

	return errChan
}
