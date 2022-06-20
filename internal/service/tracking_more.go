package service

import (
	"context"
	"encoding/json"
	tm "github.com/aakosarev/tracking-service/pkg/client/tracking_more"
	"github.com/aakosarev/tracking-service/pkg/logging"
)

type TrackingMoreService struct {
	Client *tm.Client
	Logger *logging.Logger
}

func NewTrackingMoreService(client *tm.Client, logger *logging.Logger) *TrackingMoreService {
	return &TrackingMoreService{Client: client, Logger: logger}
}

func (tms *TrackingMoreService) CreateTracking(ctx context.Context, in *tm.InputDataForCreatingTracking) (out *tm.CreatingTrackingResult, err error) {
	resp, err := tms.Client.CreateTracking(ctx, in)
	if err != nil {
		return nil, err
	}
	if err = json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	return out, nil
}

func (tms *TrackingMoreService) GetTrackingResult(ctx context.Context, trackingNumber string) (out *tm.TrackingResult, err error) {
	resp, err := tms.Client.GetResult(ctx, trackingNumber)
	if err != nil {
		return nil, err
	}
	if err = json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	return out, nil
}
