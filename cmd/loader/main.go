package main

import (
	"context"
	"errors"
	"log"
	"os"

	"github.com/alexflint/go-arg"

	"github.com/bcap/stock-data/config"
	"github.com/bcap/stock-data/runner"
)

type Args struct {
	Config string `arg:"-c,--config,required" help:"main app config file"`
}

func main() {
	args := Args{}
	arg.MustParse(&args)

	cfg, err := config.Load(args.Config)
	if err != nil {
		log.Printf("couldnt load config: %v", err)
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = runner.New(cfg).Run(ctx)
	if err != nil {
		var multi *runner.ErrMultiple
		if errors.As(err, &multi) {
			log.Printf("run finished with %d errors:", len(multi.Errors))
			for _, err := range multi.Errors {
				log.Printf("  %v", err)
			}
		} else {
			log.Printf("run failed: %v", err)
		}
		os.Exit(1)
	}

	log.Printf("run finished succesfully")
}
