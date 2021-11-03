package xlog

import (
	"runtime"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
)

var (

	// qualified package name, cached at first use
	thisPackage string

	// Positions in the call stack when tracing to report the calling method
	minimumCallerDepth int

	// Used for caller information initialisation
	callerInitOnce sync.Once
)

const (
	logrusPackage          = "github.com/sirupsen/logrus"
	maximumCallerDepth int = 25
	knownLogrusFrames  int = 7
)

func init() {
	// start at the bottom of the stack before the package-name cache is primed
	minimumCallerDepth = 1
}

// NewCallerHook creates a new CallerHook.
func NewCallerHook(levels []logrus.Level) *CallerHook {
	return &CallerHook{
		levels: levels,
	}
}

// CallerHook is the logrus.Hook that sets up the logrus.Entry.Caller field.
// Sets up the entry.Logger
type CallerHook struct {
	levels []logrus.Level
}

// Levels gets the levels that this hook is available for.
// Implements logrus.Hook.
func (c CallerHook) Levels() []logrus.Level {
	return c.levels
}

// Fire sets up the Caller field in the logrus.Entry.
// Implements logrus.Hook.
func (c CallerHook) Fire(entry *logrus.Entry) error {
	entry.Caller = getCallerX()
	return nil
}

// getCallerX retrieves the name of the first non-logrus calling function
func getCallerX() *runtime.Frame {
	// cache this package's fully-qualified name
	callerInitOnce.Do(func() {
		pcs := make([]uintptr, maximumCallerDepth)
		_ = runtime.Callers(0, pcs)

		// dynamic get the package name and the minimum caller depth
		for i := 0; i < maximumCallerDepth; i++ {
			funcName := runtime.FuncForPC(pcs[i]).Name()
			if strings.Contains(funcName, "getCallerX") {
				thisPackage = getPackageName(funcName)
				break
			}
		}

		minimumCallerDepth = knownLogrusFrames
	})

	// Restrict the lookback frames to avoid runaway lookups
	pcs := make([]uintptr, maximumCallerDepth)
	depth := runtime.Callers(minimumCallerDepth, pcs)
	frames := runtime.CallersFrames(pcs[:depth])

	for f, again := frames.Next(); again; f, again = frames.Next() {
		pkg := getPackageName(f.Function)

		// If the caller isn't part of this package, we're done
		if pkg != thisPackage && pkg != logrusPackage {
			return &f //nolint:scopelint
		}
	}

	// if we got here, we failed to find the caller's context
	return nil
}

// getPackageName reduces a fully qualified function name to the package name
// There really ought to be a better way...
func getPackageName(f string) string {
	for {
		lastPeriod := strings.LastIndex(f, ".")
		lastSlash := strings.LastIndex(f, "/")
		if lastPeriod > lastSlash {
			f = f[:lastPeriod]
		} else {
			break
		}
	}

	return f
}
