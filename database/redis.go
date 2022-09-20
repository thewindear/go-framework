package database

import (
    "github.com/go-redis/redis/v8"
    "github.com/thewindear/go-web-framework/config"
)

func NewRedis(cfg *config.Framework) (*redis.Client, error) {
    option := &redis.Options{
        Addr: cfg.Redis.GenAddr(),
        DB:   cfg.Redis.DB,
    }
    if cfg.Redis.IsUsernameValid() {
        option.Username = cfg.Redis.Username
    }
    if cfg.Redis.IsPasswordValid() {
        option.Password = cfg.Redis.Password
    }
    rdb := redis.NewClient(option)
    return rdb, nil
}
