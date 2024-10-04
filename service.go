package main

// spell-checker:ignore protocolbuffers tccoord connectrpc

import (
	"context"
	"fmt"
	"log"

	pb "buf.build/gen/go/krelinga/proto/protocolbuffers/go/krelinga/video/tccoord/v1"
	"connectrpc.com/connect"
	"github.com/krelinga/tc-coord/internal/workflows"
	temporal_enums "go.temporal.io/api/enums/v1"
	"go.temporal.io/api/workflowservice/v1"
	"go.temporal.io/sdk/client"
	temporal_converter "go.temporal.io/sdk/converter"
	"go.temporal.io/sdk/temporal"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	errReusedId                 = status.Error(codes.AlreadyExists, "id already in use")
	errQueueTooLarge            = status.Error(codes.Unimplemented, "queue too large")
	errExecutionMissingDir      = status.Error(codes.Internal, "execution missing dir")
	errExecutionCorruptDir      = status.Error(codes.Internal, "execution has corrupt dir")
	errCouldNotCheckForReusedId = status.Error(codes.Internal, "could not check for reused id")
)

type tcCoord struct {
	temporalClient client.Client
}

func (server *tcCoord) EnqueueDir(ctx context.Context, req *connect.Request[pb.EnqueueDirRequest]) (*connect.Response[pb.EnqueueDirResponse], error) {
	countReq := &workflowservice.CountWorkflowExecutionsRequest{
		Query: fmt.Sprintf("WorkflowType = '%s' AND WorkflowId = '%s'", "Directory", req.Msg.Id),
	}
	countResp, err := server.temporalClient.CountWorkflow(ctx, countReq)
	if err != nil {
		log.Print(err)
		return nil, errCouldNotCheckForReusedId
	}
	if countResp.GetCount() > 0 {
		return nil, errReusedId
	}

	opts := client.StartWorkflowOptions{
		ID:                    req.Msg.Id,
		WorkflowIDReusePolicy: temporal_enums.WORKFLOW_ID_REUSE_POLICY_REJECT_DUPLICATE,
		TypedSearchAttributes: temporal.NewSearchAttributes(workflows.DirKey.ValueSet(req.Msg.Dir)),
		TaskQueue:             workflows.TaskQueue,
	}
	input := &workflows.DirectoryInput{
		Dir: req.Msg.Dir,
	}
	_, err = server.temporalClient.ExecuteWorkflow(ctx, opts, workflows.Directory, input)

	return nil, err
}

func (server *tcCoord) GetQueue(ctx context.Context, req *connect.Request[pb.GetQueueRequest]) (*connect.Response[pb.GetQueueResponse], error) {
	tReq := &workflowservice.ListWorkflowExecutionsRequest{
		PageSize: 1000,
		Query:    fmt.Sprintf("WorkflowType = '%s'", "Directory"),
	}
	tResp, err := server.temporalClient.ListWorkflow(ctx, tReq)
	if err != nil {
		return nil, err
	}
	if len(tResp.NextPageToken) > 0 {
		return nil, errQueueTooLarge
	}
	resp := &pb.GetQueueResponse{}
	for _, e := range tResp.Executions {
		dirPayload, ok := e.SearchAttributes.IndexedFields["dir"]
		if !ok {
			return nil, errExecutionMissingDir
		}
		entry := &pb.QueueEntry{
			Id: e.Execution.WorkflowId,
		}
		if err := temporal_converter.GetDefaultDataConverter().FromPayload(dirPayload, &entry.Dir); err != nil {
			return nil, errExecutionCorruptDir
		}
		switch e.Status {
		case temporal_enums.WORKFLOW_EXECUTION_STATUS_COMPLETED:
			entry.Status = pb.QueueEntryStatus_QUEUE_ENTRY_STATUS_DONE
		case temporal_enums.WORKFLOW_EXECUTION_STATUS_RUNNING:
			entry.Status = pb.QueueEntryStatus_QUEUE_ENTRY_STATUS_PROCESSING
		default:
			entry.Status = pb.QueueEntryStatus_QUEUE_ENTRY_STATUS_ERROR
		}
		resp.Queue = append(resp.Queue, entry)
	}
	return &connect.Response[pb.GetQueueResponse]{Msg: resp}, nil
}
