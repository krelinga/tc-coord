package main

import (
	"context"
	"testing"

	pb "buf.build/gen/go/krelinga/proto/protocolbuffers/go/krelinga/video/tccoord/v1"
	"connectrpc.com/connect"
	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/testing/protocmp"
)

func TestTcCoord(t *testing.T) {
	// Create a new service instance
	service := &tcCoord{}

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
}
