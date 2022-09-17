package test

import (
    "context"
    "github.com/thewindear/go-web-framework/components"
    "go.uber.org/zap"
    "testing"
)

type Blog struct {
    ID     uint   `gorm:"id" json:"ID"`
    ShopId uint   `gorm:"shop_id" json:"shopId"`
    Title  string `gorm:"title" json:"title"`
    Images string `gorm:"images" json:"images"`
}

func (im Blog) TableName() string {
    return "tb_blog"
}

func TestNewMysql(t *testing.T) {
    db, err := components.NewMysql(cfg.Framework, components.NewLog(cfg.Framework))
    if err != nil {
        TLog.Error("connect db error", zap.String("DBError", err.Error()))
    }
    var blog Blog
    ctx := context.WithValue(context.Background(), "requestId", "wefe2f-23f32f23-f23f32-fewd")
    db.WithContext(ctx).Model(&blog).Where("id = 4").First(&blog)
    TLog.Info(blog.Title)
}
