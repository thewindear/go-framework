package log

import (
    "fmt"
    "github.com/thewindear/go-web-framework/etc"
    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
    "gopkg.in/natefinch/lumberjack.v2"
    "os"
    "time"
)

func NewLog(cfg *etc.Framework) *zap.Logger {
    encoderConfig := zapcore.EncoderConfig{
        MessageKey:     "msg",
        LevelKey:       "level",
        TimeKey:        "time",
        NameKey:        "name",
        SkipLineEnding: false,
        EncodeLevel:    zapcore.CapitalLevelEncoder,
        EncodeTime:     timeEncoderOfLayout(cfg.Web.AppName, cfg.Web.Env),
        StacktraceKey:  "trace",
        LineEnding:     zapcore.DefaultLineEnding,
        EncodeDuration: zapcore.SecondsDurationEncoder,
        EncodeName:     zapcore.FullNameEncoder,
    }
    var writerSync zapcore.WriteSyncer
    if cfg.Log != nil && cfg.Log.FileName != "" && !cfg.Web.IsDev() {
        writer := &lumberjack.Logger{
            Filename:   cfg.Log.FileName,
            MaxSize:    cfg.Log.MaxSize,
            MaxAge:     cfg.Log.MaxAge,
            MaxBackups: cfg.Log.MaxBackups,
            LocalTime:  true,
            Compress:   cfg.Log.Compress,
        }
        writerSync = zapcore.AddSync(writer)
    } else {
        writerSync = zapcore.AddSync(os.Stdout)
    }
    zapCore := zapcore.NewCore(zapcore.NewConsoleEncoder(encoderConfig), writerSync, zap.NewAtomicLevelAt(cfg.Log.GetLogLevel()))
    log := zap.New(zapCore, zap.AddStacktrace(zap.ErrorLevel))
    if cfg.Web.IsDev() {
        log = log.WithOptions(zap.Development())
    }
    zap.ReplaceGlobals(log)
    return log
}

func timeEncoderOfLayout(appName, env string) zapcore.TimeEncoder {
    return func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
        enc.AppendString(time.Now().Format("2006-01-02 15:04:05.000"))
        enc.AppendString(fmt.Sprintf("[%s][%s]", env, appName))
    }
}
