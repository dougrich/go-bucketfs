package bucketfs_test

import (
	"testing"

	"github.com/dougrich/go-bucketfs"
)

func TestPlaceholder(t *testing.T) {
	if bucketfs.Placeholder() != 2 {
		t.Fatal("Expected a placeholder value of 2")
	}
}
