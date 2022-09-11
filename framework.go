package go_web_framework

import (
    "errors"
    "fmt"
    "github.com/go-redis/redis/v8"
    "github.com/gofiber/fiber/v2"
    "github.com/thewindear/go-web-framework/dao"
    "github.com/thewindear/go-web-framework/etc"
    "github.com/thewindear/go-web-framework/web"
    "go.uber.org/zap"
    "gopkg.in/yaml.v3"
    "gorm.io/gorm"
    "log"
    "os"
)

type Framework struct {
    Cfg   *etc.Cfg
    Web   *fiber.App
    Mysql *gorm.DB
    Rdb   *redis.Client
    Log   *zap.Logger
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

func IniCfg(cfgFile string) (*etc.Cfg, error) {
    content, err := os.ReadFile(cfgFile)
    if err != nil {
        return nil, fmt.Errorf("open config file fialure: %s", err)
    }
    var cfg etc.Cfg
    if err = yaml.Unmarshal(content, &cfg); err != nil {
        return nil, fmt.Errorf("initialize config failure: %s", err)
    }
    return &cfg, nil
}

func NewFramework(cfgFile string) (*Framework, error) {
    cfg, err := IniCfg(cfgFile)
    if err != nil {
        return nil, err
    }
    if cfg.Web == nil {
        return nil, errors.New("web config is empty")
    }
    framework := new(Framework)
    if cfg.Log != nil {
        framework.Log = web.NewLog(cfg)
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
