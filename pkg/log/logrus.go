package log

import (
	"bytes"
	"fmt"
	"sort"
	"sync"

	"github.com/sirupsen/logrus"
)

const (
	green = 32
)

func NewLogrus() *logrus.Logger {
	logger := logrus.New()
	logger.Formatter = &logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "15:04:05",
	}

	return logger
}

func NewPlan() *logrus.Logger {
	logger := logrus.New()
	logger.Formatter = &PlanFormatter{}

	return logger
}

// PlanFormatter formats logs into text
type PlanFormatter struct {
	// Whether the logger's out is to a terminal
	isTerminal bool

	sync.Once
}

func (f *PlanFormatter) init(entry *logrus.Entry) {
	if entry.Logger != nil {
		f.isTerminal = true // checkIfTerminal(entry.Logger.Out)
	}
}

// Format renders a single log entry
func (f *PlanFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	// prefixFieldClashes(entry.Data, f.FieldMap)

	keys := make([]string, 0, len(entry.Data))
	for k := range entry.Data {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	f.Do(func() { f.init(entry) })

	if f.isTerminal {
		f.printColored(b, entry, keys)
	} else {
		f.appendKeyValue(b, logrus.FieldKeyLevel, entry.Level.String())
		if entry.Message != "" {
			f.appendKeyValue(b, logrus.FieldKeyMsg, entry.Message)
		}
		for _, key := range keys {
			f.appendKeyValue(b, key, entry.Data[key])
		}
	}

	b.WriteByte('\n')
	return b.Bytes(), nil
}

func (f *PlanFormatter) printColored(b *bytes.Buffer, entry *logrus.Entry, keys []string) {
	levelText := "PLAN"

	fmt.Fprintf(b, "       \x1b[%dm%s\x1b[0m => %-44s ", green, levelText, entry.Message)
	for _, k := range keys {
		v := entry.Data[k]
		fmt.Fprintf(b, " \x1b[%dm%s\x1b[0m=", green, k)
		f.appendValue(b, v)
	}
}

func (f *PlanFormatter) needsQuoting(text string) bool {
	for _, ch := range text {
		if !((ch >= 'a' && ch <= 'z') ||
			(ch >= 'A' && ch <= 'Z') ||
			(ch >= '0' && ch <= '9') ||
			ch == '-' || ch == '.' || ch == '_' || ch == '/' || ch == '@' || ch == '^' || ch == '+') {
			return true
		}
	}
	return false
}

func (f *PlanFormatter) appendKeyValue(b *bytes.Buffer, key string, value interface{}) {
	if b.Len() > 0 {
		b.WriteByte(' ')
	}
	b.WriteString(key)
	b.WriteByte('=')
	f.appendValue(b, value)
}

func (f *PlanFormatter) appendValue(b *bytes.Buffer, value interface{}) {
	stringVal, ok := value.(string)
	if !ok {
		stringVal = fmt.Sprint(value)
	}

	if !f.needsQuoting(stringVal) {
		b.WriteString(stringVal)
	} else {
		b.WriteString(fmt.Sprintf("%q", stringVal))
	}
}
