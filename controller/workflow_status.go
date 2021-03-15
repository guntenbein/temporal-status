package controller

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	temporal_status "temporal_starter"

	"github.com/gorilla/mux"
	temporal_errors "go.temporal.io/api/serviceerror"
	"go.temporal.io/sdk/client"
)

func WorkflowStatusHandleFunc(temporalClient client.Client) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		workflowId, ok := vars["workflowId"]
		if !ok {
			writeError(rw, errors.New("workflowId not found"), http.StatusInternalServerError)
			return
		}
		queryResp, err := temporalClient.QueryWorkflow(context.Background(), workflowId, "", temporal_status.WorkflowQueryTypeStatus)
		if err != nil {
			_, ok := err.(*temporal_errors.NotFound)
			if ok {
				writeError(rw, errors.New("workflow with this workflowId not found"), http.StatusNotFound)
				return
			}
			writeError(rw, err, http.StatusInternalServerError)
			return
		}
		var status temporal_status.Status
		if err := queryResp.Get(&status); err != nil {
			writeError(rw, err, http.StatusInternalServerError)
			return
		}
		writeOk(rw, fmt.Sprintf("%s: %d%%", status.Message, status.Percentage))
	}
}
