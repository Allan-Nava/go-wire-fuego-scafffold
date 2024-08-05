package app_log


import (
	"context"
	"go.uber.org/zap"
	"log/slog"
)

type LoggerAdapter struct {
	logger *zap.SugaredLogger
}

func NewLoggerAdapter(logger *zap.SugaredLogger) *LoggerAdapter {
	return &LoggerAdapter{logger: logger}
}

func (a *LoggerAdapter) Enabled(_ context.Context, level slog.Level) bool {
	switch level {
	case slog.LevelDebug:
		return a.logger.Desugar().Core().Enabled(zap.DebugLevel)
	case slog.LevelInfo:
		return a.logger.Desugar().Core().Enabled(zap.InfoLevel)
	case slog.LevelWarn:
		return a.logger.Desugar().Core().Enabled(zap.WarnLevel)
	case slog.LevelError:
		return a.logger.Desugar().Core().Enabled(zap.ErrorLevel)
	default:
		return false
	}
}

// Handle logs the record using the adapted zap.SugaredLogger
func (a *LoggerAdapter) Handle(_ context.Context, record slog.Record) error {
	attrs := convertAttrs(record)
	switch record.Level {
	case slog.LevelDebug:
		a.logger.Debugw(record.Message, attrs...)
	case slog.LevelInfo:
		a.logger.Infow(record.Message, attrs...)
	case slog.LevelWarn:
		a.logger.Warnw(record.Message, attrs...)
	case slog.LevelError:
		a.logger.Errorw(record.Message, attrs...)
	default:
		return nil
	}
	return nil
}

func (a *LoggerAdapter) WithAttrs(attrs []slog.Attr) slog.Handler {
	newLogger := a.logger.With(convertAttrsFromSlice(attrs)...)
	return &LoggerAdapter{logger: newLogger}
}

func (a *LoggerAdapter) WithGroup(name string) slog.Handler {
	newLogger := a.logger.Named(name)
	return &LoggerAdapter{logger: newLogger}
}

func convertAttrs(record slog.Record) []interface{} {
	var kvs []interface{}
	record.Attrs(func(attr slog.Attr) bool {
		kvs = append(kvs, attr.Key, attr.Value)
		return true
	})
	return kvs
}

func convertAttrsFromSlice(attrs []slog.Attr) []interface{} {
	var kvs []interface{}
	for _, attr := range attrs {
		kvs = append(kvs, attr.Key, attr.Value)
	}
	return kvs
}
