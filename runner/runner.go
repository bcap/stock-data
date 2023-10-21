package runner

import (
	"context"
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
	if r.shouldRun(config.LoadExchanges) {
		waitables = append(waitables, r.exchanges(ctx)...)
	}
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

func (r *Runner) fetchAndUpload(ctx context.Context, fetch func() ([]byte, error), s3Bucket string, s3Path string) error {
	data, err := fetch()
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
