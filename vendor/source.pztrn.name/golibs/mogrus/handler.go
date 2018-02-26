package mogrus

import (
	// stdlib
	"io"
	"strings"
	"sync"

	// github
	"github.com/sirupsen/logrus"
)

type LoggerHandler struct {
	// Logrus instances
	instances      map[string]*logrus.Logger
	instancesMutex sync.Mutex
}

// Adds output for logger handler.
// This actually creates new Logrus's Logger instance and configure
// it to write to given writer.
// To configure debug level you should pass it's name as debugLvl.
// Valid values: "", "debug" (the default), "info", "warn", "error"
func (lh *LoggerHandler) CreateOutput(name string, writer io.Writer, colors bool, debugLvl string) {
	// Formatter.
	logrus_formatter := new(logrus.TextFormatter)
	logrus_formatter.DisableTimestamp = false
	logrus_formatter.FullTimestamp = true
	logrus_formatter.QuoteEmptyFields = true
	logrus_formatter.TimestampFormat = "2006-01-02 15:04:05.000000000"

	if colors {
		logrus_formatter.DisableColors = false
		logrus_formatter.ForceColors = true
	} else {
		logrus_formatter.DisableColors = true
		logrus_formatter.ForceColors = false
	}

	logrus_instance := logrus.New()
	logrus_instance.Out = writer
	// Defaulting to debug.
	logrus_instance.Level = logrus.DebugLevel
	if debugLvl == "info" {
		logrus_instance.Level = logrus.InfoLevel
	} else if debugLvl == "warn" {
		logrus_instance.Level = logrus.WarnLevel
	} else if debugLvl == "error" {
		logrus_instance.Level = logrus.ErrorLevel
	}
	logrus_instance.Formatter = logrus_formatter

	lh.instancesMutex.Lock()
	lh.instances[name] = logrus_instance

	for _, logger := range lh.instances {
		logger.Debugln("Added new logger output:", name)
	}
	lh.instancesMutex.Unlock()
}

// Formats string by replacing "{{ var }}"'s with data from passed map.
func (lh *LoggerHandler) FormatString(data string, replacers map[string]string) string {
	for k, v := range replacers {
		data = strings.Replace(data, "{{ "+k+" }}", v, -1)
	}

	return data
}

// Initializes logger handler.
// It will only initializes LoggerHandler structure, see CreateOutput()
// for configuring output for this logger handler.
func (lh *LoggerHandler) Initialize() {
	lh.instances = make(map[string]*logrus.Logger)
}

// Removes previously created output.
// If output isn't found - doing nothing.
func (lh *LoggerHandler) RemoveOutput(output_name string) {
	lh.instancesMutex.Lock()
	_, found := lh.instances[output_name]
	if found {
		delete(lh.instances, output_name)
	}
	lh.instancesMutex.Unlock()
}

// Proxy for Logrus's Logger.Debug() function.
func (lh *LoggerHandler) Debug(args ...interface{}) {
	lh.instancesMutex.Lock()
	for _, logger := range lh.instances {
		logger.Debug(args...)
	}
	lh.instancesMutex.Unlock()
}

// Proxy for Logrus's Logger.Debugf() function.
func (lh *LoggerHandler) Debugf(format string, args ...interface{}) {
	lh.instancesMutex.Lock()
	for _, logger := range lh.instances {
		logger.Debugf(format, args...)
	}
	lh.instancesMutex.Unlock()
}

// Proxy for Logrus's Logger.Debugln() function.
func (lh *LoggerHandler) Debugln(args ...interface{}) {
	lh.instancesMutex.Lock()
	for _, logger := range lh.instances {
		logger.Debugln(args...)
	}
	lh.instancesMutex.Unlock()
}

// Proxy for Logrus's Logger.Error() function.
func (lh *LoggerHandler) Error(args ...interface{}) {
	lh.instancesMutex.Lock()
	for _, logger := range lh.instances {
		logger.Error(args...)
	}
	lh.instancesMutex.Unlock()
}

// Proxy for Logrus's Logger.Errorf() function.
func (lh *LoggerHandler) Errorf(format string, args ...interface{}) {
	lh.instancesMutex.Lock()
	for _, logger := range lh.instances {
		logger.Errorf(format, args...)
	}
	lh.instancesMutex.Unlock()
}

// Proxy for Logrus's Logger.Errorln() function.
func (lh *LoggerHandler) Errorln(args ...interface{}) {
	lh.instancesMutex.Lock()
	for _, logger := range lh.instances {
		logger.Errorln(args...)
	}
	lh.instancesMutex.Unlock()
}

// Proxy for Logrus's Logger.Fatal() function.
func (lh *LoggerHandler) Fatal(args ...interface{}) {
	lh.instancesMutex.Lock()
	for _, logger := range lh.instances {
		logger.Fatal(args...)
	}
	lh.instancesMutex.Unlock()
}

// Proxy for Logrus's Logger.Fatalf() function.
func (lh *LoggerHandler) Fatalf(format string, args ...interface{}) {
	lh.instancesMutex.Lock()
	for _, logger := range lh.instances {
		logger.Fatalf(format, args...)
	}
	lh.instancesMutex.Unlock()
}

// Proxy for Logrus's Logger.Fatalln() function.
func (lh *LoggerHandler) Fatalln(args ...interface{}) {
	lh.instancesMutex.Lock()
	for _, logger := range lh.instances {
		logger.Fatalln(args...)
	}
	lh.instancesMutex.Unlock()
}

// Proxy for Logrus's Logger.Info() function.
func (lh *LoggerHandler) Info(args ...interface{}) {
	lh.instancesMutex.Lock()
	for _, logger := range lh.instances {
		logger.Info(args...)
	}
	lh.instancesMutex.Unlock()
}

// Proxy for Logrus's Logger.Infof() function.
func (lh *LoggerHandler) Infof(format string, args ...interface{}) {
	lh.instancesMutex.Lock()
	for _, logger := range lh.instances {
		logger.Infof(format, args...)
	}
	lh.instancesMutex.Unlock()
}

// Proxy for Logrus's Logger.Infoln() function.
func (lh *LoggerHandler) Infoln(args ...interface{}) {
	lh.instancesMutex.Lock()
	for _, logger := range lh.instances {
		logger.Infoln(args...)
	}
	lh.instancesMutex.Unlock()
}

// Proxy for Logrus's Logger.Print() function.
func (lh *LoggerHandler) Print(args ...interface{}) {
	lh.instancesMutex.Lock()
	for _, logger := range lh.instances {
		logger.Print(args...)
	}
	lh.instancesMutex.Unlock()
}

// Proxy for Logrus's Logger.Printf() function.
func (lh *LoggerHandler) Printf(format string, args ...interface{}) {
	lh.instancesMutex.Lock()
	for _, logger := range lh.instances {
		logger.Printf(format, args...)
	}
	lh.instancesMutex.Unlock()
}

// Proxy for Logrus's Logger.Println() function.
func (lh *LoggerHandler) Println(args ...interface{}) {
	lh.instancesMutex.Lock()
	for _, logger := range lh.instances {
		logger.Println(args...)
	}
	lh.instancesMutex.Unlock()
}

// Proxy for Logrus's Logger.Warn() function.
func (lh *LoggerHandler) Warn(args ...interface{}) {
	lh.instancesMutex.Lock()
	for _, logger := range lh.instances {
		logger.Warn(args...)
	}
	lh.instancesMutex.Unlock()
}

// Proxy for Logrus's Logger.Warnf() function.
func (lh *LoggerHandler) Warnf(format string, args ...interface{}) {
	lh.instancesMutex.Lock()
	for _, logger := range lh.instances {
		logger.Warnf(format, args...)
	}
	lh.instancesMutex.Unlock()
}

// Proxy for Logrus's Logger.Warnln() function.
func (lh *LoggerHandler) Warnln(args ...interface{}) {
	lh.instancesMutex.Lock()
	for _, logger := range lh.instances {
		logger.Warnln(args...)
	}
	lh.instancesMutex.Unlock()
}
