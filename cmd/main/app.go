package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aakosarev/tracking-service/internal/config"
	"github.com/aakosarev/tracking-service/internal/service"
	"github.com/aakosarev/tracking-service/pkg/client/tracking_more"
	"github.com/aakosarev/tracking-service/pkg/logging"
	"log"
	"net/http"
	"os"
)

func main() {
	log.Print("Config initializing")
	cfg := config.GetConfig()
	log.Print("Logger initializing")
	logger := logging.GetLogger(cfg.AppConfig.LogLevel)
	tmClient := tracking_more.NewClient(&http.Client{}, cfg.TrackingMore.BaseUrl, cfg.TrackingMore.ApiKey)
	tmService := service.NewTrackingMoreService(tmClient, logger)

	inputDate := &tracking_more.InputDataForCreatingTracking{
		TrackingNumber: "4058614350",
		CourierCode:    "dhl",
	}

	tracker, err := tmService.CreateTracking(context.Background(), inputDate)

	if err != nil {
		fmt.Fprintln(os.Stderr, "Error creating tracking:", err)
		os.Exit(1)
		return
	}

	prettyJSON, err := json.MarshalIndent(tracker, "", "    ")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error creating JSON:", err)
	}
	fmt.Printf("%s\n", string(prettyJSON))

	// Get result

	trackingNumber := "4058614350"
	result, err := tmService.GetTrackingResult(context.Background(), trackingNumber)

	if err != nil {
		fmt.Fprintln(os.Stderr, "Error get result:", err)
		os.Exit(1)
		return
	}

	prettyJSON, err = json.MarshalIndent(result, "", "    ")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error creating JSON:", err)
	}
	fmt.Printf("%s\n", string(prettyJSON))

}
