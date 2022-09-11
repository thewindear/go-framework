package etc

import (
    "fmt"
)

type MysqlConfig struct {
    Host         string `yaml:"host"`
    Port         uint16 `yaml:"port"`
    Username     string `yaml:"username"`
    Password     string `yaml:"password"`
    Database     string `yaml:"database"`
    Params       string `yaml:"params"`
    Idle         int    `yaml:"idle"`
    IdleLeftTime int    `yaml:"idleLeftTime"`
    MaxConn      int    `yaml:"maxConn"`
    LeftTime     int    `yaml:"leftTime"`
    Log          bool   `yaml:"log"`
    LogLevel     string `yaml:"logLevel"`
}

func (im *MysqlConfig) GenDSN() string {
    dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s", im.Username, im.Password, im.Host, im.Port, im.Database, im.Params)
    return dsn
}
