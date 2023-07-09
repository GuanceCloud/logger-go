package logger

import (
	"os"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	zapRoot = defaultRoot()

	mtx sync.Mutex
)

func SLogger(name string) *zap.SugaredLogger {
	return zapRoot.Sugar().Named(name)
}

func defaultRoot() *zap.Logger {
	x, _ := zap.NewDevelopment()
	return x
}

func Setup(opts ...option) error {
	mtx.Lock()
	defer mtx.Unlock()

	o := &loggerOptions{}
	for _, opt := range opts {
		if opt != nil {
			opt(o)
		}
	}

	return o.setup()
}

func (o *loggerOptions) setup() error {
	highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.Level(ErrorLevel)
	})

	lowPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl < zapcore.Level(ErrorLevel)
	})

	allPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.Level(o.level)
	})

	var (
		cores []zapcore.Core

		stdoutEncoder,
		fileEncoder zapcore.Encoder
		encCfg zapcore.EncoderConfig
	)

	switch o.mode {
	case ModeProduction:
		encCfg = zap.NewProductionEncoderConfig()

	default: // ModeDevelopment
		encCfg = zap.NewDevelopmentEncoderConfig()
	}

	if o.colorLevel &&
		!o.json && // not json
		o.logPath == "" { // not diskfile
		encCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	if o.json {
		stdoutEncoder = zapcore.NewJSONEncoder(encCfg)
		fileEncoder = zapcore.NewJSONEncoder(encCfg)
	} else {
		stdoutEncoder = zapcore.NewConsoleEncoder(encCfg)
		fileEncoder = zapcore.NewConsoleEncoder(encCfg)
	}

	// send log to disk file(with or without rotate)
	if o.logPath != "" {
		var fileBasedSyncer zapcore.WriteSyncer
		if o.rotateEnabled {
			fileBasedSyncer = zapcore.AddSync(&lumberjack.Logger{
				Filename:   o.logPath,
				MaxSize:    o.rotateMaxSize,
				MaxBackups: o.rotateMaxBackup,
				MaxAge:     o.rotateMaxAge,
			})
		} else {
			fileBasedSyncer = zapcore.Lock(mustNewFileSyncer(o.logPath))
		}

		if o.errorLogPath != "" {
			// only accept low-lavel logs, some error logs separated to other place
			cores = append(cores,
				zapcore.NewCore(fileEncoder, fileBasedSyncer, lowPriority),
			)
		} else { // accept all-level logs
			cores = append(cores,
				zapcore.NewCore(fileEncoder, fileBasedSyncer, allPriority),
			)
		}
	} else { // default send log to stdout
		o.stdout = true
	}

	if o.errorLogPath != "" {
		var errSyncer zapcore.WriteSyncer

		if o.rotateEnabled {
			errSyncer = zapcore.AddSync(&lumberjack.Logger{
				Filename:   o.errorLogPath,
				MaxSize:    o.rotateMaxSize,
				MaxBackups: o.rotateMaxBackup,
				MaxAge:     o.rotateMaxAge,
			})
		} else {
			errSyncer = zapcore.Lock(mustNewFileSyncer(o.errorLogPath))
		}
		cores = append(cores,
			zapcore.NewCore(fileEncoder, errSyncer, highPriority),
		)
	}

	if o.stdout {
		stdout := zapcore.Lock(os.Stdout)
		// under stdout, both high and low priority accepted
		cores = append(cores,
			zapcore.NewCore(stdoutEncoder, stdout, allPriority),
		)
	}

	core := zapcore.NewTee(cores...)

	zapNewOptions := []zap.Option{
		zap.AddCaller(), // default enabled
	}

	if o.stackTrace {
		zapNewOptions = append(zapNewOptions,
			zap.AddStacktrace(zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
				return lvl >= zapcore.Level(o.stackTraceLevel)
			})))
	}

	zapRoot = zap.New(core, zapNewOptions...)
	return nil
}

func Sync() {
	mtx.Lock()
	defer mtx.Unlock()
	zapRoot.Sync()
}
