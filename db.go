package main

import (
	"fmt"
	"log"

	"github.com/gocql/gocql"
)

func DatabaseConnection() *gocql.Session {
	config := LoadConfig()
	cluster := gocql.NewCluster(config.SCYLLA_HOST)
	cluster.Keyspace = config.KEYSPACE
	cluster.Consistency = gocql.Quorum

	session, err := cluster.CreateSession()
	if err != nil {
		log.Fatalf("Failed to connect to the database, %v", err)
	}
	fmt.Println("Connected to ScyllaDB.")
	return session
}
