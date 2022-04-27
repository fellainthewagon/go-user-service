package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"rest-api/internal/config"
	"rest-api/internal/user"
	"rest-api/internal/user/db"
	"rest-api/pkg/client/postgresql"
	"rest-api/pkg/logging"
	"time"

	"github.com/julienschmidt/httprouter"
)

func main() {
	logger := logging.GetLogger()

	logger.Info("Create new router")
	router := httprouter.New()

	cfg := config.GetConfig()

	postgresqlClient, err := postgresql.NewClient(context.Background(), cfg.PostgreSQLConfig, 5)
	if err != nil {
		logger.Fatal(err)
	}

	// mongodbClient, err := mongodb.NewClient(context.Background(), cfg.MongoDBConfig)
	// if err != nil {
	// 	logger.Fatal(err)
	// }

	// defer func() {
	// 	if err := mongodbClient.Disconnect(context.TODO()); err != nil {
	// 	logger.Fatal(err)
	// 	}
	// }()

	// mongoDatabase := mongodbClient.Database(cfg.MongoDBConfig.Database)
	// userStorage := db.NewMongodbStorage(mongoDatabase, cfg.MongoDBConfig.Collection, logger)

	userStorage2 := db.NewPostgresqlStorage(postgresqlClient, logger)
	userService := user.NewUserService(userStorage2, logger)

	logger.Info("Create User handler")
	userHandler := user.NewHandler(userService, logger)
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
