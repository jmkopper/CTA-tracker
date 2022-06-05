// config.go
package main

import (
	"fmt"
	"os"
	"strconv"
)

type configSettings struct {
	port                int
	apiKey              string
	mongoUsername       string
	mongoPassword       string
	mongoClusterAddress string
	db                  string
	collection          string
}

func readConfig() configSettings {
	portString := os.Getenv("PORT")

	if portString == "" {
		portString = "8000"
	}

	port, err := strconv.Atoi(portString)

	if err != nil {
		panic(fmt.Sprintf("Could not parse %s to int", portString))
	}

	/*
		port: port the app runs on (uses PORT env var if set, else 8000)
		apiKey: CTA API key
		mongoUsername: mongo username
		mongoPassword: mongo password
		mongoClusterAddress: mongo cluster address, usually of the form "<clusterName>.<xyz>.mongodb.net"
		db: name of mongo database to search
		collection: name of mongo collection to search (must be a collection in above db)
	*/

	return configSettings{
		port:                port,
		apiKey:              "<CTA API key>",
		mongoUsername:       "<Mongo username>",
		mongoPassword:       "<Mongo password>",
		mongoClusterAddress: "<Mongo URI>",
		db:                  "<db name>",
		collection:          "<collection name>",
	}
}
