package workflows

import (
	"time"

	"github.com/krelinga/tc-coord/internal/activities"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

type DirectoryInput struct {
	Dir string `json:"dir"`
}

func Directory(ctx workflow.Context, input *DirectoryInput) error {
	opts := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, opts)
	dirKindInput := &activities.DirKindInput{
		Dir: input.Dir,
	}
	var dirKindOutput activities.DirKindOutput
	err := workflow.ExecuteActivity(ctx, activities.DirKind, dirKindInput).Get(ctx, &dirKindOutput)
	if err != nil {
		return err
	}
	return nil
}

func RegisterDirectory(w worker.Worker) {
	registerOpts := workflow.RegisterOptions{Name: "Directory"}
	w.RegisterWorkflowWithOptions(Directory, registerOpts)
}
