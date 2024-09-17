package main

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"

	pb "buf.build/gen/go/krelinga/proto/protocolbuffers/go/krelinga/video/tccoord/v1"
	"connectrpc.com/connect"
	"github.com/google/go-cmp/cmp"
	"github.com/krelinga/tc-coord/internal/mocks"
	"google.golang.org/protobuf/testing/protocmp"
)

func TestTcCoord(t *testing.T) {
	mockBe := mocks.NewTCServiceClient(t)
	mockBe.EXPECT().StartAsyncTranscode(mock.Anything, mock.Anything).Return(nil, nil).Maybe()

	// Create a new service instance
	service := newTcCoord(mockBe)

	t.Run("EmptyQueue", func(t *testing.T) {
		resp, err := service.GetQueue(context.Background(), &connect.Request[pb.GetQueueRequest]{})
		if err != nil {
			t.Fatalf("GetQueue failed: %v", err)
		}
		expected := &pb.GetQueueResponse{
			Queue: []*pb.QueueEntry{},
		}
		if !cmp.Equal(resp.Msg, expected, protocmp.Transform()) {
			t.Errorf("GetQueue returned unexpected response: %v", cmp.Diff(resp, expected, protocmp.Transform()))
		}
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
		if !cmp.Equal(resp.Msg, expected, protocmp.Transform()) {
			t.Errorf("GetQueue returned unexpected response: %v", cmp.Diff(resp, expected, protocmp.Transform()))
		}
	})

	t.Run("EnqueueDirReusedId", func(t *testing.T) {
		_, err := service.EnqueueDir(context.Background(), &connect.Request[pb.EnqueueDirRequest]{
			Msg: &pb.EnqueueDirRequest{
				Id:  "testid",
				Dir: "testdir",
			},
		})
		if !errors.Is(err, errReusedId) {
			t.Errorf("EnqueueDir did not return expected error: %v", err)
		}
	})
}
