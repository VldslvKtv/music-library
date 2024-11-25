package logger

import "log/slog"

// Улучшение "читаемости" ошибки
func Err(err error) slog.Attr {
	return slog.Attr{
		Key:   "error",
		Value: slog.StringValue(err.Error()),
	}
}
