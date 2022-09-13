package test

import (
    goWebFramework "github.com/thewindear/go-web-framework"
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