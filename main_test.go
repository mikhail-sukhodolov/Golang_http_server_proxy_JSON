package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const cType = "application/json; charset=utf-8"

func TestCreate(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(Create))
	defer s.Close()
	var str User = User{Name: "Ann", Age: 23}
	k, _ := json.Marshal(str)
	post, err := http.Post(s.URL+"/create", cType, strings.NewReader(string(k)))
	if err != nil {
		t.Fatal(err)
	}
	var f *User
	json.Unmarshal(k, &f)

	z := S.users[0]

	if z.Name != f.Name {
		t.Errorf("Запрос не соответствует")
	}
	if z.Age != f.Age {
		t.Errorf("Запрос не соответствует")
	}
	if post.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK; got %v", post.StatusCode)
	}
}
