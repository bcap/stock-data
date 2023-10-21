package runner

import (
	"context"
	"fmt"

	"github.com/bcap/stock-data/config"
)

func (r *Runner) fundamentals(ctx context.Context) []waitable {
	s3Bucket := string(r.Config.Fundamentals.S3Bucket)

	work := func(ticker config.Ticker) error {
		s3Path := fmt.Sprintf("%s/fundamentals.%s.json", r.Config.Fundamentals.S3Prefix, ticker)
		fetch := func() ([]byte, error) { return r.eodhdClient.Fundamentals(ctx, ticker) }
		return r.fetchAndUpload(ctx, fetch, s3Bucket, s3Path)
	}

	waitables := []waitable{}
	for _, ticker := range r.Config.Tickers {
		ticker := ticker
		w := r.launch(ctx, func(ctx context.Context) error {
			if err := work(ticker); err != nil {
				return fmt.Errorf("fundamentals(%s) failed: %w", ticker, err)
			}
			return nil
		})
		waitables = append(waitables, w)
	}
	return waitables
}
