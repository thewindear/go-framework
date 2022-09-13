package test

import (
    "github.com/gofiber/fiber/v2"
    goWebFramework "github.com/thewindear/go-web-framework"
    "github.com/thewindear/go-web-framework/etc"
    "log"
    "testing"
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
                framework.GetLog(ctx.Context()).Info("hello world")
                return ctx.Send([]byte("hello world"))
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