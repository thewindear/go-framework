package test

import (
    goWebFramework "github.com/thewindear/go-web-framework"
    "testing"
)

func TestInitCfg(t *testing.T) {
    cfg, err := goWebFramework.IniCfg("../config.template.yaml")
    if err != nil {
        t.Fatal(err)
    } else {
        t.Log(cfg.Web.AppName)
    }
}

func TestFramework(t *testing.T) {
    framework, err := goWebFramework.NewFramework("../config.template.yaml")
    if err != nil {
        t.Fatal(err)
    } else {
        framework.Run()
    }
}
