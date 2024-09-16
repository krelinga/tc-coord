package main

// spell-checker:ignore protocolbuffers tccoord connectrpc

import (
	"context"

	pb "buf.build/gen/go/krelinga/proto/protocolbuffers/go/krelinga/video/tccoord/v1"
	"connectrpc.com/connect"
)

type tcCoord struct {
}

func (server *tcCoord) EnqueueDir(ctx context.Context, req *connect.Request[pb.EnqueueDirRequest]) (*connect.Response[pb.EnqueueDirResponse], error) {
	return nil, nil
}

func (server *tcCoord) GetQueue(ctx context.Context, req *connect.Request[pb.GetQueueRequest]) (*connect.Response[pb.GetQueueResponse], error) {
	return &connect.Response[pb.GetQueueResponse]{
		Msg: &pb.GetQueueResponse{},
	}, nil
}
