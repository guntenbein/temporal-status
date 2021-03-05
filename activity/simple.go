package activity

import (
	"context"
	"time"
)

func LongTermActivity(ctx context.Context) error {
	for i := 0; i < 10; i++ {
		time.Sleep(time.Second)
	}
	return nil
}
