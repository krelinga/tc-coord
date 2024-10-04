package main

import (
	"context"
	"testing"

	pb "buf.build/gen/go/krelinga/proto/protocolbuffers/go/krelinga/video/tccoord/v1"
	"connectrpc.com/connect"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	temporal_testsuite "go.temporal.io/sdk/testsuite"
)

func TestTcCoord(t *testing.T) {
	// Create a new service instance
	devTempOpts := temporal_testsuite.DevServerOptions{}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	devTemp, err := temporal_testsuite.StartDevServer(ctx, devTempOpts)
	require.NoError(t, err)
	service := &tcCoord{
		temporalClient: devTemp.Client(),
	}

	t.Run("EmptyQueue", func(t *testing.T) {
		resp, err := service.GetQueue(context.Background(), &connect.Request[pb.GetQueueRequest]{})
		if err != nil {
			t.Fatalf("GetQueue failed: %v", err)
		}
		expected := &pb.GetQueueResponse{}
		assert.Equal(t, expected, resp.Msg)
	})

	t.Run("EnqueueDirUniqueId", func(t *testing.T) {
		_, err := service.EnqueueDir(context.Background(), &connect.Request[pb.EnqueueDirRequest]{
			Msg: &pb.EnqueueDirRequest{
				Id:  "testid",
				Dir: "testdir",
			},
		})
		if err != nil {
			t.Fatalf("EnqueueDir failed: %v", err)
		}

		resp, err := service.GetQueue(context.Background(), &connect.Request[pb.GetQueueRequest]{})
		if err != nil {
			t.Fatalf("GetQueue failed: %v", err)
		}
		expected := &pb.GetQueueResponse{
			Queue: []*pb.QueueEntry{
				{
					Id:  "testid",
					Dir: "testdir",
				},
			},
		}
		assert.Equal(t, expected, resp.Msg)
	})

	t.Run("EnqueueDirReusedId", func(t *testing.T) {
		_, err := service.EnqueueDir(context.Background(), &connect.Request[pb.EnqueueDirRequest]{
			Msg: &pb.EnqueueDirRequest{
				Id:  "testid",
				Dir: "testdir",
			},
		})
		assert.ErrorIs(t, err, errReusedId)
	})
}
