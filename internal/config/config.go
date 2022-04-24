package config

import (
	"rest-api/pkg/logging"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	IsDebug bool `yaml:"is_debug"`
	Listen  struct {
		Type   string `yaml:"type"`
		BindIp string `yaml:"bind_ip"`
		Port   string `yaml:"port"`
	} `yaml:"listen"`
	MongoDB struct {
		Username   string `yaml:"username"`
		Password   string `yaml:"password"`
		Database   string `yaml:"database"`
		Collection string `yaml:"collection"`
		AuthDB     string `yaml:"auth_db"`
		Host       string `yaml:"host"`
		Port       string `yaml:"port"`
		AtlasURI   string `yaml:"atlas_uri"`
	} `yaml:"mongodb"`
	SecretCrypto string `yaml:"secret_crypto"`
}

var instance *Config // singleton
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		logger := logging.GetLogger()
		logger.Info("Reading Config...")

		instance = &Config{}
		if err := cleanenv.ReadConfig("config.yml", instance); err != nil {
			help, _ := cleanenv.GetDescription(instance, nil)
			logger.Info(help)
			logger.Fatal(err)
		}
	})

	return instance
}
