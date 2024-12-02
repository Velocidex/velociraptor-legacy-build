package utils

import (
	"context"
	"time"
)

func WithTimeoutCause(ctx context.Context, duration time.Duration, err error) (
	context.Context, func()) {
	return context.WithTimeout(ctx, duration)
}
