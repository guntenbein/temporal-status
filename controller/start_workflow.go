package controller

import (
	"context"
	"net/http"
	temporal_status "temporal_starter"
	"temporal_starter/workflow"

	"go.temporal.io/sdk/client"
)

func StartWorkflowHandleFunc(temporalClient client.Client) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		weID, err := executeWorkflow(context.Background(), temporalClient)
		if err != nil {
			writeError(rw, err, http.StatusInternalServerError)
			return
		}
		writeOk(rw, weID)
	}
}

func executeWorkflow(ctx context.Context, temporalClient client.Client) (workflowExecutionID string, err error) {
	workflowOptions := client.StartWorkflowOptions{
		TaskQueue: temporal_status.WorkflowQueue,
	}
	we, err := temporalClient.ExecuteWorkflow(ctx, workflowOptions, workflow.StatusWorkflow)
	if err != nil {
		return "", err
	}
	return we.GetID(), nil
}
