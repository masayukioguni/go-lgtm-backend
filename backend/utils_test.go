package backend

import (
	"reflect"
	"testing"
)

func TestGetRandomName(t *testing.T) {
	name := GetRandomName("a")

	wantName := "0cc175b9c0f1b6a831c399e269772661"
	if !reflect.DeepEqual(name, wantName) {
		t.Errorf("TestGetRandomName returned %+v, want %+v", name, wantName)
	}

	name = GetRandomName("b")
	wantName = "92eb5ffee6ae2fec3ad71c777531578f"
	if !reflect.DeepEqual(name, wantName) {
		t.Errorf("TestGetRandomName returned %+v, want %+v", name, wantName)
	}

}
