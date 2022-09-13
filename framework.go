package go_web_framework

import (
    "context"
    "errors"
    "fmt"
    "github.com/go-redis/redis/v8"
    "github.com/gofiber/fiber/v2"
    "github.com/thewindear/go-web-framework/dao"
    "github.com/thewindear/go-web-framework/etc"
    log2 "github.com/thewindear/go-web-framework/log"
    "github.com/thewindear/go-web-framework/web"
    "go.uber.org/zap"
    "gopkg.in/yaml.v3"
    "gorm.io/gorm"
    "log"
    "os"
)

type Framework struct {
    Cfg   *etc.Framework
    web   *fiber.App
    mysql *gorm.DB
    rdb   *redis.Client
    log   *zap.Logger
}

func (im *Framework) GetDB(ctx context.Context) *gorm.DB {
    return im.mysql.WithContext(ctx)
}

func (im *Framework) GetLog(ctx context.Context) *zap.Logger {
    if im.Cfg.Web.CtxFields != nil && len(im.Cfg.Web.CtxFields) > 0 {
        var fields []zap.Field
        for _, key := range im.Cfg.Web.CtxFields {
            val := ctx.Value(key)
            if val != nil {
                fields = append(fields, zap.String(key, fmt.Sprintf("%s", val)))
            }
        }
        return im.log.With(fields...)
    }
    return im.log
}

func (im *Framework) GetRdb(ctx context.Context) *redis.Client {
    return im.rdb.WithContext(ctx)
}

func (im *Framework) SetRouter(routeHandle func(app *fiber.App)) {
    routeHandle(im.web)
}

func (im *Framework) Run() {
    defer im.shutdown()
    log.Fatalf("listen server failure: %s", im.web.Listen(im.Cfg.Web.ServerAddr))
}

func (im *Framework) shutdown() {
    sqlDB, _ := im.mysql.DB()
    var err error
    if err = sqlDB.Close(); err != nil {
        im.log.Error("close mysql error: " + err.Error())
    }
    if err = im.rdb.Close(); err != nil {
        im.log.Error("close redis error: " + err.Error())
    }
    if im.Cfg.Log != nil && im.Cfg.Log.FileName != "" {
        if err = im.log.Sync(); err != nil {
            im.log.Error("sync log file error: " + err.Error())
        }
    }
    log.Println("bye bye server shutdown...")
}

func DefaultInitCfg(cfgFile string) (*etc.Cfg, error) {
    var cfg etc.Cfg
    err := InitCfg(cfgFile, &cfg)
    if err != nil {
        return nil, err
    }
    return &cfg, nil
}

func InitCfg[T any](cfgFile string, obj *T) error {
    content, err := os.ReadFile(cfgFile)
    if err != nil {
        return fmt.Errorf("open config file fialure: %s", err)
    }
    if err = yaml.Unmarshal(content, &obj); err != nil {
        return fmt.Errorf("initialize config failure: %s", err)
    }
    return nil
}

func NewFramework(cfgFile string, cfg *etc.Framework) (*Framework, error) {
    var err error
    if cfg == nil {
        _cfg := &etc.Cfg{}
        err = InitCfg(cfgFile, _cfg)
        if err != nil {
            return nil, err
        }
        cfg = _cfg.Framework
    }
    
    if cfg.Web == nil {
        return nil, errors.New("web config is empty")
    }
    framework := new(Framework)
    if cfg.Log != nil {
        framework.log = log2.NewLog(cfg)
    }
    if cfg.Mysql != nil {
        if framework.mysql, err = dao.NewMysql(cfg, framework.log); err != nil {
            return nil, err
        }
    }
    if cfg.Redis != nil {
        if framework.rdb, err = dao.NewRedis(cfg); err != nil {
            return nil, err
        }
    }
    if framework.web, err = web.NewWeb(cfg, framework.log); err != nil {
        return nil, err
    }
    framework.Cfg = cfg
    return framework, nil
}
