package runner

import (
	"context"
	"fmt"
	"log"

	"github.com/bcap/stock-data/config"
)

func (r *Runner) fundamentals(ctx context.Context) []waitable {
	work := func(ticker config.Ticker) error {
		data, err := r.eodhdClient.Fundamentals(ctx, ticker)
		if err != nil {
			return err
		}

		s3Bucket := string(r.Config.Fundamentals.S3Bucket)
		s3Path := fmt.Sprintf("%s/fundamentals.%s.json", r.Config.Fundamentals.S3Prefix, ticker)

		etag, err := r.s3Client.Put(ctx, s3Bucket, s3Path, data)
		if err != nil {
			return err
		}

		log.Printf("s3://%s/%s: %s", s3Bucket, s3Path, etag)
		return nil
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
