package main

import (
	"log"
	"net"
	"net/http"
	"rest-api/internal/user"
	"time"

	"github.com/julienschmidt/httprouter"
)

func main() {
	log.Println("Create New Router")
	router := httprouter.New()

	log.Println("Create User Handler")
	userHandler := user.NewHandler()
	userHandler.Register(router)

	run(router)
}

func run(router *httprouter.Router) {
	listener, err := net.Listen("tcp", ":5000")
	if err != nil {
		panic(err)
	}

	server := &http.Server{
		Handler:      router,
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
	}

	log.Fatalln(server.Serve(listener))
}
