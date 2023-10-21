package runner

import (
	"context"
	"fmt"
	"log"
	"time"

	awsConfig "github.com/aws/aws-sdk-go-v2/config"

	"github.com/bcap/stock-data/aws/s3"
	"github.com/bcap/stock-data/config"
	"github.com/bcap/stock-data/eodhd"
	"github.com/bcap/stock-data/executor"
)

type void = struct{}

type waitable = func(context.Context) error

type Runner struct {
	Config      config.Config
	eodhdClient *eodhd.Client
	s3Client    *s3.S3
	executor    *executor.Executor[void]
}

func New(cfg config.Config) *Runner {
	awsCfg, _ := awsConfig.LoadDefaultConfig(
		context.Background(),
		awsConfig.WithRegion(cfg.AWS.Region),
		awsConfig.WithSharedConfigProfile(cfg.AWS.Profile),
	)
	return &Runner{
		Config:      cfg,
		eodhdClient: eodhd.NewClient(cfg.EODHD.ApiKey),
		s3Client:    s3.New(awsCfg),
		executor:    executor.New[struct{}](cfg.MaxParallel),
	}
}

func (r *Runner) Run(ctx context.Context) error {
	stopped := make(chan struct{})
	defer close(stopped)

	go func() {
		for {
			select {
			case <-time.After(5 * time.Second):
				log.Printf("Runner status: %d running tasks, %d pending tasks", r.executor.Running(), r.executor.Pending())
			case <-stopped:
				return
			}
		}
	}()

	waitables := []waitable{}
	if r.shouldRun(config.LoadFundamentals) {
		waitables = append(waitables, r.fundamentals(ctx)...)
	}
	if r.shouldRun(config.LoadHistoricalIntraday) {
		waitables = append(waitables, r.historicalIntraday(ctx)...)
	}

	errors := r.collect(ctx, waitables)
	if len(errors) > 0 {
		return &ErrMultiple{Errors: errors}
	}
	return nil
}

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

func (r *Runner) historicalIntraday(ctx context.Context) []waitable {
	interval := r.Config.HistoricalIntraday.Interval
	s3Bucket := string(r.Config.HistoricalIntraday.S3Bucket)
	s3PathDatelayout := "20060102-150405"

	work := func(ticker config.Ticker, from time.Time, to time.Time) error {
		fromS := from.Format(s3PathDatelayout)
		toS := to.Format(s3PathDatelayout)

		s3Path := fmt.Sprintf(
			"%s/interval-%s/historical-intraday-%s.%s.%s.%s.json",
			r.Config.HistoricalIntraday.S3Prefix, interval, interval, ticker, fromS, toS,
		)

		data, err := r.eodhdClient.HistoricalIntraDay(ctx, ticker, interval, time.Time(from), time.Time(to))
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

func (r *Runner) shouldRun(lo config.LoadOption) bool {
	_, ok := r.Config.Load[lo]
	return ok
}

func (r *Runner) launch(ctx context.Context, fn func(context.Context) error) waitable {
	resultFn := r.executor.Launch(ctx, func(ctx context.Context) (void, error) {
		return void{}, fn(ctx)
	})
	return func(ctx context.Context) error {
		_, err := resultFn(ctx)
		return err
	}
}
