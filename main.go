package main

import (
	"github.com/masayukioguni/go-lgtm-backend/backend"
)

const (
	Dial       = "mongodb://localhost"
	DB         = "test-go-lgtm-server"
	Collection = "test_collection"
)

func main() {
	s := backend.NewServer(&backend.Config{
		MongoHost:       Dial,
		MongoDataBase:   DB,
		MongoCollection: Collection,
		LogFilePath:     "./log",
	})
	s.Run()
}
