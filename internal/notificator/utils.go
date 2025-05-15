package notificator

import (
	"context"
	"time"
)

func waitForTickerOrCtx(ctx context.Context, ticker *time.Ticker) {
	select {
	case <-ticker.C:
	case <-ctx.Done():
	}
}

func isCancelled(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		return true
	default:
		return false
	}
}
