package main

import (
	"encoding/json"
	"fmt"
	"github.com/aakosarev/tracking-service/internal/config"
	"github.com/aakosarev/tracking-service/pkg/client/tracking_more"
	"net/http"
	"os"
)

func main() {
	cfg := config.GetConfig()
	client := tracking_more.NewClient(&http.Client{}, cfg.TrackingMore.BaseUrl, cfg.TrackingMore.ApiKey)

	// Create tracking
	/*
		tracker, err := client.CreateTracking(
			&tracking_more.InputDataForCreatingTracking{
				TrackingNumber: "1Z6R57A00492491127",
				CourierCode:    "ups",
			},
		)

		if err != nil {
			fmt.Fprintln(os.Stderr, "Error creating tracking:", err)
			os.Exit(1)
			return
		}

		prettyJSON, err := json.MarshalIndent(tracker, "", "    ")
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error creating JSON:", err)
		}
		fmt.Printf("%s\n", string(prettyJSON)) */

	// Get result

	result, err := client.GetResult("1Z6R57A00492491127")

	if err != nil {
		fmt.Fprintln(os.Stderr, "Error get result:", err)
		os.Exit(1)
		return
	}

	prettyJSON, err := json.MarshalIndent(result, "", "    ")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error creating JSON:", err)
	}

	fmt.Printf("%s\n", string(prettyJSON))
}
