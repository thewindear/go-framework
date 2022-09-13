package test

import (
    "github.com/gofiber/fiber/v2"
    goWebFramework "github.com/thewindear/go-web-framework"
    "github.com/thewindear/go-web-framework/etc"
    "go.uber.org/zap"
    "log"
    "testing"
    "time"
)

func TestInitCfg(t *testing.T) {
    c, err := goWebFramework.DefaultInitCfg("../config.template.yaml")
    if err != nil {
        t.Fatal(err)
    } else {
        t.Log(c)
    }
}

func TestFramework(t *testing.T) {
    framework, err := goWebFramework.NewFramework("../config.template.yaml", nil)
    if err != nil {
        t.Fatal(err)
    } else {
        framework.SetRouter(func(app *fiber.App) {
            app.Get("/", func(ctx *fiber.Ctx) error {
                l := framework.GetLog(ctx.Context())
                l.Info("hello world")
                if ctx.Query("key", "") != "" {
                    time.Sleep(time.Second * 3)
                }
                var blog Blog
                framework.GetDB(ctx.Context()).Model(&blog).Where("id = 4").First(&blog)
                l.Info("blog info ", zap.String("title", blog.Title))
                return ctx.JSON(blog)
            })
        })
        framework.Run()
    }
}

type CustomConfig struct {
    Username      string `json:"username"`
    etc.Framework `json:"framework"`
}

func TestCustomConfig(t *testing.T) {
    var customConfig CustomConfig
    _ = goWebFramework.InitCfg[CustomConfig]("../config.template.yaml", &customConfig)
    log.Println(customConfig)
}