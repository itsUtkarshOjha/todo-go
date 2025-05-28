package main

import (
	"log"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	PORT        string
	SCYLLA_HOST string
	KEYSPACE    string
	REDIS_HOST  string
}

func LoadConfig() Config {
	viper.AddConfigPath("./")
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("ReadInConfig error: %v", err)
	}

	config := Config{
		PORT:        viper.GetString("PORT"),
		SCYLLA_HOST: viper.GetString("SCYLLA_HOST"),
		KEYSPACE:    viper.GetString("KEYSPACE"),
		REDIS_HOST:  viper.GetString("REDIS_HOST"),
	}
	return config
}
