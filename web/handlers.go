package web

import (
    "github.com/gofiber/fiber/v2"
    go_web_framework "github.com/thewindear/go-web-framework"
    "go.uber.org/zap"
)

func ErrorHandler(components *go_web_framework.Components) fiber.ErrorHandler {
    return func(ctx *fiber.Ctx, err error) error {
        if err == nil {
            return nil
        }
        var wrapError *RespError
        if tmpErr, ok := err.(*RespError); ok {
            wrapError = tmpErr
        } else {
            wrapError = Error(tmpErr)
        }
        log := components.GetLogWithContext(ctx.Context())
        log.Error("response error", zap.String("error details", wrapError.Error()))
        return ctx.Status(wrapError.HttpStatus).JSON(wrapError)
    }
    
}
