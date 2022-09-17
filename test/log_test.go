package test

import (
    "github.com/thewindear/go-web-framework/components"
    "go.uber.org/zap"
    "testing"
)

func TestLog(t *testing.T) {
    log2 := components.NewLog(cfg.Framework)
    log2.Info("hello world", zap.String("username", "root"))
    log2.Error("abcd")
}
