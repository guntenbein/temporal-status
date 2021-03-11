package controller

import (
	"context"
	"log"
	"net/http"
	temporal_status "temporal_starter"
	"temporal_starter/workflow"

	"go.temporal.io/sdk/client"
)

func StartWorkflowHandleFunc(temporalClient client.Client) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		err := executeWorkflow(context.Background(), temporalClient)
		if err != nil {
			writeError(rw, err)
			return
		}
		writeOk(rw)
	}
}

func executeWorkflow(ctx context.Context, temporalClient client.Client) error {
	workflowOptions := client.StartWorkflowOptions{
		TaskQueue: temporal_status.WorkflowQueue,
	}
	_, err := temporalClient.ExecuteWorkflow(ctx, workflowOptions, workflow.StatusWorkflow)
	if err != nil {
		return err
	}
	return nil
}

func writeOk(rw http.ResponseWriter) {
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	rw.WriteHeader(http.StatusOK)
}

func writeError(rw http.ResponseWriter, err error) {
	log.Print(err.Error())
	rw.WriteHeader(http.StatusInternalServerError)
	if _, errWrite := rw.Write([]byte(err.Error())); errWrite != nil {
		log.Print("error writing the HTTP response: " + errWrite.Error())
	}
}
