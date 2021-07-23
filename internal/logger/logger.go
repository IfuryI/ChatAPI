package logger

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Logger структуры логгера
type Logger struct {
	*logrus.Logger
}

// NewAccessLogger инизиализация структуры логгера
func NewAccessLogger() *Logger {
	lg := logrus.New()
	return &Logger{lg}
}

// StartReq начало запроса
func (l *Logger) StartReq(r http.Request, rid string) {
	l.WithFields(logrus.Fields{
		"id":        rid,
		"usr_addr":  r.RemoteAddr,
		"req_URI":   r.RequestURI,
		"method":    r.Method,
		"usr_agent": r.UserAgent(),
	}).Info("request started")
}

// EndReq конец запроса
func (l *Logger) EndReq(r http.Request, start time.Time, rid string) {
	l.WithFields(logrus.Fields{
		"id":        rid,
		"usr_addr":  r.RemoteAddr,
		"req_URI":   r.RequestURI,
		"method":    r.Method,
		"usr_agent": r.UserAgent(),
		"μs":        time.Since(start).Microseconds(),
	}).Info("request ended")
}

// HTTPInfo статус запроса
func (l *Logger) HTTPInfo(ctx *gin.Context, msg string, status int) {
	l.WithFields(logrus.Fields{
		"status": status,
	}).Info(msg)
}

// LogWarning предупреждение
func (l *Logger) LogWarning(ctx *gin.Context, pkg string, funcName string, msg string) {
	l.WithFields(logrus.Fields{
		"package":  pkg,
		"function": funcName,
	}).Warn(msg)
}

// LogError информация об ошибке
func (l *Logger) LogError(ctx *gin.Context, pkg string, funcName string, err error) {
	l.WithFields(logrus.Fields{
		"package":  pkg,
		"function": funcName,
	}).Error(err)
}
