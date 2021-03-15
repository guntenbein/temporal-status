package signals

import (
	"context"
	"fmt"

	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/client"
)

const PercentageSignalName = "signal-percentage"

type Percentage struct {
	client client.Client
}

func MakePercentageSignaller(client client.Client) Percentage {
	return Percentage{client: client}
}

func (p Percentage) Percentage(ctx context.Context, percentage int32) error {
	if percentage < 0 || percentage > 100 {
		return fmt.Errorf("wrong percentage value: %d", percentage)
	}
	we := activity.GetInfo(ctx).WorkflowExecution
	return p.client.SignalWorkflow(ctx, we.RunID, we.ID, PercentageSignalName, percentage)
}
