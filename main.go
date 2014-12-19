package main

import (
	"github.com/masayukioguni/go-lgtm-backend/backend"
)

func main() {
	s := backend.NewServer(&backend.Config{
		LogFilePath: "./log",
	})
	s.Run()
}
