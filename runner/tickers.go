package runner

import (
	"context"
	"fmt"
)

func (r *Runner) listTickers(ctx context.Context) []waitable {
	work := func(exchange string) error {
		s3Bucket := string(r.Config.ListTickers.S3Bucket)
		s3Path := fmt.Sprintf("%s/tickers.%s.json", r.Config.ListTickers.S3Prefix, exchange)
		fetch := func() ([]byte, error) { return r.eodhdClient.Tickers(ctx, exchange, r.start) }
		return r.fetchAndUpload(ctx, fetch, s3Bucket, s3Path)
	}
	waitables := []waitable{}
	for _, exchange := range r.Config.Exchanges {
		exchange := exchange
		w := r.launch(ctx, func(ctx context.Context) error {
			if err := work(exchange); err != nil {
				return fmt.Errorf("listTickers(%s) failed: %w", exchange, err)
			}
			return nil
		})
		waitables = append(waitables, w)
	}
	return waitables
}
