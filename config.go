package main

import (
	"os"
	"strconv"
)

type Config struct {
	PORT            string
	SCYLLA_HOST     string
	SCYLLA_PORT     int
	SCYLLA_USERNAME string
	SCYLLA_PASSWORD string
	KEYSPACE        string
	REDIS_HOST      string
	REDIS_PASSWORD  string
}

func LoadConfig() Config {
	// viper.AddConfigPath("./")
	// viper.SetConfigFile(".env")
	// viper.SetConfigType("env")

	// viper.AutomaticEnv()

	// viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// if err := viper.ReadInConfig(); err != nil {
	// 	log.Printf("ReadInConfig error: %v", err)
	// }
	scyllaPort, _ := strconv.Atoi(os.Getenv("SCYLLA_PORT"))

	config := Config{
		PORT:            os.Getenv("PORT"),
		SCYLLA_HOST:     os.Getenv("SCYLLA_HOST"),
		SCYLLA_PORT:     scyllaPort,
		SCYLLA_USERNAME: os.Getenv("SCYLLA_USERNAME"),
		SCYLLA_PASSWORD: os.Getenv("SCYLLA_PASSWORD"),
		KEYSPACE:        os.Getenv("KEYSPACE"),
		REDIS_HOST:      os.Getenv("REDIS_HOST"),
		REDIS_PASSWORD:  os.Getenv("REDIS_PASSWORD"),
	}
	return config
}
