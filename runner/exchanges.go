package runner

import (
	"context"
	"fmt"
	"log"
)

func (r *Runner) exchanges(ctx context.Context) []waitable {
	work := func() error {
		s3Bucket := string(r.Config.Exchanges.S3Bucket)
		s3Path := fmt.Sprintf("%s/exchanges.json", r.Config.Exchanges.S3Prefix)

		data, err := r.eodhdClient.Exchanges(ctx)
		if err != nil {
			return err
		}

		etag, err := r.s3Client.Put(ctx, s3Bucket, s3Path, data)
		if err != nil {
			return err
		}

		log.Printf("s3://%s/%s: %s", s3Bucket, s3Path, etag)
		return nil
	}

	w := r.launch(ctx, func(ctx context.Context) error {
		if err := work(); err != nil {
			return fmt.Errorf("exchanges() failed: %w", err)
		}
		return nil
	})
	return []waitable{w}
}
