package etc

type Cfg struct {
    Web        *WebConfig            `yaml:"web"`
    Mysql      *MysqlConfig          `yaml:"mysql"`
    Redis      *RedisConfig          `yaml:"redis"`
    Log        *LogConfig            `yaml:"log"`
    Keys       map[string]string     `yaml:"keys"`
    ClientKeys map[string]*ClientKey `yaml:"clientKeys"`
}

type ClientKey struct {
    Name         string `yaml:"name"`
    ClientId     string `yaml:"clientId"`
    ClientSecret string `yaml:"clientSecret"`
}

func (im *Cfg) KeyExists(name string) bool {
    _, ok := im.Keys[name]
    return ok
}

func (im *Cfg) ClientKeyExists(name string) bool {
    _, ok := im.ClientKeys[name]
    return ok
}

func (im *Cfg) GetClientKeysByName(name string) *ClientKey {
    if key, ok := im.ClientKeys[name]; ok {
        return key
    }
    return nil
}
