package activity

import (
	"context"
	"temporal_starter/heartbeats"
	"time"
)

func LongTermActivity(ctx context.Context) error {
	heartbeats := heartbeats.StartHeartbeat(ctx, time.Second)
	defer heartbeats.Stop()

	for i := 0; i < 10; i++ {
		time.Sleep(time.Second)
	}
	return nil
}
