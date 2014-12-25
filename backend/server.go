package backend

import (
	"fmt"
	"github.com/fukata/golang-stats-api-handler"
	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"
	"github.com/zenazn/goji/web/middleware"
	"io/ioutil"
	"net/http"
	"runtime"
)

type Server struct {
	Config             *Config
	Worker             []Worker
	UploadFileContexts chan *UploadFileContext
}

type Config struct {
	LogFilePath   string
	FluentHost    string
	FluentPort    int
	FluentTagName string
	S3Bucket      string
}

func NewServer(Config *Config) *Server {
	server := &Server{
		Config: Config,
	}

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

func (server *Server) index(w http.ResponseWriter, r *http.Request) {

	w.WriteHeader(http.StatusOK)
}

func (server *Server) Prepare(m *web.Mux) {
	m.Get("/", server.index)
	m.Get("/stats", stats_api.Handler)
	m.Post("/images", server.postImage)
}

func (server *Server) Run() {

	server.Prepare(goji.DefaultMux)

	goji.Use(middleware.Recoverer)
	goji.Use(middleware.NoCache)

	server.UploadFileContexts = make(chan *UploadFileContext, runtime.NumCPU())
	server.Worker = make([]Worker, runtime.NumCPU())
	for _, v := range server.Worker {
		v.Task = server.UploadFileContexts
		go v.Run()
	}

	goji.Serve()

}
