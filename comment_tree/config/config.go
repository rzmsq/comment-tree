package config

import (
	"CommentTree/comment_tree/pkg/http_server/config"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	AppConfig    appConfig           `yaml:"app_config" required:"true"`
	ServerConfig config.ServerConfig `yaml:"server_config" required:"true"`
	DBConfig     DBConfig            `yaml:"db_config" required:"true"`
}

func MustLoadConfig(pathConfig string) *Config {
	cfg := &Config{}
	if err := cleanenv.ReadConfig(pathConfig, cfg); err != nil {
		panic(err)
	}
	return cfg
}
