package gcp_test

import (
	"context"
	"io/ioutil"
	"os"
	"testing"

	"cloud.google.com/go/storage"
	"github.com/dougrich/go-bucketfs"
	gcpbucket "github.com/dougrich/go-bucketfs/gcp"
)

func TestGCPBucketRead(t *testing.T) {

	ctx := context.Background()
	testID := os.Getenv("TEST_ID")
	projectID := os.Getenv("CLOUDSDK_CORE_PROJECT")
	if testID == "" {
		t.Fatalf("TEST_ID environment variable required for tests")
	}
	if projectID == "" {
		t.Fatalf("CLOUDSDK_CORE_PROJECT environment variable required for tests")
	}
	testID = testID + "-reader"

	client, err := storage.NewClient(ctx)
	if err != nil {
		t.Fatalf("Error opening storage, %v", err)
	}

	err = client.Bucket(testID).Create(ctx, projectID, nil)
	if err != nil {
		t.Fatalf("Error creating the bucket, %v", err)
	}
	defer func() {
		if err := client.Bucket(testID).Object("reader").Delete(ctx); err != nil {
			if err != storage.ErrObjectNotExist {
				t.Errorf("Unable to delete storage bucket %s, %v", testID, err)
			}
		}
		if err := client.Bucket(testID).Delete(ctx); err != nil {
			t.Errorf("Unable to delete storage bucket %s, %v", testID, err)
		}
	}()

	var file bucketfs.BucketFile = gcpbucket.New(client.Bucket(testID).Object("reader"))

	// read an empty file out
	reader, err := file.NewReader(ctx)
	if err != nil {
		t.Fatalf("Unexpected err from NewReader, %v", err)
	}
	if reader != nil {
		t.Fatalf("Reader should be nil, file is not initialized, %v", reader)
	}

	// write a test file to the bucket
	writer := client.Bucket(testID).Object("reader").NewWriter(ctx)
	_, err = writer.Write([]byte("Hello world"))
	if err != nil {
		t.Fatalf("Unexpected err writing test data to bucket, %v", err)
	}
	writer.Close()

	// read the test file out
	reader, err = file.NewReader(ctx)
	if err != nil {
		t.Fatalf("Unexpected err from NewReader, %v", err)
	}
	if reader == nil {
		t.Fatalf("Reader should not be nil, file is not initialized, %v", reader)
	}
	bytes, err := ioutil.ReadAll(reader)
	if err != nil {
		t.Fatalf("Unexpected err reading all, %v", err)
	}
	testfile := string(bytes)
	if testfile != "Hello world" {
		t.Fatalf("Unexpected testfile contents, expected 'Hello world', got '%s'", testfile)
	}
}

func TestGCPBucketWrite(t *testing.T) {

	ctx := context.Background()
	testID := os.Getenv("TEST_ID")
	projectID := os.Getenv("CLOUDSDK_CORE_PROJECT")
	if testID == "" {
		t.Fatalf("TEST_ID environment variable required for tests")
	}
	if projectID == "" {
		t.Fatalf("CLOUDSDK_CORE_PROJECT environment variable required for tests")
	}
	testID = testID + "-writer"

	client, err := storage.NewClient(ctx)
	if err != nil {
		t.Fatalf("Error opening storage, %v", err)
	}

	err = client.Bucket(testID).Create(ctx, projectID, nil)
	if err != nil {
		t.Fatalf("Error creating the bucket, %v", err)
	}
	defer func() {
		if err := client.Bucket(testID).Object("writer").Delete(ctx); err != nil {
			if err != storage.ErrObjectNotExist {
				t.Errorf("Unable to delete storage bucket %s, %v", testID, err)
			}
		}
		if err := client.Bucket(testID).Delete(ctx); err != nil {
			t.Errorf("Unable to delete storage bucket %s, %v", testID, err)
		}
	}()

	var file bucketfs.BucketFile = gcpbucket.New(client.Bucket(testID).Object("writer"))

	// read an empty file out
	writer, err := file.NewWriter(ctx)
	if err != nil {
		t.Fatalf("Unexpected err from NewReader, %v", err)
	}
	if writer == nil {
		t.Fatalf("Expected writer not to be null")
	}
	_, err = writer.Write([]byte("Hello world"))
	if err != nil {
		t.Fatalf("Unexpected err writing test data to bucket, %v", err)
	}
	writer.Close()

	// read the test file out
	reader, err := client.Bucket(testID).Object("writer").NewReader(ctx)
	if err != nil {
		t.Fatalf("Unexpected err from NewReader, %v", err)
	}
	if reader == nil {
		t.Fatalf("Reader should not be nil, file is not initialized, %v", reader)
	}
	bytes, err := ioutil.ReadAll(reader)
	if err != nil {
		t.Fatalf("Unexpected err reading all, %v", err)
	}
	testfile := string(bytes)
	if testfile != "Hello world" {
		t.Fatalf("Unexpected testfile contents, expected 'Hello world', got '%s'", testfile)
	}
}
