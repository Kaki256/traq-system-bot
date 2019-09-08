package traq_system_bot

import (
	"cloud.google.com/go/logging"
	"context"
	mrpb "google.golang.org/genproto/googleapis/api/monitoredres"
	"net/http"
	"os"
)

var (
	logger       *logging.Logger
	projectID    string
	functionName string
	region       string
)

func loggerInit() {
	projectID = os.Getenv("GCP_PROJECT")
	functionName = os.Getenv("FUNCTION_NAME")
	region = os.Getenv("FUNCTION_REGION")
	client, err := logging.NewClient(context.Background(), projectID)
	if err != nil {
		panic(err)
	}

	logger = client.Logger("cloudfunctions.googleapis.com/cloud-functions", logging.CommonResource(&mrpb.MonitoredResource{
		Type: "cloud_function",
		Labels: map[string]string{
			"function_name": functionName,
			"project_id":    projectID,
			"region":        region,
		},
	}))

}

func infoL(r *http.Request, payload interface{}) {
	log(r, payload, logging.Info)
}

func errorL(r *http.Request, payload interface{}) {
	log(r, payload, logging.Error)
}

func log(r *http.Request, payload interface{}, severity logging.Severity) {
	logger.Log(logging.Entry{
		Severity: severity,
		Labels:   map[string]string{"execution_id": r.Header.Get("Function-Execution-Id")},
		Payload:  payload,
		Trace:    r.Header.Get("X-Cloud-Trace-Context"),
	})
}
