package configs

import (
	"github.com/TrungBui59/test_muopdb/config"
	"gopkg.in/yaml.v2"
	"os"
)

type Config struct {
	MuopDBConfig MuopDBConfig `yaml:"muopdb"`
}

func NewConfig(configPath string) (Config, error) {
	var (
		cfg         = Config{}
		configBytes = config.DefaultConfig
		err         error
	)

	if configPath != "" {
		configBytes, err = os.ReadFile(configPath)
		if err != nil {
			return Config{}, err
		}
	}

	err = yaml.Unmarshal(configBytes, &cfg)
	if err != nil {
		return Config{}, err
	}
	return cfg, nil
}
