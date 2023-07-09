package logger

import (
	T "testing"

	"go.uber.org/zap"

	"github.com/stretchr/testify/assert"
)

func TestOptions(t *T.T) {
	t.Run("default-logger", func(t *T.T) {
		l := SLogger(t.Name())
		l.Info("info log by default logger")
		l.Debug("debug log by default logger")
		l.Error("erorr log by default logger")
	})

	t.Run("default-send-to-stdout", func(t *T.T) {
		assert.NoError(t, Setup())

		l := SLogger(t.Name())
		l.Info("this is a info log")
		l.Error("this is a error log")
		l.Debug("this is a debug log")
		l.Warn("this is a warn log")
	})

	t.Run("only-log-to-file", func(t *T.T) {
		assert.NoError(t, Setup(
			WithPath("a.log"),
		))

		l := SLogger(t.Name())
		l.Info("this is a info log")
		l.Error("this is a error log")
		l.Debug("this is a debug log")
		l.Warn("this is a warn log")
	})

	t.Run("color", func(t *T.T) {
		assert.NoError(t, Setup(
			WithStdout(true),
			WithColorLevel(true),
		))

		l := SLogger(t.Name())
		l.Info("this is a info log")
		l.Error("this is a error log")
		l.Debug("this is a debug log")
		l.Warn("this is a warn log")
	})

	t.Run("color-in-json", func(t *T.T) {
		assert.NoError(t, Setup(
			WithStdout(true),
			WithColorLevel(true),
			WithJSONEncoding(true),
		))

		l := SLogger(t.Name())
		l.Info("this is a info log")
		l.Error("this is a error log")
		l.Debug("this is a debug log")
		l.Warn("this is a warn log")
	})

	t.Run("high-level", func(t *T.T) {
		assert.NoError(t, Setup(
			WithStdout(true),
			WithColorLevel(true),
			WithLevel(ErrorLevel),
		))

		l := SLogger(t.Name())
		l.Info("this is a info log")
		l.Error("this is a error log")
		l.Debug("- this is a debug log")
		l.Warn("this is a warn log")
	})

	t.Run("debug-level", func(t *T.T) {
		assert.NoError(t, Setup(
			WithStdout(true),
			WithColorLevel(true),
			WithLevel(DebugLevel),
		))

		l := SLogger(t.Name())
		l.Info("this is a info log")
		l.Error("this is a error log")
		l.Debug("- this is a debug log")
		l.Warn("this is a warn log")
	})

	t.Run("de-sugar", func(t *T.T) {
		assert.NoError(t, Setup(
			WithStdout(true),
			WithColorLevel(true),
			WithLevel(DebugLevel),
		))

		l := SLogger(t.Name())
		ds := l.Desugar()
		ds.Info("this is a info log")
		ds.Error("this is a error log")
		ds.Debug("- this is a debug log")
		ds.Warn("this is a warn log")
		ds.Warn("this is a warn log with fields",
			zap.String(`name`, `foo`),
			zap.Any(`any`, l),
			zap.Namespace(t.Name()),
		)
	})
}
