// Copyright (c) 2022 Sorint.lab S.p.A.
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

package utils

import (
	"io"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/sirupsen/logrus"

	"github.com/ercole-io/ercole/v2/logger"
)

// CustomLoggingHandler return a loggingHandler using ercole formatting
func CustomLoggingHandler(router http.Handler, log logger.Logger) http.Handler {
	return handlers.CustomLoggingHandler(os.Stdout, router, createLogFormatter(log))
}

func createLogFormatter(log logger.Logger) func(writer io.Writer, params handlers.LogFormatterParams) {
	return func(_ io.Writer, params handlers.LogFormatterParams) {
		req := params.Request
		fields := logrus.Fields{
			"endpoint":     req.Method + " " + req.URL.String(),
			"userAgent":    req.UserAgent(),
			"accept":       req.Header.Get("Accept"),
			"serverSocket": req.Host,
			"clientSocket": req.RemoteAddr,
			"timeStamp":    params.TimeStamp,
			"statusCode":   params.StatusCode,
			"size":         params.Size,
		}

		switch l := log.(type) {
		case *logger.LogrusLogger:
			l.WithFields(fields).Info("HTTP Request")
		default:
			log.Info("HTTP Request %v", fields)
		}
	}
}
