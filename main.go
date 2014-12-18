package main

import (
	"github.com/masayukioguni/go-lgtm-backend/server"
)

const (
	Dial       = "mongodb://localhost"
	DB         = "test-go-lgtm-server"
	Collection = "test_collection"
)

func main() {
	s := server.NewServer(&server.Config{
		MongoHost:       Dial,
		MongoDataBase:   DB,
		MongoCollection: Collection,
		LogFilePath:     "./log",
	})
	s.Run()
}
