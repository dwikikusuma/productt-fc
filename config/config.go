package config

import (
	"github.com/spf13/viper"
	"log"
)

func LoadConfig() Config {
	var cfg Config
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./files/config")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("failed to load config %s", err)
	}

	err = viper.Unmarshal(&cfg)
	if err != nil {
		log.Fatalf("failed to load config %s", err)
	}

	return cfg
}
