package configs

// DBConfig ...
type DBConfig struct {
	URL string `json:"url"`
}

// ServerConfig ...
type ServerConfig struct {
	BindAddr string `json:"bind_addr"`
	Name     string `json:"name"`
	Version  string `json:"version"`
}

// Config ...
type Config struct {
	Database *DBConfig     `json:"db"`
	Server   *ServerConfig `json:"server"`
}

// NewConfig ...
func NewConfig(name, version string) *Config {
	return &Config{
		Database: &DBConfig{},
		Server: &ServerConfig{
			Name:    name,
			Version: version,
		},
	}
}
