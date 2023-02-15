package main

import (
	"proxy"
	chi "github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"os/signal"
	"syscall"
)

func main() {

	go func() {
		ReplicaOne()
	}()
	go func() {
		ReplicaTwo()
	}()
	go func() {
		proxy.ProxyTwoReplicasRun()
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

}
func ReplicaOne() {
	go func() {
		log.Fatalln(http.ListenAndServe(":8080", Handler()))
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

}
func ReplicaTwo() {
	go func() {
		log.Fatalln(http.ListenAndServe(":8081", Handler()))
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit
}

func Handler() http.Handler {
	r := chi.NewRouter()

	r.Post("/create", Create)
	r.Post("/makefriends", MakeFriends)
	r.Delete("/delete", DeleteUser)
	r.Get("/friends/{id}", GetFriends)
	r.Put("/set_age/{id}", SetAge)
	return r

}
