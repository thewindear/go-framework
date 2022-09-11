package test

import (
    web2 "github.com/thewindear/go-web-framework/web"
    "log"
    "testing"
)

func TestNewWeb(t *testing.T) {
    web, _ := web2.NewWeb(cfg, web2.NewLog(cfg))
    log.Fatalln(web.Listen(cfg.Web.ServerAddr))
}
