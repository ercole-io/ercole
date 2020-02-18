package utils

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
)

// NewLogger return a logrus.Logger initialized with ercole log standard
func NewLogger(componentName string) *logrus.Logger {
	logger := logrus.New()
	logger.SetFormatter(&ercoleFormatter{ComponentName: componentName[0:4]})
	logger.SetReportCaller(true)

	return logger
}

// ercoleFormatter custom formatter for ercole that formats logs into text
type ercoleFormatter struct {
	ComponentName string
}

// Format renders a single log entry
func (f *ercoleFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	levelColor := getColorByLevel(entry)
	levelText := strings.ToUpper(entry.Level.String())[0:4]
	caller := getCaller(entry)
	message := strings.TrimSuffix(entry.Message, "\n")

	var msg bytes.Buffer
	msg.WriteString(
		fmt.Sprintf("\x1b[%dm[%s][%s][%s]\x1b[0m[%s] %-44s",
			levelColor,
			entry.Time.Format("06-01-02 15:04"),
			f.ComponentName,
			levelText,
			caller,
			message))

	return append(msg.Bytes(), '\n'), nil
}

func getColorByLevel(entry *logrus.Entry) int {
	const gray = 37
	const yellow = 33
	const red = 31
	const blue = 36

	switch entry.Level {
	case logrus.DebugLevel, logrus.TraceLevel:
		return gray
	case logrus.WarnLevel:
		return yellow
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		return red
	default:
		return blue
	}
}

func getCaller(entry *logrus.Entry) string {
	if !entry.HasCaller() {
		return ""
	}

	defer func() string {
		recover()
		return ""
	}()

	caller := strings.SplitAfter(entry.Caller.File, "ercole-services/")[1]

	return fmt.Sprintf("%s:%d", caller, entry.Caller.Line)
}
