package test

import (
    "github.com/thewindear/go-web-framework/web"
    "go.uber.org/zap"
    "testing"
)

func TestLog(t *testing.T) {
    log := web.NewLog(cfg)
    log.Info("hello world", zap.String("username", "root"))
    log.Error("abcd")
}
