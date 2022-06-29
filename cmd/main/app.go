package main

import (
	"fmt"
	"github.com/aakosarev/tracking-service/internal/config"
	"github.com/aakosarev/tracking-service/internal/parcels_app"
	"github.com/aakosarev/tracking-service/internal/tracking_more"
	"github.com/aakosarev/tracking-service/pkg/dynamo"
	"github.com/aakosarev/tracking-service/pkg/logging"
	"github.com/julienschmidt/httprouter"
	"log"
	"net"
	"net/http"
	"time"
)

func main() {
	log.Println("Config initializing")
	cfg := config.GetConfig()

	log.Println("Logger initializing")
	logger := logging.GetLogger(cfg.AppConfig.LogLevel)

	logger.Println("Router initializing")
	router := httprouter.New()

	dynamoClient, err := dynamo.NewClient()
	if err != nil {
		logger.Fatal(err)
	}

	trackingMoreStorage := tracking_more.NewStorage(dynamoClient, logger)
	trackingMoreService := tracking_more.NewService(trackingMoreStorage, logger)
	trackingMoreHandler := tracking_more.NewHandler(trackingMoreService, logger)
	trackingMoreHandler.Register(router)

	parcelsAppStorage := parcels_app.NewStorage(dynamoClient, logger)
	parcelsAppService := parcels_app.NewService(parcelsAppStorage, logger)
	parcelsAppHandler := parcels_app.NewHandler(parcelsAppService, logger)
	parcelsAppHandler.Register(router)

	logger.Println("Start application")
	start(router, logger, cfg)
}

func start(router http.Handler, logger *logging.Logger, cfg *config.Config) {
	var server *http.Server

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", cfg.Listen.BindIP, cfg.Listen.Port))
	if err != nil {
		logger.Fatal(err)
	}

	server = &http.Server{
		Handler:      router,
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  30 * time.Second,
	}

	if err := server.Serve(listener); err != nil {
		logger.Fatal(err)
	}
}
