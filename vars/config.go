package vars

import "fmt"

type config struct {
	BotConfig botConfig `yaml:"botconfig"`
	Sandbox   bool
	Redis     redisConfig `yaml:"redis"`
}

type botConfig struct {
	AppId string `yaml:"appid"`
	Token string `yaml:"token"`
}

type redisConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	DB       int    `yaml:"db"`
	Password string `yaml:"password"`
}

func (c *config) GetToken() string {
	return fmt.Sprintf("%s%s.%s", "Bot ", c.BotConfig.AppId, c.BotConfig.Token)
}

func (c *config) GetRedisAddr() string {
	return fmt.Sprintf("%s:%d", c.Redis.Host, c.Redis.Port)
}

var Config = &config{}
