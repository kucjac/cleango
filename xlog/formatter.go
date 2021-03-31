package xlog

import (
	"bytes"
	"fmt"
	"sort"
	"strings"

	"github.com/sirupsen/logrus"
)

// TextFormatter is struct implementing Format interface
// this is useful for fomating logs indifferent environments.
// This formater will format logs for terminal and testing.
type TextFormatter struct {
	// Force disabling colors. For a TTY colors are enabled by default.
	UseColors bool
	// Color scheme to use.
	scheme *compiledColorScheme
}

// NewTextFormatter will create new formtter for logrus
func NewTextFormatter(colors bool) *TextFormatter {
	f := &TextFormatter{}
	f.scheme = noColorsColorScheme
	if colors {
		f.scheme = defaultCompiledColorScheme
	}
	return f
}

// Format to make interface `logrus.Formatter` happy. It will take entry
// and covert it into byte stream.
func (f *TextFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	x, errSet := ExtractHTTPField(entry.Data)
	switch errSet {
	case nil:
		entry.Message = fmt.Sprintf(
			"[%s] %d | %8v | %15s",
			x.RequestMethod,
			x.Status,
			x.Latency,
			x.RequestURL,
		)
		delete(entry.Data, HTTPRequestKey)
	case NoHTTPOpt:
	default:
		return nil, errSet
	}
	var b *bytes.Buffer
	var keys = make([]string, 0, len(entry.Data))
	for k := range entry.Data {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}
	f.printColored(b, entry, keys)
	b.WriteByte('\n')
	return b.Bytes(), nil
}

func (f *TextFormatter) printColored(b *bytes.Buffer, entry *logrus.Entry, keys []string) {
	var levelColor func(string) string
	var levelText string
	switch entry.Level {
	case logrus.InfoLevel:
		levelColor = f.scheme.InfoLevelColor
	case logrus.WarnLevel:
		levelColor = f.scheme.WarnLevelColor
	case logrus.ErrorLevel:
		levelColor = f.scheme.ErrorLevelColor
	case logrus.FatalLevel:
		levelColor = f.scheme.FatalLevelColor
	case logrus.PanicLevel:
		levelColor = f.scheme.PanicLevelColor
	default:
		levelColor = f.scheme.DebugLevelColor
	}
	if entry.Level != logrus.WarnLevel {
		levelText = entry.Level.String()
	} else {
		levelText = "warn"
	}
	levelText = strings.ToUpper(levelText)
	level := levelColor(fmt.Sprintf("%-5s", levelText))
	message := entry.Message
	messageFormat := "%s"
	fmt.Fprintf(b, "%s "+messageFormat, level, message)
	for _, k := range keys {
		v := entry.Data[k]
		data := fmt.Sprintf("%+v", v)
		fmt.Fprintf(b, " %s=%q", fmt.Sprintf("\"%s\"", levelColor(k)), data)
	}
}
