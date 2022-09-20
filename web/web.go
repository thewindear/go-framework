package web

import (
    "github.com/gofiber/contrib/fiberzap"
    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/limiter"
    "github.com/gofiber/fiber/v2/middleware/requestid"
    "github.com/thewindear/go-web-easy-kit/config"
    "go.uber.org/zap"
    "time"
)

func NewWeb(cfg *config.Framework, logger *zap.Logger) (*fiber.App, error) {
    fiberConfig := fiber.Config{
        ServerHeader: cfg.Web.GetServerHead(),
        Concurrency:  cfg.Web.MaxConcurrency,
        AppName:      cfg.Web.EnvAppName(),
    }
    web := fiber.New(fiberConfig)
    // 请求日志
    if cfg.Web.RequestLog != nil {
        fiberZapConfig := fiberzap.ConfigDefault
        fiberZapConfig.Fields = cfg.Web.RequestLog.Fields
        if logger != nil {
            fiberZapConfig.Logger = logger
        }
        web.Use(fiberzap.New(fiberZapConfig))
    }
    // 是否使用request-id中间件
    if cfg.Web.RequestID != nil {
        requestIDConfig := requestid.ConfigDefault
        requestIDConfig.ContextKey = "requestId"
        requestIDConfig.Header = cfg.Web.RequestID.HeaderName
        web.Use(requestid.New(requestIDConfig))
    }
    
    // 是否使用limiter中间件
    if cfg.Web.RequestLimiter != nil {
        web.Use(limiter.New(
            limiter.Config{
                Max:        cfg.Web.RequestLimiter.Max,
                Expiration: time.Duration(cfg.Web.RequestLimiter.Expiration) * time.Second,
            },
        ))
    }
    return web, nil
}
