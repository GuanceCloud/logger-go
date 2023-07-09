package logger

type option func(*loggerOptions)

// A Level is a logging priority. Higher levels are more important.
type Level int8

const (
	// same to zap's levels, but only keep 4 levels of them
	DebugLevel Level = iota - 1
	// InfoLevel is the default logging priority.
	InfoLevel
	// WarnLevel logs are more important than Info, but don't need individual
	// human review.
	WarnLevel
	// ErrorLevel logs are high-priority. If an application is running smoothly,
	// it shouldn't generate any error-level logs.
	ErrorLevel
)

type loggerOptions struct {
	level,
	stackTraceLevel Level

	mode Mode

	logPath,
	errorLogPath string

	stackTrace bool
	stdout     bool
	json       bool
	colorLevel bool

	// lumberjack rotate settings
	rotateEnabled   bool
	rotateMaxSize   int  // default 32MB
	rotateMaxBackup int  // default 5
	rotateMaxAge    int  // default 30day
	rotateCompress  bool // default true
}

// WithJSONEncoding encode log format in JSON.
func WithJSONEncoding(on bool) option {
	return func(o *loggerOptions) {
		o.json = on
	}
}

// WithStackTrace enable stack trace that log level >= lvl
func WithStackTrace(on bool, lvl Level) option {
	return func(o *loggerOptions) {
		o.stackTraceLevel = lvl
		o.stackTrace = on
	}
}

// A Mode is a logging settings. development mode logging are more
// readable than production mode.
type Mode int8

const (
	// ModeDevelopment set mode to development
	ModeDevelopment Mode = iota

	// ModeDevelopment set mode to production
	ModeProduction
)

// WithMode set global logging mode.
func WithMode(mode Mode) option {
	return func(o *loggerOptions) {
		o.mode = mode
	}
}

// WithStdout send log to stdout even the log path already set.
func WithStdout(on bool) option {
	return func(o *loggerOptions) {
		o.stdout = on
	}
}

// WithRotate setup log file rotate on default settings.
func WithRotate(on bool) option {
	return func(o *loggerOptions) {
		o.rotateEnabled = on
		if on {
			o.rotateMaxSize = 32 * 1024 * 1024
			o.rotateMaxBackup = 5
			o.rotateMaxAge = 30
			o.rotateCompress = true
		}
	}
}

// WithErrorLogPath set log file on ERROR logs
func WithErrorLogPath(p string) option {
	return func(o *loggerOptions) {
		o.errorLogPath = p
	}
}

// WithPath set log path.
func WithPath(p string) option {
	return func(o *loggerOptions) {
		o.logPath = p
	}
}

// WithColorLevel set color on log levels. Only applied to
// stdout and non-JSON encoder.
//
// Why: color in json and diskfile makes it hard to parse for
//      logging system, it's only for human readablility.
func WithColorLevel(on bool) option {
	return func(o *loggerOptions) {
		o.colorLevel = on
	}
}

// WithLevel set root logger level.
func WithLevel(lvl Level) option {
	return func(o *loggerOptions) {
		o.level = lvl
	}
}
