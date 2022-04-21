package main

import (
	"net"
	"net/http"
	"rest-api/internal/user"
	"rest-api/pkg/logging"
	"time"

	"github.com/julienschmidt/httprouter"
)

func main() {
	logger := logging.GetLogger()

	logger.Info("Create new router")
	router := httprouter.New()

	logger.Info("Create User handler")
	userHandler := user.NewHandler(logger)
	userHandler.Register(router)

	run(router)
}

func run(router *httprouter.Router) {
	logger := logging.GetLogger()

	listener, err := net.Listen("tcp", ":5000")
	if err != nil {
		panic(err)
	}

	server := &http.Server{
		Handler:      router,
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
	}

	logger.Info("Server started to listen port: 5000")
	logger.Fatalln(server.Serve(listener))
}
