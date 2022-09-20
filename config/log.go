package config

import (
    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
    gormlogger "gorm.io/gorm/logger"
    "strings"
)

type LogConfig struct {
    Level      string `yaml:"level"`
    FileName   string `yaml:"fileName"`
    MaxSize    int    `yaml:"maxSize"`
    MaxAge     int    `yaml:"maxAge"`
    MaxBackups int    `yaml:"maxBackups"`
    Compress   bool   `yaml:"compress"`
}

func (im *LogConfig) GetGormLogLevel() gormlogger.LogLevel {
    switch strings.ToLower(im.Level) {
    case "info":
        return gormlogger.Info
    case "error":
        return gormlogger.Error
    case "warn":
        return gormlogger.Warn
    default:
        return gormlogger.Silent
    }
}

func (im *LogConfig) GetLogLevel() zapcore.Level {
    switch strings.ToLower(im.Level) {
    case "info":
        return zap.InfoLevel
    case "warn":
        return zap.WarnLevel
    case "error":
        return zap.ErrorLevel
    default:
        return zap.DebugLevel
    }
}
