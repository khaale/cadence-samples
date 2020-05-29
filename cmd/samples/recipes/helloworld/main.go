package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/pborman/uuid"
	"github.com/uber-common/cadence-samples/cmd/samples/common"
	"go.uber.org/cadence/client"
	"go.uber.org/cadence/worker"
)

// This needs to be done as part of a bootstrap step when the process starts.
// The workers are supposed to be long running.
func startWorkers(h *common.SampleHelper) {
	// Configure worker options.
	workerOptions := worker.Options{
		MetricsScope: h.Scope,
		Logger:       h.Logger,
	}
	h.StartWorkers(h.Config.DomainName, ApplicationName, workerOptions)
}

func startWorkflow(h *common.SampleHelper) {
	workflowOptions := client.StartWorkflowOptions{
		ID:                              "helloworld_" + uuid.New(),
		TaskList:                        ApplicationName,
		ExecutionStartToCloseTimeout:    time.Minute * 20,
		DecisionTaskStartToCloseTimeout: time.Minute * 20,
	}
	h.StartWorkflow(workflowOptions, Workflow, h.Config.ParallelActivities)
}

func runWebServer(h *common.SampleHelper) {
	http.HandleFunc("/run", func(w http.ResponseWriter, r *http.Request) {
		startWorkflow(h)
		fmt.Fprintf(w, "Workflow started.")
	})
	http.ListenAndServe(":80", nil)
}

func main() {
	var mode string
	flag.StringVar(&mode, "m", "trigger", "Mode is worker or trigger or launcher.")
	var metricsPort int
	flag.IntVar(&metricsPort, "p", 8080, "Port for exposing metrics")
	flag.Parse()

	var h common.SampleHelper
	h.SetupServiceConfig()
	h.SetupMetrics(metricsPort)

	switch mode {
	case "worker":
		startWorkers(&h)

		// The workers are supposed to be long running process that should not exit.
		// Use select{} to block indefinitely for samples, you can quit by CMD+C.
		select {}
	case "trigger":
		startWorkflow(&h)
	case "launcher":
		runWebServer(&h)
	}
}
