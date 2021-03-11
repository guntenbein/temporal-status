package workflow

import (
	temporal_status "temporal_starter"
	"temporal_starter/activity"
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

func StatusWorkflow(ctx workflow.Context) (err error) {
	status := "STARTED"
	err = workflow.SetQueryHandler(ctx,
		temporal_status.WorkflowQueryTypeStatus, func(input []byte) (temporal_status.Status, error) {
			return temporal_status.Status{Message: status}, nil
		})
	if err != nil {
		return
	}
	defer func() {
		status = "FINISHED"
	}()

	ctx = withActivityOptions(ctx, temporal_status.WorkflowQueue)
	status = "PROCESSING ACTIVITY 1"
	err = workflow.ExecuteActivity(ctx, activity.LongTermActivity).Get(ctx, nil)
	if err != nil {
		return
	}

	status = "PROCESSING ACTIVITY 2"
	err = workflow.ExecuteActivity(ctx, activity.LongTermActivity).Get(ctx, nil)
	if err != nil {
		return
	}

	status = "PROCESSING ACTIVITY 3"
	return workflow.ExecuteActivity(ctx, activity.LongTermActivity).Get(ctx, nil)
}

func withActivityOptions(ctx workflow.Context, queue string) workflow.Context {
	ao := workflow.ActivityOptions{
		TaskQueue:              queue,
		ScheduleToStartTimeout: 24 * time.Hour,
		StartToCloseTimeout:    24 * time.Hour,
		HeartbeatTimeout:       time.Second * 5,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:        time.Second,
			BackoffCoefficient:     2.0,
			MaximumInterval:        time.Minute * 5,
			NonRetryableErrorTypes: []string{"BusinessError"},
		},
	}
	ctxOut := workflow.WithActivityOptions(ctx, ao)
	return ctxOut
}
