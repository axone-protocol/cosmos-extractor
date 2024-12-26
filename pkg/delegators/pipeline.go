package delegators

import (
	"path"

	"github.com/teambenny/goetl"
)

const (
	DelegatorsCSVFilename = "delegators.csv"
	ChainsCSVFilename     = "chains.csv"
)

type pipelines struct {
	pipelines []goetl.PipelineIface
}

func Pipeline(chainName, src, dst string) (goetl.PipelineIface, error) {
	readDelegators, err := NewDelegatorsReader(chainName, src)
	if err != nil {
		return nil, err
	}
	writeDelegators, err := NewCSVWriter(path.Join(dst, DelegatorsCSVFilename))
	if err != nil {
		return nil, err
	}

	readChain, err := NewChainReader(chainName, src)
	if err != nil {
		return nil, err
	}
	writeChain, err := NewCSVWriter(path.Join(dst, ChainsCSVFilename))
	if err != nil {
		return nil, err
	}

	return &pipelines{
		pipelines: []goetl.PipelineIface{
			goetl.NewPipeline(readChain, writeChain),
			goetl.NewPipeline(readDelegators, writeDelegators),
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
