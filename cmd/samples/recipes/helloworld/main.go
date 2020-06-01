package main

import (
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/pborman/uuid"
	"github.com/uber-common/cadence-samples/cmd/samples/common"
	"go.uber.org/cadence/client"
	"go.uber.org/cadence/worker"
)

func getTaskListName(taskListNumber int) string {
	return fmt.Sprintf("%s_%d", ApplicationName, taskListNumber)
}

// This needs to be done as part of a bootstrap step when the process starts.
// The workers are supposed to be long running.
func startWorkers(h *common.SampleHelper) {
	// Configure worker options.
	workerOptions := worker.Options{
		MetricsScope: h.Scope,
		Logger:       h.Logger,
	}
	for i := 0; i < h.Config.TaskListCount; i++ {
		h.StartWorkers(
			h.Config.DomainName,
			getTaskListName(i),
			workerOptions)
	}
}

func startWorkflow(h *common.SampleHelper) {
	workflowOptions := client.StartWorkflowOptions{
		ID:                              "helloworld_" + uuid.New(),
		TaskList:                        getTaskListName(rand.Intn(h.Config.TaskListCount)),
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
	http.ListenAndServe(":8081", nil)
	h.Logger.Info("Launcher listens on localhost:8081/run")
}

func main() {
	var mode string
	flag.StringVar(&mode, "m", "worker", "Mode is worker or trigger or launcher.")
	var metricsPort int
	flag.IntVar(&metricsPort, "p", 8080, "Port for exposing metrics")
	flag.Parse()

	var h common.SampleHelper
	h.SetupServiceConfig()
	h.SetupMetrics(metricsPort)

	switch mode {
	case "worker":
		startWorkers(&h)
		runWebServer(&h)
	case "trigger":
		startWorkflow(&h)
	case "launcher":
		runWebServer(&h)
	}
}
