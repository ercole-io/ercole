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
	"io"
	"os"
	"path/filepath"
)

// Logger interface for a logger implementation
type Logger interface {
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	Panicf(format string, args ...interface{})

	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
	Panic(args ...interface{})

	setLevel(Level)
	setOutput(output io.Writer)
	setExitFunc(exitFunc func(int))
}

type LoggerOption func(Logger) error

func LogDirectory(logDirectory string) LoggerOption {
	return func(logger Logger) error {
		path := filepath.Join(logDirectory, "ercole-agent.log")
		logger.Infof("Logging on %s", path)

		f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}

		logger.setOutput(f)

		return nil
	}
}

func LogLevel(level Level) LoggerOption {
	return func(logger Logger) error {
		logger.setLevel(level)

		return nil
	}
}

func LogVerbosely(verbose bool) LoggerOption {
	if verbose {
		return LogLevel(DebugLevel)
	}

	return func(logger Logger) error { return nil }
}

func SetExitFunc(exitFunc func(int)) LoggerOption {
	return func(logger Logger) error {
		logger.setExitFunc(exitFunc)

		return nil
	}
}
