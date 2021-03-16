package workflow

import (
	"log"
	"sync/atomic"
	temporal_status "temporal_starter"
	"temporal_starter/activity"
	"temporal_starter/signals"
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

const finish = 1

func StatusWorkflow(ctx workflow.Context) (err error) {
	status := "STARTED"
	var percentage int32
	err = workflow.SetQueryHandler(ctx,
		temporal_status.WorkflowQueryTypeStatus, func(input []byte) (temporal_status.Status, error) {
			return temporal_status.Status{Message: status, Percentage: atomic.LoadInt32(&percentage)}, nil
		})
	if err != nil {
		return
	}
	defer func() {
		status = "FINISHED"
	}()

	var finishFlag int32
	workflow.Go(ctx, func(ctx workflow.Context) {
		err := workflow.Await(ctx, func() bool {
			var pcn int32
			ok := workflow.GetSignalChannel(ctx, signals.PercentageSignalName).ReceiveAsync(&pcn)
			if ok {
				atomic.StoreInt32(&percentage, pcn)
			}
			return atomic.LoadInt32(&finishFlag) == finish
		})
		if err != nil {
			log.Print(err.Error())
		}
	})
	defer func() {
		atomic.StoreInt32(&finishFlag, 1)
	}()

	return statusWorkflow(ctx, &status)
}

func statusWorkflow(ctx workflow.Context, status *string) (err error) {
	ctx = withActivityOptions(ctx, temporal_status.WorkflowQueue)
	*status = "PROCESSING ACTIVITY 1"
	err = workflow.ExecuteActivity(ctx, activity.Handler{}.LongTermActivity).Get(ctx, nil)
	if err != nil {
		return
	}

	*status = "PROCESSING ACTIVITY 2"
	err = workflow.ExecuteActivity(ctx, activity.Handler{}.LongTermActivity).Get(ctx, nil)
	if err != nil {
		return
	}

	*status = "PROCESSING ACTIVITY 3"
	return workflow.ExecuteActivity(ctx, activity.Handler{}.LongTermActivity).Get(ctx, nil)
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
