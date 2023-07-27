package loggers

import (
	"context"
	"notes-server/constants"
	"notes-server/utils"
	"time"

	colorable "github.com/mattn/go-colorable"
	logrus "github.com/sirupsen/logrus"
	"github.com/snowzach/rotatefilehook"
	"github.com/spf13/viper"
)

type Logger struct {
	logger *logrus.Logger
}

func NewLogger() *Logger {
	return &Logger{
		logger: Init(),
	}
}

// Init - Initialize the logger
func Init() *logrus.Logger {
	var logLevel logrus.Level
	switch viper.GetString("LOG_LEVEL") {
	case "INFO":
		logLevel = logrus.InfoLevel
	case "ERROR":
		logLevel = logrus.ErrorLevel
	case "DEBUG":
		logLevel = logrus.DebugLevel
	case "WARN":
		logLevel = logrus.WarnLevel
	default:
		logLevel = logrus.WarnLevel
	}
	logger := logrus.New()
	rotateFileHook, err := rotatefilehook.NewRotateFileHook(rotatefilehook.RotateFileConfig{
		Filename:   "loggers/logs/console.log",
		MaxSize:    50,
		MaxBackups: 3,
		MaxAge:     28,
		Level:      logLevel,
		Formatter: &logrus.JSONFormatter{
			TimestampFormat: time.RFC822,
		},
	})
	if err != nil {
		logger.Fatalf("Failed to initialize file rotate hook: %v", err)
	}
	logger.SetLevel(logLevel)
	logger.SetOutput(colorable.NewColorableStdout())
	logger.SetFormatter(&logrus.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: time.RFC822,
	})
	logger.AddHook(rotateFileHook)
	return logger
}

func (l *Logger) Info(ctx context.Context, args ...interface{}) {
	requestID := utils.GetRequestIDFromCtx(ctx)
	if requestID != "" {
		l.logger.WithField(constants.RequestIDKey, requestID).Info(args...)
	}
	l.logger.Info(args...)
}

func (l *Logger) Warn(ctx context.Context, args ...interface{}) {
	requestID := utils.GetRequestIDFromCtx(ctx)
	if requestID != "" {
		l.logger.WithField(constants.RequestIDKey, requestID).Warn(args...)
	}
	l.logger.Warn(args...)
}

func (l *Logger) Debug(ctx context.Context, args ...interface{}) {
	requestID := utils.GetRequestIDFromCtx(ctx)
	if requestID != "" {
		l.logger.WithField(constants.RequestIDKey, requestID).Debug(args...)
	}
	l.logger.Debug(args...)
}

func (l *Logger) Error(ctx context.Context, args ...interface{}) {
	requestID := utils.GetRequestIDFromCtx(ctx)
	if requestID != "" {
		l.logger.WithField(constants.RequestIDKey, requestID).Error(args...)
	}
	l.logger.Error(args...)
}

func (l *Logger) Fatal(ctx context.Context, args ...interface{}) {
	requestID := utils.GetRequestIDFromCtx(ctx)
	if requestID != "" {
		l.logger.WithField(constants.RequestIDKey, requestID).Fatal(args...)
	}
	l.logger.Fatal(args...)
}
