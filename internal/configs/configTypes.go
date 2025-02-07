package configs

type MuopDBConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type HttpConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type GeminiConfig struct {
	APIKey string `yaml:"api_key"`
}
