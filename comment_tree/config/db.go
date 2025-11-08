package config

type DBConfig struct {
	Addr string `yaml:"addr" env_default:"postgres://admin:admin@localhost:5432/master?sslmode=disable"`
}
