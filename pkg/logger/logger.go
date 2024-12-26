package logger

import (
	"fmt"
	"io"
	"strings"

	"github.com/sourcegraph/conc/iter"
	"github.com/teambenny/goetl/logger"

	"cosmossdk.io/log"
)

type etlLogger struct {
	logger log.Logger
}

func NewETLLogger(l log.Logger) logger.ETLNotifier {
	return &etlLogger{
		logger: l,
	}
}

func InstallETLLogger(l log.Logger) {
	logger.Notifier = NewETLLogger(l)
	logger.SetOutput(io.Discard)
}

// ETLNotifier is an interface for receiving log events from goetl.
func (l *etlLogger) ETLNotify(lvl int, _ []byte, v ...interface{}) {
	s := func() string {
		return strings.Join(iter.Map(v, func(it *interface{}) string {
			return fmt.Sprint(*it)
		}), " ")
	}

	switch lvl {
	case logger.LevelDebug:
		l.logger.Debug(s())
	case logger.LevelInfo:
		l.logger.Info(s())
	case logger.LevelError:
		l.logger.Error(s())
	case logger.LevelStatus:
		l.logger.Info(s())
	case logger.LevelSilent:
		// shh
	}
}
