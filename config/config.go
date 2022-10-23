package config

type Cfg struct {
	Application map[string]interface{} `yaml:"application"`
	Framework   *Framework             `yaml:"framework"`
}

func (im *Cfg) GetAppCfg(key string, def interface{}) interface{} {
	if v, has := im.Application[key]; has {
		return v
	} else {
		return def
	}
}

type Framework struct {
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

func (im *Framework) KeyExists(name string) bool {
	_, ok := im.Keys[name]
	return ok
}

func (im *Framework) ClientKeyExists(name string) bool {
	_, ok := im.ClientKeys[name]
	return ok
}

func (im *Framework) GetClientKeysByName(name string) *ClientKey {
	if key, ok := im.ClientKeys[name]; ok {
		return key
	}
	return nil
}
