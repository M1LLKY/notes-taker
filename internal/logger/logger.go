package logger

import (
	"context"
	"github.com/sirupsen/logrus"
	"notes-taker/internal/auth"
	"notes-taker/internal/requestmeta"
	"os"
	"time"
)

var logger *Logger

type Logger struct {
	Logrus *logrus.Logger
}

func Init(level string) {
	log := logrus.New()
	log.SetOutput(os.Stdout)
	log.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02T15:04:05.000-0700",
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "time",    // ключ для времени
			logrus.FieldKeyLevel: "level",   // ключ для уровня логирования
			logrus.FieldKeyMsg:   "message", // ключ для сообщения
		},
	})
	parsedLevel, err := logrus.ParseLevel(level)
	if err != nil {
		parsedLevel = logrus.InfoLevel
	}
	log.SetLevel(parsedLevel)
	logger = &Logger{
		Logrus: log,
	}
}

func Get() *Logger {
	return logger
}

// Info записывает информационное сообщение с полями и дополнительными данными из контекста.
func (l *Logger) Info(ctx context.Context, msg string, fields ...logrus.Fields) {
	f := fieldsOrEmpty(fields...)
	l.logWithUserData(ctx).WithFields(f).Info(msg)
}

// Error записывает сообщение об ошибке с полями и дополнительными данными из контекста.
func (l *Logger) Error(ctx context.Context, msg string, fields ...logrus.Fields) {
	f := fieldsOrEmpty(fields...)
	l.logWithUserData(ctx).WithFields(f).Error(msg)
}

func (l *Logger) logWithUserData(ctx context.Context) *logrus.Entry {
	fields := logrus.Fields{}

	if requestData, ok := requestmeta.FromContext(ctx); ok {
		fields["cutoff"] = time.Since(requestData.StartTime).Seconds()
		fields["handler-name"] = requestData.Method + " " + requestData.URL
	}

	if userData, ok := auth.FromContext(ctx); ok {
		if sub, err := userData.GetSubject(); err == nil {
			fields["sub"] = sub
		}
	}

	return l.Logrus.WithFields(fields)
}

func fieldsOrEmpty(fields ...logrus.Fields) logrus.Fields {
	if len(fields) > 0 {
		return fields[0]
	}
	return logrus.Fields{}
}
