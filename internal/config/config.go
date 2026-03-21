package config

type ConfigProvider interface {
	GetPort() int
	GetEnv() string
	GetDBDSN() string
	GetJWTSecret() string
}

type Config struct {
	Port int
	Env  string
	DB   struct {
		DSN string
	}
	Auth struct {
		JWTSecret string
	}
}

func New(port int, env string, dsn string) *Config {
	cfg := &Config{
		Port: port,
		Env:  env,
	}
	cfg.DB.DSN = dsn
	return cfg
}

func (c *Config) GetPort() int {
	return c.Port
}

func (c *Config) GetEnv() string {
	return c.Env
}

func (c *Config) GetDBDSN() string {
	return c.DB.DSN
}

func (c *Config) GetJWTSecret() string {
	return c.Auth.JWTSecret
}
