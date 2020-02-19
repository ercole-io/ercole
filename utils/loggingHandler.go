package utils

import (
	"io"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// CustomLoggingHandler return a loggingHandler using ercole formatting
func CustomLoggingHandler(router *mux.Router, log *logrus.Logger) http.Handler {
	return handlers.CustomLoggingHandler(os.Stdout, router, createLogFormatter(log))
}

func createLogFormatter(log *logrus.Logger) func(writer io.Writer, params handlers.LogFormatterParams) {
	return func(_ io.Writer, params handlers.LogFormatterParams) {
		req := params.Request
		log.
			WithFields(logrus.Fields{
				"endpoint":     req.Method + " " + req.URL.String(),
				"userAgent":    req.UserAgent(),
				"accept":       req.Header.Get("Accept"),
				"serverSocket": req.Host,
				"clientSocket": req.RemoteAddr,
				"timeStamp":    params.TimeStamp,
				"statusCode":   params.StatusCode,
				"size":         params.Size,
			}).
			Info("HTTP Request")
	}
}
