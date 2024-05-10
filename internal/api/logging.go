package api

import (
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

func (s *Server) logging(statusCode *int, r *http.Request) {
	s.logger.WithFields(logrus.Fields{
		"request_time":   time.Now().Format("2006-01-02 15:04:05"),
		"request_method": r.Method,
		"request_url":    r.URL,
		"responce_code":  *statusCode,
	}).Debug("http request served")
}
