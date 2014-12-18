package server

import (
	"encoding/json"
	"fmt"
	"github.com/fukata/golang-stats-api-handler"
	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"
	"github.com/zenazn/goji/web/middleware"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"runtime"
)

type Server struct {
	Config *Config

	Store              *Store
	Worker             []Worker
	UploadFileContexts chan *UploadFileContext
}

type Config struct {
	MongoHost       string
	MongoDataBase   string
	MongoCollection string
	LogFilePath     string
}

func NewServer(Config *Config) *Server {
	server := &Server{
		Config: Config,
	}

	server.Store, _ = NewStore(Config.MongoHost, Config.MongoDataBase, Config.MongoCollection)

	return server
}

func (server *Server) postImage(c web.C, w http.ResponseWriter, r *http.Request) {

	file, header, err := r.FormFile("image")
	defer file.Close()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, err)
		return
	}

	var data []byte
	data, err = ioutil.ReadAll(file)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, err)
		return
	}

	server.UploadFileContexts <- &UploadFileContext{Buf: data, Filename: header.Filename}

	w.WriteHeader(http.StatusOK)
}

func (server *Server) list(c web.C, w http.ResponseWriter, r *http.Request) {
	model, _ := server.Store.All()
	j, _ := json.Marshal(model)
	fmt.Fprintf(w, string(j))
}

func (server *Server) index(w http.ResponseWriter, r *http.Request) {

	w.WriteHeader(http.StatusOK)
}

func (server *Server) Run() {

	f, _ := os.OpenFile(server.Config.LogFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	defer f.Close()

	log.SetOutput(f)

	goji.Use(middleware.Recoverer)
	goji.Use(middleware.NoCache)

	goji.Get("/", server.index)
	goji.Get("/list", server.list)
	goji.Get("/stats", stats_api.Handler)

	goji.Post("/images", server.postImage)

	server.UploadFileContexts = make(chan *UploadFileContext, runtime.NumCPU())
	server.Worker = make([]Worker, runtime.NumCPU())
	for _, v := range server.Worker {
		v.Task = server.UploadFileContexts
		v.Dial = server.Config.MongoHost
		v.DBName = server.Config.MongoDataBase
		v.Collection = server.Config.MongoCollection
		go v.Run()
	}

	goji.Serve()

}
