package main

import (
	"fmt"
	"github.com/aakosarev/tracking-service/internal/config"
	"github.com/aakosarev/tracking-service/internal/tracking"
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

	trackingService := tracking.Service{}

	trackingHandler := tracking.Handler{
		Logger:  logger,
		Service: trackingService,
	}
	trackingHandler.Register(router)

	logger.Println("Start application")
	start(router)
}

func start(router http.Handler) {
	var server *http.Server

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", "127.0.0.1", "1234"))
	if err != nil {
		fmt.Println(err)
	}

	server = &http.Server{
		Handler:      router,
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  30 * time.Second,
	}

	if err := server.Serve(listener); err != nil {
		fmt.Println(err)
	}
}
