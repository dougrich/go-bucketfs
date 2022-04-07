package bucketfs

import (
	"context"
	"io"
)

type Writer interface {
	io.Writer
	Close() error
}

type BucketFile interface {
	NewReader(ctx context.Context) (io.Reader, error)
	NewWriter(ctx context.Context) (Writer, error)
}
