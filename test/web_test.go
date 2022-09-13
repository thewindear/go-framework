package test

import (
    "github.com/thewindear/go-web-framework/log"
    web2 "github.com/thewindear/go-web-framework/web"
    "testing"
)

func TestNewWeb(t *testing.T) {
    web, _ := web2.NewWeb(cfg.Framework, log.NewLog(cfg.Framework))
    t.Fatal(web.Listen(cfg.Framework.Web.ServerAddr))
}
