package crud

type Config struct {
	ProtectedMethods map[string]bool
}

type Option func(*Config)

func DefaultConfig() *Config {
	return &Config{ProtectedMethods: make(map[string]bool)}
}

func Protect(methods ...string) Option {
	return func(c *Config) {
		for _, m := range methods {
			c.ProtectedMethods[m] = true
		}
	}
}

func ProtectAll() Option {
	return Protect("GET", "POST", "PUT", "PATCH", "DELETE")
}
