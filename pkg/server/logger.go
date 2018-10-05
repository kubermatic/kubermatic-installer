package server

import (
	"encoding/json"
	"time"

	"github.com/sirupsen/logrus"
)

type logItem []byte

// newLogger returns a logrus instance prepared to write
// its entries HTML-formatted to the response writer.
func newLogger(upstream logrus.FieldLogger, buf chan logItem) *logrus.Logger {
	logger := logrus.New()

	// write messages to the buffered channel by default
	logger.SetOutput(&writer{buf})
	logger.Formatter = &formatter{}

	// ... and also copy them to the CLI logger
	logger.AddHook(&splitter{upstream})

	return logger
}

// splitter is a logrus.Hook that for every event
// send a copy to another logrus Logger
type splitter struct {
	upstream logrus.FieldLogger
}

func (s *splitter) Levels() []logrus.Level {
	return logrus.AllLevels
}

// see https://github.com/sirupsen/logrus/pull/783 and weep
func (s *splitter) Fire(e *logrus.Entry) error {
	switch e.Level {
	case logrus.PanicLevel:
		s.upstream.Panic(e.Message)
	case logrus.FatalLevel:
		s.upstream.Fatal(e.Message)
	case logrus.ErrorLevel:
		s.upstream.Error(e.Message)
	case logrus.WarnLevel:
		s.upstream.Warn(e.Message)
	case logrus.InfoLevel:
		s.upstream.Info(e.Message)
	case logrus.DebugLevel:
		s.upstream.Debug(e.Message)
	}
	return nil
}

// writer is a wrapper around the HTTP response writer
// that flushes after each Write operation.
type writer struct {
	buf chan logItem
}

func (w *writer) Write(p []byte) (n int, err error) {
	w.buf <- p
	return len(p), nil
}

// formatter is an HTML-enabled formatter for logrus entries.
type formatter struct{}

type logMessage struct {
	Type    string       `json:"type"`
	Date    time.Time    `json:"date"`
	Level   logrus.Level `json:"level"`
	Message string       `json:"message"`
}

func (f *formatter) Format(entry *logrus.Entry) ([]byte, error) {
	msg := logMessage{
		Type:    "log",
		Date:    entry.Time,
		Level:   entry.Level,
		Message: entry.Message,
	}

	encoded, _ := json.Marshal(msg)

	return encoded, nil
}
