// Copyright (c) 2020 Sorint.lab S.p.A.
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package logger

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"runtime"
	"sort"
	"strings"

	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
)

// LogrusLogger struct to compose logger with logrus that satisfy Logger interface
type LogrusLogger struct {
	*logrus.Logger
}

// SetLevel to inner field log
func (l *LogrusLogger) setLevel(level Level) {
	l.Logger.Level = logrus.Level(level)
}

func (l *LogrusLogger) setOutput(output io.Writer) {
	l.Logger.SetOutput(output)
}

func (l *LogrusLogger) setExitFunc(exitFunc func(int)) {
	l.Logger.ExitFunc = exitFunc
}

// NewLogger return a LogrusLogger initialized with ercole log standard
// If it fail, it will exit
func NewLogger(componentName string, options ...LoggerOption) Logger {
	var newLogger LogrusLogger
	newLogger.Logger = logrus.New()

	if len(componentName) > 4 {
		componentName = componentName[0:4]
	}

	newLogger.SetFormatter(&ercoleFormatter{
		ComponentName: componentName,
		isColored:     runtime.GOOS != "windows",
	})
	newLogger.SetReportCaller(true)
	newLogger.SetOutput(os.Stdout)

	for _, option := range options {
		err := option(&newLogger)
		if err != nil {
			fmt.Printf("Can't initialize %s logger: %s", componentName, err)
			os.Exit(1)
		}
	}

	return &newLogger
}

// ercoleFormatter custom formatter for ercole that formats logs into text
type ercoleFormatter struct {
	ComponentName string
	isColored     bool
}

// Format renders a single log entry
func (f *ercoleFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var scolored func(format string, a ...interface{}) string

	if f.isColored {
		colorSprintFunc := color.New(getColorByLevel(Level(entry.Level))).SprintFunc()
		scolored = func(format string, a ...interface{}) string {
			a = append([]interface{}{format}, a...)
			return colorSprintFunc(a...)
		}
	} else {
		scolored = fmt.Sprintf
	}

	levelText := strings.ToUpper(entry.Level.String())[0:4]
	caller := getCaller(entry)
	message := strings.TrimSuffix(entry.Message, "\n")

	var logBuffer bytes.Buffer

	logBuffer.WriteString(
		scolored(
			fmt.Sprintf("[%s][%s][%s]",
				entry.Time.Format("06-01-02 15:04:05"),
				f.ComponentName,
				levelText)))

	logBuffer.WriteString(fmt.Sprintf("[%s] %-50s", caller, message))

	for _, k := range getKeysInOrder(entry.Data) {
		logBuffer.WriteString(fmt.Sprintf("%s=%v ", scolored(k), entry.Data[k]))
	}

	return append(logBuffer.Bytes(), '\n'), nil
}

func getCaller(entry *logrus.Entry) string {
	if !entry.HasCaller() {
		return ""
	}

	caller := entry.Caller.File

	removeFrom := "ercole/"
	if strings.Contains(caller, removeFrom) {
		caller = "./" + caller[strings.Index(caller, removeFrom)+len(removeFrom):]
	}

	return fmt.Sprintf("%s:%d", caller, entry.Caller.Line)
}

func getKeysInOrder(entryData logrus.Fields) []string {
	manuallyOrderedKeys := []string{"endpoint", "statusCode"}

	for i := 0; i < len(manuallyOrderedKeys); i++ {
		if _, ok := entryData[manuallyOrderedKeys[i]]; !ok {
			manuallyOrderedKeys = remove(manuallyOrderedKeys, i)
			i--
		}
	}

	var entryDataKeys []string

	for k := range entryData {
		if !contains(manuallyOrderedKeys, k) {
			entryDataKeys = append(entryDataKeys, k)
		}
	}

	sort.Strings(entryDataKeys)

	return append(manuallyOrderedKeys, entryDataKeys...)
}

// contains return true if a contains x, otherwise false.
func contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}

	return false
}

// remove return slice without element at position i, mantaining order
func remove(slice []string, i int) []string {
	return append(slice[:i], slice[i+1:]...)
}
