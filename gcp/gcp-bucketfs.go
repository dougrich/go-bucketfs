package gcp

import (
	"context"
	"io"

	"cloud.google.com/go/storage"
	"github.com/dougrich/go-bucketfs"
)

type GCPBucketFile struct {
	obj *storage.ObjectHandle
}

func New(obj *storage.ObjectHandle) *GCPBucketFile {
	return &GCPBucketFile{
		obj,
	}
}

func (b *GCPBucketFile) NewReader(ctx context.Context) (io.Reader, error) {
	reader, err := b.obj.NewReader(ctx)
	if err == storage.ErrObjectNotExist {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return reader, nil
}

func (b *GCPBucketFile) NewWriter(ctx context.Context) (bucketfs.Writer, error) {
	return b.obj.NewWriter(ctx), nil
}
