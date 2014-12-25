package backend

import (
	"github.com/zenazn/goji/web"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func NewMockServer() *httptest.Server {
	s := NewServer(&Config{})
	m := web.New()
	s.Prepare(m)

	return httptest.NewServer(m)
}

func TestServer_Index(t *testing.T) {
	ts := NewMockServer()
	defer ts.Close()

	res, err := http.Get(ts.URL + "/")
	if err != nil {
		t.Errorf("server index returned %+v", err)
	}

	if !reflect.DeepEqual(res.StatusCode, http.StatusOK) {
		t.Errorf("server index returned %+v, want %+v", res.StatusCode, http.StatusOK)
	}
}

func TestServer_Stats(t *testing.T) {
	ts := NewMockServer()
	defer ts.Close()

	res, err := http.Get(ts.URL + "/stats")
	if err != nil {
		t.Errorf("server stats returned %+v", err)
	}

	if !reflect.DeepEqual(res.StatusCode, http.StatusOK) {
		t.Errorf("server stats returned %+v, want %+v", res.StatusCode, http.StatusOK)
	}
}
