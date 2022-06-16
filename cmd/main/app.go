package main

import (
	"encoding/json"
	"fmt"
	"github.com/aakosarev/tracking-service/internal/config"
	"github.com/aakosarev/tracking-service/pkg/client/tracking_more"
	"github.com/aakosarev/tracking-service/pkg/logging"
	"net/http"
	"os"
)

func main() {
	logging.Init("info")
	cfg := config.GetConfig()
	client := tracking_more.NewClient(&http.Client{}, cfg.TrackingMore.BaseUrl, cfg.TrackingMore.ApiKey)

	tracker, err := client.CreateTracker(
		&tracking_more.InputData{
			TrackingNumber: "1075418304",
			CourierCode:    "dhl",
		},
	)

	if err != nil {
		fmt.Fprintln(os.Stderr, "error creating tracker:", err)
		os.Exit(1)
		return
	}

	prettyJSON, err := json.MarshalIndent(tracker, "", "    ")
	if err != nil {
		fmt.Fprintln(os.Stderr, "error creating JSON:", err)
	}
	fmt.Printf("%s\n", string(prettyJSON))
}
