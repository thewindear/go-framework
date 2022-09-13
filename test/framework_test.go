package test

import (
    goWebFramework "github.com/thewindear/go-web-framework"
    "github.com/thewindear/go-web-framework/etc"
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
        framework.Run()
    }
}

type ExtendCfg struct {
    Username string `yaml:"username"`
    etc.Cfg  `yaml:"framework"`
}

func TestExtendCfg(t *testing.T) {
    var extendCfg ExtendCfg
    err := goWebFramework.InitCfg("../config.template.yaml", &extendCfg)
    if err != nil {
        t.Fatal(err)
    } else {
        t.Log(extendCfg.Framework)
    }
}