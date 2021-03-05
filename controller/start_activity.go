package controller

import (
	"context"
	"net/http"
	"temporal_starter"
	"temporal_starter/workflow"

	"go.temporal.io/sdk/client"
)

func StartWorkflowHandleFunc(temporalClient client.Client) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		var err error
		defer func() {
			if err != nil {
				writeError(rw, err)
			}
		}()
		ctx := context.Background()
		if err != nil {
			return
		}
		ctxWithRequestMetadata := injectFromHeaders(ctx, req)
		output, err := executeWorkflow(ctxWithRequestMetadata, temporalClient)
		if err != nil {
			return
		}
		err = writeOutput(rw, output)
		if err != nil {
			return
		}
	}
}

func executeWorkflow(ctx context.Context, temporalClient client.Client) (resp interface{}, err error) {
	workflowOptions := client.StartWorkflowOptions{
		TaskQueue: temporal_starter.WorkflowQueue,
	}
	workflowRun, err := temporalClient.ExecuteWorkflow(ctx, workflowOptions, workflow.StatusWorkflow)
	if err != nil {
		return
	}
	var workflowResp interface{}
	err = workflowRun.Get(ctx, &workflowResp)
	if err != nil {
		return
	}
	return workflowResp, nil
}
