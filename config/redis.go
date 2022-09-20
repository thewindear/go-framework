package config

import "fmt"

type RedisConfig struct {
    Host     string `yaml:"host"`
    Port     uint16 `yaml:"port"`
    Username string `yaml:"username"`
    Password string `yaml:"password"`
    DB       int    `yaml:"DB"`
}

func (im *RedisConfig) GenAddr() string {
    return fmt.Sprintf("%s:%d", im.Host, im.Port)
}

func (im *RedisConfig) IsUsernameValid() bool {
    return len(im.Username) > 0 && im.Username != ""
}

func (im *RedisConfig) IsPasswordValid() bool {
    return len(im.Password) > 0 && im.Password != ""
}
