package components

import (
    "crypto/md5"
    "fmt"
    gonanoid "github.com/matoous/go-nanoid"
    "io"
    "strings"
)

const (
    RandStrLevelEasy = iota
    RandStrLevelNormal
    RandStrLevelHard
)

var Utils = new(utils)

type utils struct{}

func (im utils) CryptMD5(args ...string) string {
    w := md5.New()
    for _, arg := range args {
        _, _ = io.WriteString(w, arg)
    }
    return fmt.Sprintf("%x", w.Sum(nil))
}

// GetEmailUsername 获取邮箱用户名
func (im utils) GetEmailUsername(email string) string {
    usernameAndDomain := strings.Split(email, "@")
    return usernameAndDomain[0]
}

// RandomStr 随机生成指定等级和长度的字符串
func (im utils) RandomStr(level, size int) (randStr string) {
    seed := "0123456789qazxswedcvfrtgbnhyyujmkiolp"
    if level >= RandStrLevelNormal {
        seed += "QAZXSWEDCVFRTGBNHYYUJMKIOLP-"
    }
    if level >= RandStrLevelHard {
        seed += "!@#$%^&*()_+"
    }
    randStr, _ = gonanoid.Generate(seed, size)
    return randStr
}
