package test

import (
    "github.com/gofiber/fiber/v2"
    goWebFramework "github.com/thewindear/go-web-framework"
    "github.com/thewindear/go-web-framework/etc"
    "go.uber.org/zap"
    "gorm.io/gorm"
    "log"
    "testing"
)

type UserService struct {
    *goWebFramework.SvcContext
}

func (im UserService) GetBlog() *Blog {
    var blog Blog
    if im.DB().Model(&blog).Where("id = 4").First(&blog).Error == gorm.ErrRecordNotFound {
        return nil
    }
    im.Log().Info("blog info ", zap.String("title", blog.Title))
    return &blog
}

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
                userService := &UserService{goWebFramework.NewDefaultSvcContext(ctx.Context(), framework)}
                blog := userService.GetBlog()
                if blog == nil {
                    return ctx.SendStatus(404)
                }
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