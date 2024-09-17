package main

// spell-checker:ignore protocolbuffers tccoord connectrpc

import (
	"context"

	be_rpc "buf.build/gen/go/krelinga/proto/connectrpc/go/krelinga/video/tcserver/v1/tcserverv1connect"
	pb "buf.build/gen/go/krelinga/proto/protocolbuffers/go/krelinga/video/tccoord/v1"
	"connectrpc.com/connect"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	errReusedId = status.Error(codes.AlreadyExists, "id already in use")
)

type tcCoord struct {
	queue   map[string]*pb.QueueEntry
	backend be_rpc.TCServiceClient
}

func newTcCoord(backend be_rpc.TCServiceClient) *tcCoord {
	return &tcCoord{
		queue:   make(map[string]*pb.QueueEntry),
		backend: backend,
	}
}

func (server *tcCoord) EnqueueDir(ctx context.Context, req *connect.Request[pb.EnqueueDirRequest]) (*connect.Response[pb.EnqueueDirResponse], error) {
	if _, alreadyExists := server.queue[req.Msg.Id]; alreadyExists {
		return nil, errReusedId
	}
	server.queue[req.Msg.Id] = &pb.QueueEntry{
		Id:  req.Msg.Id,
		Dir: req.Msg.Dir,
	}
	return &connect.Response[pb.EnqueueDirResponse]{}, nil
}

func (server *tcCoord) GetQueue(ctx context.Context, req *connect.Request[pb.GetQueueRequest]) (*connect.Response[pb.GetQueueResponse], error) {
	var queue []*pb.QueueEntry
	for _, entry := range server.queue {
		queue = append(queue, entry)
	}
	return &connect.Response[pb.GetQueueResponse]{
		Msg: &pb.GetQueueResponse{
			Queue: queue,
		},
	}, nil
}
