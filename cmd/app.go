package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"rest-api/internal/config"
	"rest-api/internal/user"
	"rest-api/pkg/logging"
	"time"

	"github.com/julienschmidt/httprouter"
)

func main() {
	logger := logging.GetLogger()

	logger.Info("Create new router")
	router := httprouter.New()

	cfg := config.GetConfig()

	logger.Info("Create User handler")
	userHandler := user.NewHandler(logger)
	userHandler.Register(router)

	run(router, cfg)
}

func run(router *httprouter.Router, cfg *config.Config) {
	logger := logging.GetLogger()

	var (
		listener  net.Listener
		listenErr error
	)

	if cfg.Listen.Type == "sock" {
		appDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			logger.Fatal(err)
		}
		socketPath := path.Join(appDir, "app.sock")

		listener, listenErr = net.Listen("unix", socketPath)
		logger.Infof("Server started to listen unix socket: %s", socketPath)
	} else {
		listener, listenErr = net.Listen(
			"tcp",
			fmt.Sprintf("%s:%s", cfg.Listen.BindIp, cfg.Listen.Port),
		)
		logger.Infof(
			"Server started to listen: %s:%s",
			cfg.Listen.BindIp,
			cfg.Listen.Port,
		)
	}

	if listenErr != nil {
		logger.Fatal(listenErr)
	}

	server := &http.Server{
		Handler:      router,
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
	}

	logger.Fatalln(server.Serve(listener))
}
