package log

import (
	"io"
	"os"

	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Logger = lumberjack.Logger

var (
	logger  zerolog.Logger
	lwriter zerolog.LevelWriter
)

func init() {
	lwriter = zerolog.MultiLevelWriter(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: "2006/01/02 15:04:05"})
	logger = zerolog.New(lwriter).With().Timestamp().Logger().Level(zerolog.DebugLevel)
}

func Output(w *Logger, level int) *zerolog.Logger {
	// logger = zerolog.New(zerolog.MultiLevelWriter(_console, w)).With().Timestamp().Logger()
	lwriter = zerolog.MultiLevelWriter(lwriter, w)
	logger = logger.Output(lwriter).Level(zerolog.Level(level))
	return &logger
}

func Writer() io.Writer {
	return lwriter
}

func Log() *zerolog.Logger {
	return &logger
}
