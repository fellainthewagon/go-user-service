package logging

import (
	"fmt"
	"io"
	"os"
	"path"
	"runtime"

	"github.com/sirupsen/logrus"
)

type writerHook struct {
	Writer    []io.Writer
	LogLevels []logrus.Level
}

func (hook *writerHook) Fire(e *logrus.Entry) error {
	line, err := e.String()
	if err != nil {
		return err
	}

	for _, w := range hook.Writer {
		w.Write([]byte(line))
	}

	return err
}

func (hook *writerHook) Levels() []logrus.Level {
	return hook.LogLevels
}

// ============================ //
var e *logrus.Entry

type Logger struct {
	*logrus.Entry
}

func GetLogger() Logger {
	return Logger{e}
}

func (l *Logger) GetLoggerWithField(k string, v interface{}) Logger {
	return Logger{l.WithField(k, v)}
}

// ============================ //

func init() {
	logger := logrus.New()

	logger.SetReportCaller(true)
	logger.Formatter = &logrus.TextFormatter{
		CallerPrettyfier: func(f *runtime.Frame) (function string, file string) {
			fileName := path.Base(f.File)

			return fmt.Sprintf("%s()", f.Function), fmt.Sprintf("%s:%d", fileName, f.Line)
		},
		DisableColors: true,
		FullTimestamp: true,
	}

	err := os.MkdirAll("logs", 0631)
	if err != nil {
		panic(err)
	}

	allLogsFile, err := os.OpenFile("logs/all.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0640)
	if err != nil {
		panic(err)
	}

	logger.SetOutput(io.Discard)

	logger.AddHook(&writerHook{
		Writer:    []io.Writer{allLogsFile, os.Stdout},
		LogLevels: logrus.AllLevels,
	})
	logger.SetLevel(logrus.TraceLevel)

	e = logrus.NewEntry(logger)
}
