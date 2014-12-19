package backend

import (
	"bytes"
	"github.com/masayukioguni/go-lgtm-model"
	"github.com/mitchellh/goamz/aws"
	"github.com/mitchellh/goamz/s3"
	"github.com/t-k/fluent-logger-golang/fluent"
	"io"
	"log"
	"os"
	"path/filepath"
)

const (
	defaultFluentHost = "127.0.0.1"
	defaultFluentPort = 24224
	defaultTagName    = "lgtm.image"
)

type UploadFileContext struct {
	Filename string
	Buf      []byte
}

type Worker struct {
	Task chan *UploadFileContext
}

func (w *Worker) Upload(name string, context *UploadFileContext) error {

	out, err := os.Create("/tmp/test/" + name)
	defer out.Close()

	if err != nil {
		return err
	}

	_, err = io.Copy(out, bytes.NewBuffer(context.Buf))
	if err != nil {
		return err
	}

	return nil
}

func (w *Worker) S3Upload(name string, context *UploadFileContext) error {

	auth, err := aws.EnvAuth()
	if err != nil {
		log.Fatalf("aws.EnvAuth() Failed %s \n", err)
		return err
	}
	client := s3.New(auth, aws.USEast)

	bucket := client.Bucket("go-lgtm")

	err = bucket.Put(name, context.Buf, "image/jpeg", s3.BucketOwnerFull)

	if err != nil {
		log.Fatalf("bucket.Put Failed %s \n", err)
		return err
	}

	return nil
}

func (w *Worker) Run() {
	logger, err := fluent.New(fluent.Config{
		FluentPort: defaultFluentPort,
		FluentHost: defaultFluentHost})

	if err != nil {
		log.Fatalf("w.Upload Failed %s \n", err)
	}

	defer logger.Close()

	for {
		select {
		case job := <-w.Task:
			name := GetRandomName(filepath.Ext(job.Filename)) + filepath.Ext(job.Filename)

			err = w.S3Upload(name, job)
			if err != nil {
				log.Fatalf("w.S3Upload Failed %s \n", err)
			}

			item := &model.Image{Name: name}
			log.Printf("%v\n", item)

			logger.Post(defaultTagName, item)
		}
	}
}
