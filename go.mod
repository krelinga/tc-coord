module github.com/krelinga/tc-coord

go 1.22.6

require (
	buf.build/gen/go/krelinga/proto/connectrpc/go v1.16.2-20240914043159-8e4e49d91022.1
	buf.build/gen/go/krelinga/proto/protocolbuffers/go v1.34.2-20240914043159-8e4e49d91022.2
	connectrpc.com/connect v1.16.2
	github.com/google/go-cmp v0.6.0
)

require (
	google.golang.org/grpc v1.66.2
	google.golang.org/protobuf v1.34.2
)

require (
	golang.org/x/sys v0.21.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240604185151-ef581f913117 // indirect
)
