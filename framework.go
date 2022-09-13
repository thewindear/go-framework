package go_web_framework

import (
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
    Web   *fiber.App
    Mysql *gorm.DB
    Rdb   *redis.Client
    Log   *zap.Logger
}

func (im *Framework) SetRouter(routeHandle func(app *fiber.App)) {
    routeHandle(im.Web)
}

func (im *Framework) Run() {
    defer im.shutdown()
    log.Fatalf("listen server failure: %s", im.Web.Listen(im.Cfg.Web.ServerAddr))
}

func (im *Framework) shutdown() {
    sqlDB, _ := im.Mysql.DB()
    var err error
    if err = sqlDB.Close(); err != nil {
        im.Log.Error("close mysql error: " + err.Error())
    }
    if err = im.Rdb.Close(); err != nil {
        im.Log.Error("close redis error: " + err.Error())
    }
    if im.Cfg.Log != nil && im.Cfg.Log.FileName != "" {
        if err = im.Log.Sync(); err != nil {
            im.Log.Error("sync log file error: " + err.Error())
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

func InitCfg[T any](cfgFile string, obj T) error {
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
        framework.Log = log2.NewLog(cfg)
    }
    if cfg.Mysql != nil {
        if framework.Mysql, err = dao.NewMysql(cfg, framework.Log); err != nil {
            return nil, err
        }
    }
    if cfg.Redis != nil {
        if framework.Rdb, err = dao.NewRedis(cfg); err != nil {
            return nil, err
        }
    }
    if framework.Web, err = web.NewWeb(cfg, framework.Log); err != nil {
        return nil, err
    }
    framework.Cfg = cfg
    return framework, nil
}
