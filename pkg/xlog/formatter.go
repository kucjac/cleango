package xlog

import (
	"bytes"
	"fmt"
	"io"
	"runtime"
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/kucjac/cleango/cgerrors"
)

// TextFormatter is struct implementing Format interface
// this is useful for formatting logs in different environments.
// This formatter will format logs for terminal and testing.
type TextFormatter struct {
	// TimestampFormat sets the format used for marshaling timestamps.
	// The format to use is the same than for time.Format or time.Parse from the standard
	// library.
	// The standard Library already provides a set of predefined format.
	TimestampFormat string

	// Force disabling colors. For a TTY colors are enabled by default.
	UseColors bool

	// DisableTimestamp allows disabling automatic timestamps in output
	DisableTimestamp bool

	// DataKey allows users to put all the log entry parameters into a nested dictionary at a given key.
	DataKey string

	// CallerPrettyfier can be set by the user to modify the content
	// of the function and file keys in the json data when ReportCaller is
	// activated. If any of the returned value is the empty string the
	// corresponding key will be removed from fields.
	CallerPrettyfier func(*runtime.Frame) (function string, file string)

	// Color scheme to use.
	scheme *compiledColorScheme
}

// NewTextFormatter creates new logrus based text formatter.
func NewTextFormatter(colors bool) *TextFormatter {
	f := &TextFormatter{
		scheme:          noColorsColorScheme,
		TimestampFormat: time.RFC3339,
	}
	if colors {
		f.scheme = defaultCompiledColorScheme
	}
	return f
}

// Format implements `logrus.Formatter`.
func (f *TextFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	x, err := ExtractHTTPField(entry.Data)
	if err != nil && !cgerrors.Is(err, NoHTTPOpt) {
		return nil, err
	}
	if err == nil {
		entry.Message = fmt.Sprintf("[%s] %d | %8v | %15s", x.RequestMethod, x.Status, x.Latency, x.RequestURL)
		delete(entry.Data, HTTPRequestKey)
	}

	data := make(logrus.Fields, len(entry.Data)+4)
	for k, v := range entry.Data {
		switch v := v.(type) {
		case error:
			data[k] = v.Error()
		default:
			data[k] = v
		}
	}
	if f.DataKey != "" {
		newData := make(logrus.Fields, 4)
		newData[f.DataKey] = data
		data = newData
	}

	tsFormat := f.TimestampFormat
	if tsFormat == "" {
		tsFormat = time.RFC3339
	}

	if !f.DisableTimestamp {
		data[logrus.FieldKeyTime] = entry.Time.Format(tsFormat)
	}

	if entry.HasCaller() {
		funcVal := entry.Caller.Function
		fileVal := fmt.Sprintf("%s:%d", entry.Caller.File, entry.Caller.Line)
		if f.CallerPrettyfier != nil {
			funcVal, fileVal = f.CallerPrettyfier(entry.Caller)
		}
		if funcVal != "" {
			data["func"] = funcVal
		}
		if fileVal != "" {
			data["file"] = fileVal
		}
	}

	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}
	f.printColored(b, entry, data)
	b.WriteByte('\n')
	return b.Bytes(), nil
}

func (f *TextFormatter) printColored(b io.Writer, entry *logrus.Entry, data logrus.Fields) {
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
	for k, v := range data {
		data := fmt.Sprintf("%+v", v)
		fmt.Fprintf(b, " %s=%q", fmt.Sprintf("\"%s\"", levelColor(k)), data)
	}
}
