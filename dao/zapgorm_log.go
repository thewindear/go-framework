package dao

import (
    "context"
    "errors"
    "go.uber.org/zap"
    "gorm.io/gorm"
    gormlogger "gorm.io/gorm/logger"
    "path/filepath"
    "runtime"
    "strings"
    "time"
)

var loggerCtxFields []string

func InitLoggerCtxFields(fields []string) {
    if len(fields) > 0 {
        loggerCtxFields = append(loggerCtxFields, fields[:]...)
    }
}

type Logger struct {
    ZapLogger                 *zap.Logger
    LogLevel                  gormlogger.LogLevel
    SlowThreshold             time.Duration
    SkipCallerLookup          bool
    IgnoreRecordNotFoundError bool
}

func NewZapGormLog(zapLogger *zap.Logger, level gormlogger.LogLevel) *Logger {
    return &Logger{
        ZapLogger:                 zapLogger,
        LogLevel:                  level,
        SlowThreshold:             100 * time.Millisecond,
        SkipCallerLookup:          false,
        IgnoreRecordNotFoundError: false,
    }
}

func (l *Logger) SetAsDefault() {
    gormlogger.Default = l
}

func (l *Logger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
    return &Logger{
        ZapLogger:                 l.ZapLogger,
        SlowThreshold:             l.SlowThreshold,
        LogLevel:                  level,
        SkipCallerLookup:          l.SkipCallerLookup,
        IgnoreRecordNotFoundError: l.IgnoreRecordNotFoundError,
    }
}

func (l *Logger) Info(ctx context.Context, str string, args ...interface{}) {
    if l.LogLevel < gormlogger.Info {
        return
    }
    l.makeCtxLog(ctx).Sugar().Infof(str, args...)
}

func (l *Logger) Warn(ctx context.Context, str string, args ...interface{}) {
    if l.LogLevel < gormlogger.Warn {
        return
    }
    l.makeCtxLog(ctx).Sugar().Warnf(str, args...)
}

func (l *Logger) Error(ctx context.Context, str string, args ...interface{}) {
    if l.LogLevel < gormlogger.Error {
        return
    }
    l.makeCtxLog(ctx).Sugar().Errorf(str, args...)
}

func (l *Logger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
    if l.LogLevel <= gormlogger.Silent {
        return
    }
    elapsed := time.Since(begin)
    switch {
    case err != nil && l.LogLevel >= gormlogger.Error:
        sql, rows := fc()
        if errors.Is(err, gorm.ErrRecordNotFound) {
            l.makeCtxLog(ctx).Info("gorm sql trace", zap.Error(err), zap.Duration("elapsed", elapsed), zap.Int64("rows", rows), zap.String("sql", sql))
        } else {
            l.makeCtxLog(ctx).Error("gorm sql trace", zap.Error(err), zap.Duration("elapsed", elapsed), zap.Int64("rows", rows), zap.String("sql", sql))
        }
    case l.SlowThreshold != 0 && elapsed > l.SlowThreshold && l.LogLevel >= gormlogger.Warn:
        sql, rows := fc()
        l.makeCtxLog(ctx).Warn("gorm sql trace", zap.Bool("slowSql", true), zap.Duration("elapsed", elapsed), zap.Int64("rows", rows), zap.String("sql", sql))
    case l.LogLevel >= gormlogger.Info:
        sql, rows := fc()
        l.makeCtxLog(ctx).Info("gorm sql trace", zap.Duration("elapsed", elapsed), zap.Int64("rows", rows), zap.String("sql", sql))
    }
}

func (l *Logger) makeCtxLog(ctx context.Context) *zap.Logger {
    if len(loggerCtxFields) > 0 {
        var fields []zap.Field
        for _, key := range loggerCtxFields {
            if val := ctx.Value(key); val != nil {
                fields = append(fields, zap.String(key, val.(string)))
            }
        }
        return l.logger().With(fields...)
    }
    return l.logger()
}

var (
    gormPackage = filepath.Join("gorm.io", "gorm")
)

func (l *Logger) logger() *zap.Logger {
    for i := 2; i < 15; i++ {
        _, file, _, ok := runtime.Caller(i)
        switch {
        case !ok:
        case strings.HasSuffix(file, "_test.go"):
        case strings.Contains(file, gormPackage):
        default:
            return l.ZapLogger.WithOptions(zap.AddCallerSkip(i))
        }
    }
    return l.ZapLogger
}
