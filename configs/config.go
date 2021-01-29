package configs

type Config struct {
	Server   *ServerConfig   `json:"server"`
	Metric   *MetricConfig   `json:"metric"`
	Amazon   *AmazonConfig   `json:"amazon"`
	Database *DatabaseConfig `json:"database"`
}

func NewConfig(name, version string) *Config {
	return &Config{
		Server: &ServerConfig{
			Name:    name,
			Version: version,
		},
		Metric: &MetricConfig{},
		Amazon: &AmazonConfig{
			DirtyRegion: &AmazonConnectConfig{},
			CleanRegion: &AmazonConnectConfig{},
		},
		Database: &DatabaseConfig{},
	}
}

type ServerConfig struct {
	Addr    string `json:"addr"`
	Name    string `json:"name"`
	Version string `json:"version"`
}

type MetricConfig struct {
	Addr string `json:"addr"`
}

type DatabaseConfig struct {
	URL string `json:"url"`
}

type AmazonConfig struct {
	DirtyRegion *AmazonConnectConfig `json:"dirty_region"`
	CleanRegion *AmazonConnectConfig `json:"clean_region"`
}

type AmazonConnectConfig struct {
	Region    string `json:"region"`
	AccessKey string `json:"access_key"`
	SecretKey string `json:"secret_key"`
	Endpoint  string `json:"endpoint"`
}
