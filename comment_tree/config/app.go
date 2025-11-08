package config

type appConfig struct {
	Env string `yaml:"env" env_default:"development"`
}
