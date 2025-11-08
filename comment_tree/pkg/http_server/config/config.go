package config

import "time"

type ServerConfig struct {
	Addr            string        `yaml:"addr" env_default:"localhost:8080"`
	MaxReadTimeout  time.Duration `yaml:"max_read_timeout" env_default:"10s"`
	MaxWriteTimeout time.Duration `yaml:"max_write_timeout" env_default:"10s"`
}
