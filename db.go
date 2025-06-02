package main

import (
	"fmt"
	"log"

	"github.com/gocql/gocql"
)

func DatabaseConnection() *gocql.Session {
	config := LoadConfig()
	cluster := gocql.NewCluster(config.SCYLLA_HOST)
	cluster.Port = config.SCYLLA_PORT
	cluster.Keyspace = config.KEYSPACE
	cluster.Consistency = gocql.Quorum
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: config.SCYLLA_USERNAME,
		Password: config.SCYLLA_PASSWORD,
	}

	session, err := cluster.CreateSession()
	if err != nil {
		log.Fatalf("Failed to connect to the database, %v", err)
	}
	fmt.Println("Connected to ScyllaDB.")
	return session
}
