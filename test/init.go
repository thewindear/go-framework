package test

import (
    "github.com/thewindear/go-web-framework/components"
    "github.com/thewindear/go-web-framework/etc"
)

var cfg = &etc.Cfg{
    Framework: &etc.Framework{
        Log:   logCfg,
        Mysql: mysqlCfg,
        Web:   webCfg,
    },
}

var TLog = components.NewLog(cfg.Framework)

var logCfg = &etc.LogConfig{
    Level:      "info",
    FileName:   "",
    MaxSize:    0,
    MaxAge:     0,
    MaxBackups: 0,
    Compress:   false,
}

var mysqlCfg = &etc.MysqlConfig{
    Host:         "localhost",
    Port:         3306,
    Username:     "root",
    Password:     "Kb7DPGVY98Dv64S97M73gW7GKZjCusje",
    Database:     "test",
    Params:       "charset=utf8mb4&parseTime=True&loc=Local",
    Idle:         10,
    IdleLeftTime: 36400,
    MaxConn:      100,
    LeftTime:     60 * 60 * 2,
    Log:          true,
}

var webCfg = &etc.WebConfig{
    Env:            "dev",
    ServerAddr:     ":8080",
    AppName:        "blog-system",
    DomainName:     "blog.com",
    MaxConcurrency: 30000,
    RequestLimiter: &etc.RequestLimiter{
        Max:        20,
        Expiration: 3600,
    },
    RequestID: &etc.RequestID{
        HeaderName: "X-Request-ID",
    },
    CtxFields: []string{"requestId"},
    RequestLog: &etc.RequestLog{
        Fields: []string{"latency", "requestId", "status", "method", "url", "queryParams"},
    },
}
