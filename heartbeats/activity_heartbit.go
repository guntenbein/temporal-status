package heartbeats

import (
	"context"
	"time"

	"go.temporal.io/sdk/activity"
)

type SimpleHeartbeat struct {
	cancelc chan struct{}
}

func StartHeartbeat(ctx context.Context, period time.Duration) (h *SimpleHeartbeat) {
	h = &SimpleHeartbeat{make(chan struct{})}
	timer := time.NewTimer(period)
	go func() {
		for {
			select {
			case <-h.cancelc:
				timer.Stop()
				return
			case <-timer.C:
				activity.RecordHeartbeat(ctx)
			}
		}
	}()
	return h
}

func (h *SimpleHeartbeat) Stop() {
	close(h.cancelc)
}
