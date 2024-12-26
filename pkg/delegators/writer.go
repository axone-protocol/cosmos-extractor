package delegators

import (
	"bufio"
	"fmt"
	"os"

	"github.com/teambenny/goetl"
	"github.com/teambenny/goetl/etldata"
	"github.com/teambenny/goetl/processors"
)

type csvWriter struct {
	file          *os.File
	writer        *bufio.Writer
	processor     *processors.CSVWriter
	processorName string
}

func NewCSVWriter(dest string) (goetl.Processor, error) {
	file, err := os.OpenFile(dest, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return nil, err
	}
	writer := bufio.NewWriter(file)
	processor := processors.NewCSVWriter(writer)

	return &csvWriter{
		file:          file,
		writer:        writer,
		processor:     processor,
		processorName: fmt.Sprintf("CSVWriter<%s>", file.Name()),
	}, nil
}

func (w *csvWriter) ProcessData(d etldata.Payload, outputChan chan etldata.Payload, killChan chan error) {
	w.processor.ProcessData(d, outputChan, killChan)
}

func (w *csvWriter) Finish(_ chan etldata.Payload, _ chan error) {
	w.writer.Flush()
	w.file.Close()
}

func (w *csvWriter) String() string {
	return w.processorName
}
