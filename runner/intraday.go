package runner

import (
	"context"
	"fmt"
	"time"

	"github.com/bcap/stock-data/config"
)

func (r *Runner) historicalIntraday(ctx context.Context) []waitable {
	interval := r.Config.HistoricalIntraday.Interval
	s3Bucket := string(r.Config.HistoricalIntraday.S3Bucket)
	s3PathDatelayout := "20060102-150405"

	work := func(ticker config.Ticker, from time.Time, to time.Time) error {
		fromS := from.Format(s3PathDatelayout)
		toS := to.Format(s3PathDatelayout)
		s3Path := fmt.Sprintf(
			"%s/interval-%s/historical-intraday-%s.%s.%s.%s.json",
			r.Config.HistoricalIntraday.S3Prefix, interval, interval, fromS, toS, ticker,
		)
		fetch := func() ([]byte, error) {
			return r.eodhdClient.HistoricalIntraDay(ctx, ticker, interval, time.Time(from), time.Time(to))
		}
		return r.fetchAndUpload(ctx, fetch, s3Bucket, s3Path)
	}

	start := time.Time(r.Config.HistoricalIntraday.TimeRange.From)
	end := time.Time(r.Config.HistoricalIntraday.TimeRange.To)
	at := start

	waitables := []waitable{}
	for at.Before(end) {
		to := end
		if r.Config.HistoricalIntraday.TimeRange.Split > 0 {
			to = at.Add(r.Config.HistoricalIntraday.TimeRange.Split)
		}

		for _, ticker := range r.Config.Tickers {
			ticker := ticker
			at := at
			to := to.Add(-1 * time.Second)
			w := r.launch(ctx, func(ctx context.Context) error {
				if err := work(ticker, at, to); err != nil {
					return fmt.Errorf("historicalIntraday(%s, %s, %s) failed: %w", ticker, at, to, err)
				}
				return nil
			})
			waitables = append(waitables, w)
		}

		at = to
	}
	return waitables
}
