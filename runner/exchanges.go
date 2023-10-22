package runner

import (
	"context"
	"fmt"
)

func (r *Runner) listExchanges(ctx context.Context) []waitable {
	work := func() error {
		s3Bucket := string(r.Config.ListExchanges.S3Bucket)
		s3Path := fmt.Sprintf("%s/exchanges.json", r.Config.ListExchanges.S3Prefix)
		fetch := func() ([]byte, error) { return r.eodhdClient.Exchanges(ctx) }
		return r.fetchAndUpload(ctx, fetch, s3Bucket, s3Path)
	}
	w := r.launch(ctx, func(ctx context.Context) error {
		if err := work(); err != nil {
			return fmt.Errorf("listExchanges() failed: %w", err)
		}
		return nil
	})
	return []waitable{w}
}
