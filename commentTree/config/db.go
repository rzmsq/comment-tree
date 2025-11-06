package config

type DBConfig struct {
	MaxOpenConns int    `yaml:"max_open_conns" env_default:"10"`
	MaxIdleConns int    `yaml:"max_idle_conns" env_default:"5"`
	MasterDSN    string `yaml:"master_dsn" env_default:"master.db"`
}
