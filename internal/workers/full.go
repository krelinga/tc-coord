package workers

import (
	"github.com/krelinga/tc-coord/internal/activities"
	"github.com/krelinga/tc-coord/internal/workflows"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func RunFullWorker(c client.Client, interrupt <-chan interface{}) error {
	w := worker.New(c, workflows.TaskQueue, worker.Options{})
	// Register all workflows and activities
	workflows.RegisterDirectory(w)
	activities.RegisterDirKind(w)
	return w.Run(interrupt)
}
