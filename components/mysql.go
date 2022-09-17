package components

import (
    "github.com/thewindear/go-web-framework/etc"
    "go.uber.org/zap"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
    "time"
)

func NewMysql(cfg *etc.Framework, log *zap.Logger) (*gorm.DB, error) {
    gormCfg := &gorm.Config{
        PrepareStmt: true,
        QueryFields: true,
    }
    if cfg.Log != nil && cfg.Mysql.Log {
        zapGormLog := NewZapGormLog(log, cfg.Log.GetGormLogLevel(), cfg.Mysql.SlowSqlTime)
        InitLoggerCtxFields(cfg.Web.CtxFields)
        zapGormLog.SetAsDefault()
        gormCfg.Logger = zapGormLog
    }
    db, err := gorm.Open(mysql.Open(cfg.Mysql.GenDSN()), gormCfg)
    if err != nil {
        return nil, err
    }
    sqlDB, _ := db.DB()
    sqlDB.SetMaxIdleConns(cfg.Mysql.Idle)
    sqlDB.SetConnMaxLifetime(time.Duration(cfg.Mysql.LeftTime) * time.Second)
    sqlDB.SetMaxOpenConns(cfg.Mysql.MaxConn)
    sqlDB.SetConnMaxIdleTime(time.Duration(cfg.Mysql.IdleLeftTime) * time.Second)
    return db, nil
}
