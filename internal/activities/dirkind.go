package activities

import (
	"context"

	"github.com/krelinga/go-lib/video"
)

type DirKindInput struct {
	Dir string `json:"dir"`
}

type DirKindOutput struct {
	Kind video.DirKind `json:"kind"`
}

func DirKind(ctx context.Context, input *DirKindInput) (*DirKindOutput, error) {
	kind, err := video.GetDirKind(input.Dir)
	if err != nil {
		return nil, err
	}
	return &DirKindOutput{
		Kind: kind,
	}, nil
}
