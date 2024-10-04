package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	pb "buf.build/gen/go/krelinga/proto/protocolbuffers/go/krelinga/video/tccoord/v1"
	"connectrpc.com/connect"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.temporal.io/sdk/client"
	temporal_testsuite "go.temporal.io/sdk/testsuite"
)

func getTemporalBinaryPath(dir string) (string, error) {
	list, err := os.ReadDir(dir)
	if err != nil {
		return "", err
	}
	if len(list) != 1 {
		return "", errors.New("expected exactly one file in directory")
	}
	return filepath.Join(dir, list[0].Name()), nil
}

func initCustomSearchAttribute(ctx context.Context, binaryPath, hostport string) error {
	cmd := exec.CommandContext(ctx, binaryPath, "operator", "search-attribute", "create", "--name=dir", "--type=Keyword")
	cmd.Env = append(os.Environ(), "TEMPORAL_ADDRESS="+hostport)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

func getFreePort() (port int, err error) {
	var a *net.TCPAddr
	if a, err = net.ResolveTCPAddr("tcp", "localhost:0"); err == nil {
		var l *net.TCPListener
		if l, err = net.ListenTCP("tcp", a); err == nil {
			defer l.Close()
			return l.Addr().(*net.TCPAddr).Port, nil
		}
	}
	return
}

func TestTcCoord(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "tc-coord-test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	port, err := getFreePort()
	require.NoError(t, err)
	hostport := fmt.Sprintf("localhost:%d", port)

	// Create a new service instance
	devTempOpts := temporal_testsuite.DevServerOptions{
		CachedDownload: temporal_testsuite.CachedDownload{
			DestDir: tempDir,
		},
		ClientOptions: &client.Options{
			HostPort: hostport,
		},
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	devTemp, err := temporal_testsuite.StartDevServer(ctx, devTempOpts)
	require.NoError(t, err)
	defer devTemp.Stop()

	temporalBinPath, err := getTemporalBinaryPath(tempDir)
	require.NoError(t, err)
	t.Log("temporal binary path:", temporalBinPath)
	err = initCustomSearchAttribute(ctx, temporalBinPath, hostport)
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
