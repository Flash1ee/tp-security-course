package cfg

type Config struct {
	BindAddr    string `toml:"bind_addr"`
	LogLevel    string `toml:"log_level"`
	LogAddr     string `toml:"log_path"`
	CertKeyPath string `toml:"ssl_cert_key"`
	DatabaseURL string `toml:"database_url"`
}

func NewConfig() *Config {
	return &Config{}
}
