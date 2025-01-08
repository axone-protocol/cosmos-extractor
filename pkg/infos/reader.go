package infos

import (
	"fmt"
	"io"

	"github.com/axone-protocol/cosmos-extractor/pkg/keeper"
	"github.com/teambenny/goetl"
	"github.com/teambenny/goetl/etldata"

	"cosmossdk.io/log"
)

type infoReader struct {
	chainName string
	src       string
	logger    log.Logger
	closer    io.Closer
}

// NewInfoReader returns a new Reader that reads metadata information about a blockchain data store.
func NewInfoReader(chainName, src string, logger log.Logger) (goetl.Processor, error) {
	return &infoReader{
		chainName: chainName,
		src:       src,
		logger:    logger,
	}, nil
}

func (r *infoReader) ProcessData(_ etldata.Payload, outputChan chan etldata.Payload, killChan chan error) {
	keepers, err := keeper.OpenStore(r.src, r.logger)
	if err != nil {
		r.logger.Error(err.Error())
		killChan <- err
		return
	}
	r.closer = keepers

	payload := Info{
		Name:         r.chainName,
		StoreVersion: fmt.Sprintf("%d", keepers.Store.LastCommitID().Version),
		StoreHash:    fmt.Sprintf("%X", keepers.Store.LastCommitID().Hash),
	}

	json, err := etldata.NewJSON(payload)
	if err != nil {
		r.logger.Error(err.Error())
		killChan <- err
		return
	}

	outputChan <- json
}

func (r *infoReader) Finish(_ chan etldata.Payload, killChan chan error) {
	if r.closer != nil {
		err := r.closer.Close()
		if err != nil {
			r.logger.Error(err.Error())
			killChan <- err
		}
	}
}

func (r *infoReader) String() string {
	return "InfoReader"
}
