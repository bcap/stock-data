package s3

import (
	"bytes"
	"context"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3 struct {
	client *s3.Client
}

func New(config aws.Config) *S3 {
	return &S3{
		client: s3.NewFromConfig(config),
	}
}

func (s *S3) Put(ctx context.Context, bucket string, path string, data []byte) (string, error) {
	buf := bytes.NewBuffer(data)
	out, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:        &bucket,
		Key:           &path,
		Body:          buf,
		ContentLength: int64(len(data)),
	})
	if err != nil {
		return "", err
	}

	// aws s3 api is weird/buggy: etags can be returned as a quoted string. Remove them if thats the case
	etag := *out.ETag
	if etag[0] == '"' {
		etag = etag[1:]
	}
	if etag[len(etag)-1] == '"' {
		etag = etag[:len(etag)-1]
	}

	return etag, nil
}

func (s *S3) Get(ctx context.Context, bucket string, path string) ([]byte, error) {
	out, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &path,
	})
	if err != nil {
		return nil, err
	}
	defer out.Body.Close()

	data, err := io.ReadAll(out.Body)
	if err != nil {
		return nil, err
	}
	return data, nil
}
