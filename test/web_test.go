package test

import (
    web2 "github.com/thewindear/go-web-framework/web"
    "log"
    "testing"
)

func TestNewWeb(t *testing.T) {
    web, _ := web2.NewWeb(cfg.Framework, web2.NewLog(cfg.Framework))
    log.Fatalln(web.Listen(cfg.Framework.Web.ServerAddr))
}
