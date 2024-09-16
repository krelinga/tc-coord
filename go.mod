module github.com/krelinga/tc-coord

go 1.22.6

require (
	buf.build/gen/go/krelinga/proto/connectrpc/go v1.16.2-20240914043159-8e4e49d91022.1
	buf.build/gen/go/krelinga/proto/protocolbuffers/go v1.34.2-20240914043159-8e4e49d91022.2
	connectrpc.com/connect v1.16.2
	github.com/google/go-cmp v0.5.9
)

require google.golang.org/protobuf v1.34.2 // indirect
