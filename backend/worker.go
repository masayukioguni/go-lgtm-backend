package backend

import (
	"bytes"
	"github.com/masayukioguni/go-lgtm-model"
	"github.com/mitchellh/goamz/aws"
	"github.com/mitchellh/goamz/s3"
	"io"
	"log"
	"os"
	"path/filepath"
)

type UploadFileContext struct {
	Filename string
	Buf      []byte
}

type Worker struct {
	Task       chan *UploadFileContext
	Dial       string
	DBName     string
	Collection string
}

func (w *Worker) Insert(image *model.Image) error {
	store, err := model.NewStore(w.Dial, w.DBName, w.Collection)
	defer store.Close()

	if err != nil {
		return err
	}

	err = store.Insert(image)
	if err != nil {
		return err
	}

	return nil
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
	for {
		select {
		case job := <-w.Task:
			name := GetRandomName(filepath.Ext(job.Filename)) + filepath.Ext(job.Filename)

			err := w.Upload(name, job)
			if err != nil {
				log.Fatalf("w.Upload Failed %s \n", err)
			}

			err = w.S3Upload(name, job)
			if err != nil {
				log.Fatalf("w.S3Upload Failed %s \n", err)
			}

			err = w.Insert(&model.Image{Name: name})
			if err != nil {
				log.Fatalf("w.Insert Failed %s \n", err)
			}
		}
	}
}
