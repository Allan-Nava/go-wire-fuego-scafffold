package app_log

import (
	"fmt"
	"os"
	"strings"

	"github.com/Paxx-RnD/go-helper/helpers/string_helper"
	"go.uber.org/zap"
	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
)

func NewLogger(logLevel string) *zap.SugaredLogger {
	level := getLevel(logLevel)
	cfg := zap.Config{
		Encoding:         "json",
		Level:            zap.NewAtomicLevelAt(level),
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:   "message",
			LevelKey:     "level",
			EncodeLevel:  zapcore.CapitalColorLevelEncoder,
			TimeKey:      "time",
			EncodeTime:   zapcore.ISO8601TimeEncoder,
			CallerKey:    "caller",
			EncodeCaller: zapcore.ShortCallerEncoder,
		},
	}
	logger, _ := cfg.Build()

	logger = logger.WithOptions(
		zap.WrapCore(
			func(zapcore.Core) zapcore.Core {
				return zapcore.NewCore(zapcore.NewConsoleEncoder(cfg.EncoderConfig), zapcore.AddSync(os.Stderr), level)
			}))

	defer func(logger *zap.Logger) {
		_ = logger.Sync()
	}(logger)
	return logger.Sugar()
}

func getLevel(level string) zapcore.Level {
	lower := strings.ToLower(level)
	switch lower {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	case "fatal":
		return zapcore.FatalLevel
	case "panic":
		return zapcore.PanicLevel
	default:
		return zapcore.InfoLevel
	}
}

func GetThreadLogger(threadID string, jobID int, title string, service string, level string) *zap.SugaredLogger {
	logLevel := getLevel(level)
	cfg := zap.Config{
		Encoding:         "json",
		Level:            zap.NewAtomicLevelAt(logLevel),
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:   "message",
			LevelKey:     "level",
			EncodeLevel:  zapcore.CapitalColorLevelEncoder,
			TimeKey:      "time",
			EncodeTime:   zapcore.ISO8601TimeEncoder,
			CallerKey:    "caller",
			EncodeCaller: nil,
		},
	}
	logger, _ := cfg.Build()

	const maxLength = 20
	title = string_helper.Truncate(title, maxLength)

	enc := &prependEncoder{
		Encoder: zapcore.NewConsoleEncoder(cfg.EncoderConfig),
		pool:    buffer.NewPool(),
		text:    fmt.Sprintf("%s|%d|%s -> [%s] ", threadID[:4], jobID, title, service),
	}

	logger = logger.WithOptions(
		zap.WrapCore(
			func(zapcore.Core) zapcore.Core {
				return zapcore.NewCore(enc, zapcore.AddSync(os.Stderr), logLevel)
			}))

	defer func(logger *zap.Logger) {
		_ = logger.Sync()
	}(logger)
	return logger.Sugar()
}

type prependEncoder struct {
	zapcore.Encoder
	pool buffer.Pool
	text string
}

func (e *prependEncoder) EncodeEntry(entry zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
	buf := e.pool.Get()

	entry.Message = e.text + entry.Message
	output, err := e.Encoder.EncodeEntry(entry, fields)
	if err != nil {
		return nil, err
	}

	_, err = buf.Write(output.Bytes())
	if err != nil {
		return nil, err
	}
	return buf, nil
}
