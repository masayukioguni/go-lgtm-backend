package server

import (
	"labix.org/v2/mgo/bson"
	"reflect"
	"testing"
	"time"
)

const (
	TestPrefix     = "TestPrefix_"
	TestDial       = "mongodb://localhost"
	TestDB         = "test-go-lgtm-server"
	TestCollection = "test_collection"
)

func NewMockStore() *Store {
	s, _ := NewStore(TestDial, TestDB, TestCollection)
	s.Session.DB(TestDB).DropDatabase()
	return s
}

func CloseMockStore(s *Store) {
	_, _ = s.Session.DB(TestDB).C(TestCollection).RemoveAll(bson.M{})
	s.Close()
}

func TestStore(t *testing.T) {
	s := NewMockStore()
	defer CloseMockStore(s)

	createdAt := time.Now()

	insertValue := Image{
		Name:      "image_url",
		CreatedAt: createdAt,
	}

	_ = s.Insert(&insertValue)
	model, _ := s.All()

	wantSize := 1

	if !reflect.DeepEqual(len(model), wantSize) {
		t.Errorf("TestStore returned %+v, want %+v", len(model), wantSize)
	}

	if !reflect.DeepEqual(model[0].Name, insertValue.Name) {
		t.Errorf("TestStore  returned %+v, want %+v", model[0].Name, insertValue.Name)
	}

}

func TestStore_Multi(t *testing.T) {
	s := NewMockStore()
	defer CloseMockStore(s)

	for i := 0; i < 10; i++ {
		insertValue := Image{
			Name:      "image_url",
			CreatedAt: time.Now(),
		}
		_ = s.Insert(&insertValue)
	}

	model, _ := s.All()

	wantSize := 10

	if !reflect.DeepEqual(len(model), wantSize) {
		t.Errorf("TestStore returned %+v, want %+v", len(model), wantSize)
	}
}
