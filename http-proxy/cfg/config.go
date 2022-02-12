package cfg

type Config struct {
	BindAddr string `toml:"bind_addr"`
	LogLevel string `toml:"log_level"`
	LogAddr  string `toml:"log_path"`
}

func NewConfig() *Config {
	return &Config{}
}
