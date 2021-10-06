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
