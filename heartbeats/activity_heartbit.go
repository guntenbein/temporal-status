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
	go func() {
		for {
			select {
			case <-h.cancelc:
				return
			default:
				activity.RecordHeartbeat(ctx)
				time.Sleep(period)
			}
		}
	}()
	return h
}

func (h *SimpleHeartbeat) Stop() {
	close(h.cancelc)
}
