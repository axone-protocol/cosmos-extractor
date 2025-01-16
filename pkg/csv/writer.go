package csv

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/teambenny/goetl"
	"github.com/teambenny/goetl/etldata"
	"github.com/teambenny/goetl/etlutil"
	"github.com/teambenny/goetl/processors"
)

type Option func(*writer) error

func WithWriterHeader() Option {
	return func(w *writer) error {
		w.delegated.Parameters.WriteHeader = true

		return nil
	}
}

func WithFile(file string) Option {
	return WithFileAndFlag(file, os.O_RDWR|os.O_CREATE)
}

func WithFileAndFlag(file string, flag int) Option {
	return func(w *writer) error {
		file, err := os.OpenFile(file, flag, os.ModePerm)
		if err != nil {
			return err
		}
		bw := bufio.NewWriter(file)

		delegated := processors.NewCSVWriter(bw)

		if w.delegated != nil {
			copyCSVWriterParameters(w.delegated.Parameters, &delegated.Parameters)
		} else {
			delegated.Parameters.WriteHeader = false
		}

		w.finalizer = func() error {
			if err := bw.Flush(); err != nil {
				return err
			}
			return file.Close()
		}
		w.delegated = delegated
		w.processorName = fmt.Sprintf("CSVWriter<%s>", file.Name())

		return nil
	}
}

func WithWriter(w io.Writer) Option {
	return func(cw *writer) error {
		delegated := processors.NewCSVWriter(w)

		if cw.delegated != nil {
			copyCSVWriterParameters(cw.delegated.Parameters, &delegated.Parameters)
		} else {
			delegated.Parameters.WriteHeader = false
		}

		if closer, ok := w.(flusher); ok {
			cw.finalizer = closer.Flush
		} else {
			cw.finalizer = func() error { return nil }
		}

		cw.delegated = delegated

		if stringer, ok := w.(fmt.Stringer); ok {
			cw.processorName = fmt.Sprintf("CSVWriter<%s>", stringer.String())
		} else {
			cw.processorName = "CSVWriter<custom>"
		}

		return nil
	}
}

type flusher interface {
	Flush() error
}

type writer struct {
	finalizer     func() error
	delegated     *processors.CSVWriter
	processorName string
}

func NewCSVWriter(options ...Option) (goetl.Processor, error) {
	delegated := defaultCSVWriter()
	writer := &writer{
		finalizer:     func() error { return nil },
		delegated:     delegated,
		processorName: "CSVWriter<null>",
	}

	for _, option := range options {
		err := option(writer)
		if err != nil {
			return nil, err
		}
	}

	return writer, nil
}

func (w *writer) ProcessData(d etldata.Payload, outputChan chan etldata.Payload, killChan chan error) {
	w.delegated.ProcessData(d, outputChan, killChan)
}

func (w *writer) Finish(outputChan chan etldata.Payload, killChan chan error) {
	w.delegated.Finish(outputChan, killChan)
	killChan <- w.finalizer()
}

func (w *writer) String() string {
	return w.processorName
}

func defaultCSVWriter() *processors.CSVWriter {
	csvW := processors.NewCSVWriter(io.Discard)
	csvW.Parameters.WriteHeader = false

	return csvW
}

func copyCSVWriterParameters(from etlutil.CSVParameters, to *etlutil.CSVParameters) {
	to.WriteHeader = from.WriteHeader
	to.HeaderWritten = from.HeaderWritten
	to.SendUpstream = from.SendUpstream
}
