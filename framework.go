package go_web_easy_kit

import (
    "context"
    "fmt"
    "github.com/go-redis/redis/v8"
    "github.com/gofiber/fiber/v2"
    "github.com/pkg/errors"
    "github.com/thewindear/go-web-easy-kit/config"
    "github.com/thewindear/go-web-easy-kit/database"
    log2 "github.com/thewindear/go-web-easy-kit/log"
    "github.com/thewindear/go-web-easy-kit/web"
    "go.uber.org/zap"
    "gopkg.in/yaml.v3"
    "gorm.io/gorm"
    "log"
    "os"
    "os/signal"
    "syscall"
)

var ErrWebConfigEmpty = errors.New("web config empty")

func NewDefaultSvcContext(ctx context.Context, components *Components) *SvcContext {
    return &SvcContext{Ctx: ctx, Components: components}
}

type SvcContext struct {
    Ctx context.Context
    *Components
}

func (im *SvcContext) GetOne(model interface{}, where string, values ...interface{}) error {
    err := im.DB().Model(model).Where(where, values...).First(&model).Error
    if err != nil {
        return err
    }
    return nil
}

func (im *SvcContext) DB() *gorm.DB {
    return im.GetDBWithContext(im.Ctx)
}

func (im *SvcContext) Log() *zap.Logger {
    return im.GetLogWithContext(im.Ctx)
}

func (im *SvcContext) RDB() *redis.Client {
    return im.GetRdbWithContext(im.Ctx)
}

type Components struct {
    cfg   *config.Cfg
    mysql *gorm.DB
    rdb   *redis.Client
    log   *zap.Logger
}

func (im *Components) MakeSvc(ctx context.Context) *SvcContext {
    return NewDefaultSvcContext(ctx, im)
}

func (im *Components) GetCfg() *config.Framework {
    return im.cfg.Framework
}

func (im *Components) GetDBWithContext(ctx context.Context) *gorm.DB {
    return im.mysql.WithContext(ctx)
}

func (im *Components) GetLogWithContext(ctx context.Context) *zap.Logger {
    if im.cfg.Framework.Web.CtxFields != nil && len(im.cfg.Framework.Web.CtxFields) > 0 {
        var fields []zap.Field
        for _, key := range im.cfg.Framework.Web.CtxFields {
            val := ctx.Value(key)
            if val != nil {
                fields = append(fields, zap.String(key, fmt.Sprintf("%s", val)))
            }
        }
        return im.log.With(fields...)
    }
    return im.log
}

func (im *Components) GetRdbWithContext(ctx context.Context) *redis.Client {
    return im.rdb.WithContext(ctx)
}

type Framework struct {
    web *fiber.App
    *Components
}

func (im *Framework) GetComponents() *Components {
    return im.Components
}

func (im *Framework) GetWeb() *fiber.App {
    return im.web
}

func (im *Framework) SetHandles(handleEntry func(route fiber.Router, components *Components)) {
    handleEntry(im.web, im.Components)
}

func (im *Framework) Run() {
    go func() {
        err := im.web.Listen(im.cfg.Framework.Web.ServerAddr)
        if err != nil {
            log.Fatalf("listen server failure: %s", err)
        }
    }()
    ch := make(chan os.Signal, 1)
    signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
    _ = <-ch
    log.Println("Gracefully shutting downing...")
    _ = im.web.Shutdown()
    log.Println("Running close task...")
    im.shutdown()
    log.Println("bye bye server shutdown...")
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
    if im.cfg.Framework.Log != nil && im.cfg.Framework.Log.FileName != "" {
        if err = im.log.Sync(); err != nil {
            im.log.Error("sync log file error: " + err.Error())
        }
    }
}

func DefaultInitCfg(cfgFile string) (*config.Cfg, error) {
    var cfg config.Cfg
    err := InitCfg(cfgFile, &cfg)
    if err != nil {
        return nil, err
    }
    return &cfg, nil
}

func InitCfg[T any](cfgFile string, obj *T) error {
    content, err := os.ReadFile(cfgFile)
    if err != nil {
        return fmt.Errorf("framework: open config file fialure: %w", err)
    }
    if err = yaml.Unmarshal(content, &obj); err != nil {
        return fmt.Errorf("framework: initialize config failure: %w", err)
    }
    return nil
}

func NewFramework(cfgFile string) (*Framework, error) {
    var err error
    var cfg config.Cfg
    err = InitCfg[config.Cfg](cfgFile, &cfg)
    
    if cfg.Framework.Web == nil {
        return nil, ErrWebConfigEmpty
    }
    framework := &Framework{Components: &Components{}}
    if cfg.Framework.Log != nil {
        framework.log = log2.NewLog(cfg.Framework)
    }
    if cfg.Framework.Mysql != nil {
        if framework.mysql, err = database.NewMysql(cfg.Framework, framework.log); err != nil {
            return nil, err
        }
    }
    if cfg.Framework.Redis != nil {
        if framework.rdb, err = database.NewRedis(cfg.Framework); err != nil {
            return nil, err
        }
    }
    if framework.web, err = web.NewWeb(cfg.Framework, framework.log); err != nil {
        return nil, err
    }
    framework.cfg = &cfg
    return framework, nil
}

func ErrorHandler(components *Components) fiber.ErrorHandler {
    return func(ctx *fiber.Ctx, err error) error {
        var wrapError *web.RespError
        switch err := err.(type) {
        case nil:
            return nil
        case *web.RespError:
            wrapError = err
            break
        default:
            wrapError = web.Error(err)
        }
        ctxLog := components.GetLogWithContext(ctx.Context())
        if wrapError.HttpStatus >= fiber.StatusInternalServerError {
            errStackInfos := fmt.Sprintf("%+v", wrapError.Err)
            if errStackInfos != wrapError.Err.Error() {
                ctxLog.Info("server error stacks", zap.String("stacks", errStackInfos))
            }
            ctxLog.Error("server error", zap.String("details", wrapError.Error()))
        } else {
            ctxLog.Info("logic error", zap.String("details", wrapError.Error()))
        }
        return ctx.Status(wrapError.HttpStatus).JSON(wrapError)
    }
}
