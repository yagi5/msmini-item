package log

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger returns log struct
type Logger struct {
	L *zap.Logger
}

// Field is logging field
type Field struct {
	key string
	val string
}

// F returns field
func F(key, val string) *Field {
	return &Field{key: key, val: val}
}

// New returns zaplogger
func New() (*Logger, error) {
	// really talkative logger.
	l, err := zap.NewProduction(
		zap.AddCaller(),
		zap.AddStacktrace(zapcore.InfoLevel),
		zap.ErrorOutput(zapcore.AddSync(os.Stdout)),
		zap.Fields(zap.Field{Key: "component", Type: zapcore.SkipType, String: "msmini-item"}),
		zap.Hooks(func(ent zapcore.Entry) error {
			return nil
		}),
	)
	return &Logger{l}, err
}

// Info writes out info log
func (l *Logger) Info(msg string, fields ...*Field) {
	fs := []zap.Field{}
	for _, f := range fields {
		zf := zap.String(f.key, f.val)
		fs = append(fs, zf)
	}
	l.L.Info(msg, fs...)
}
