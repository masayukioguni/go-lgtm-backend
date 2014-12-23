package main

import (
	"github.com/joho/godotenv"
	"github.com/masayukioguni/go-lgtm-backend/backend"
	"log"
	"os"
	"strconv"
)

var applicationName = "go-lgtm-backend"
var version = "0.1"

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	LogPath := os.Getenv("LOG_PATH")

	f, err := os.OpenFile(LogPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Failed to os.OpenFile()")
	}
	defer f.Close()

	FluentHost := os.Getenv("FLUENT_HOST")
	FluentPort, _ := strconv.Atoi(os.Getenv("FLUENT_PORT"))
	FluentTagName := os.Getenv("FLUENT_TAG_NAME")
	S3Bucket := os.Getenv("S3BUCKET")

	log.SetOutput(f)

	s := backend.NewServer(&backend.Config{
		LogFilePath:   LogPath,
		FluentHost:    FluentHost,
		FluentPort:    FluentPort,
		FluentTagName: FluentTagName,
		S3Bucket:      S3Bucket,
	})
	s.Run()
}
