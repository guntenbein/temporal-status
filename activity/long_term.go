package activity

import (
	"context"
	"log"
	"temporal_starter/heartbeats"
	"time"
)

type Handler struct {
	progressReporter percentageProgressReporter
}

func MakeHandler(percentageSignaller percentageProgressReporter) Handler {
	return Handler{progressReporter: percentageSignaller}
}

type percentageProgressReporter interface {
	Percentage(ctx context.Context, percentage int32) error
}

func (h Handler) LongTermActivity(ctx context.Context) error {
	heartbeats := heartbeats.StartHeartbeat(ctx, time.Second)
	defer heartbeats.Stop()
	defer h.reportPercentage(ctx, 100)
	iterations := 10
	for i := 0; i < iterations; i++ {
		h.reportPercentage(ctx, int32(i/iterations*100))
		time.Sleep(time.Second)
	}
	return nil
}

func (h Handler) reportPercentage(ctx context.Context, percentage int32) {
	err := h.progressReporter.Percentage(ctx, percentage)
	if err != nil {
		log.Print("error reporting percentage: " + err.Error())
	}
}
