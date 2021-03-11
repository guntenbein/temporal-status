package controller

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	temporal_status "temporal_starter"
	"temporal_starter/workflow"

	"go.temporal.io/sdk/client"
)

func StartWorkflowHandleFunc(temporalClient client.Client) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		weID, err := executeWorkflow(context.Background(), temporalClient)
		if err != nil {
			writeError(rw, err)
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

func writeOk(rw http.ResponseWriter, workflowID string) {
	body, err := json.Marshal(workflowID)
	if err != nil {
		log.Print(err)
	}
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	rw.WriteHeader(http.StatusOK)
	_, err = rw.Write(body)
	if err != nil {
		log.Print(err)
	}
}

func writeError(rw http.ResponseWriter, err error) {
	log.Print(err.Error())
	rw.WriteHeader(http.StatusInternalServerError)
	if _, errWrite := rw.Write([]byte(err.Error())); errWrite != nil {
		log.Print("error writing the HTTP response: " + errWrite.Error())
	}
}
