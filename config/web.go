package config

const (
    EnvTest = "test"
    EnvProd = "prod"
    EnvDev  = "dev"
)

type WebConfig struct {
    Env            string          `yaml:"env"`
    ServerAddr     string          `yaml:"serverAddr"`
    AppName        string          `yaml:"appName"`
    DomainName     string          `yaml:"domainName"`
    MaxConcurrency int             `yaml:"maxConcurrency"`
    RequestLimiter *RequestLimiter `yaml:"requestLimiter"`
    RequestID      *RequestID      `yaml:"requestID"`
    RequestLog     *RequestLog     `yaml:"requestLog"`
    CtxFields      []string        `yaml:"ctxFields"`
}

func (im *WebConfig) GetServerHead() string {
    return im.AppName + "-" + im.Env
}

func (im *WebConfig) EnvAppName() string {
    envName := ""
    if im.AppName != "" {
        envName += im.AppName
    } else {
        envName += "web"
    }
    if im.Env != "" {
        envName += " [" + im.Env + "] "
    }
    return envName
}

type RequestLog struct {
    // Fields 打印字段
    // @see https://docs.gofiber.io/api/middleware/logger
    Fields []string `yaml:"fields"`
}

type RequestID struct {
    HeaderName string `yaml:"headerName"`
}

type RequestLimiter struct {
    Max        int `yaml:"max"`
    Expiration int `yaml:"expiration"`
}

func (im *WebConfig) IsTest() bool {
    return im.Env == EnvTest
}

func (im *WebConfig) IsDev() bool {
    return im.Env == EnvDev
}

func (im *WebConfig) IsProd() bool {
    return im.Env == EnvProd
}
