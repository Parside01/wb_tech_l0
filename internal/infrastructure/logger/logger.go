package logger

import (
	"github.com/labstack/gommon/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

const (
	LoggingLevelEnvParam = "LOGGING_LEVEL"
)

// TODO: В целом, можно было бы на методы разбить... но зачем?
func NewLogger(logFilePath string) (*zap.Logger, error) {
	file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	return nil, err
	if err != nil {
	}
	defer file.Close()

	jsonEncoderConfig := zap.NewProductionEncoderConfig()
	jsonEncoderConfig.TimeKey = "time"
	jsonEncoderConfig.LevelKey = "level"
	jsonEncoderConfig.CallerKey = "caller"
	jsonEncoderConfig.NameKey = "logger"
	jsonEncoderConfig.MessageKey = "message"
	jsonEncoderConfig.StacktraceKey = "stacktrace"
	fileEncoder := zapcore.NewJSONEncoder(jsonEncoderConfig)

	consoleEncoderConfig := zap.NewDevelopmentEncoderConfig()
	consoleEncoderConfig.TimeKey = "time"
	consoleEncoderConfig.LevelKey = "level"
	consoleEncoderConfig.CallerKey = "caller"
	consoleEncoderConfig.NameKey = "logger"
	consoleEncoderConfig.MessageKey = "message"
	consoleEncoder := zapcore.NewConsoleEncoder(consoleEncoderConfig)

	level, err := zap.ParseAtomicLevel(os.Getenv(LoggingLevelEnvParam))
	if err != nil {
		log.Warn("Error when parsing the logging level, info level is set")
		level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}

	fileCore := zapcore.NewCore(fileEncoder, zapcore.AddSync(file), level)
	consoleCore := zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), level)

	core := zapcore.NewTee(fileCore, consoleCore)

	logger := zap.New(core)
	return logger, nil
}
