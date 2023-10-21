package runner

import (
	"context"
	"fmt"
)

type ErrMultiple struct {
	Errors []error
}

func (e *ErrMultiple) Error() string {
	return fmt.Sprintf("multiple errors: %s", e.Errors)
}

func (r *Runner) collect(ctx context.Context, waitables []waitable) []error {
	errors := []error{}
	for _, w := range waitables {
		err := w(ctx)
		if err != nil {
			errors = append(errors, err)
		}
	}
	return errors
}
