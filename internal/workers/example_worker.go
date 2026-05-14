package workers

import (
	"context"
	"log/slog"
	"time"
)

func StartExampleWorker(ctx context.Context, logger *slog.Logger) {
	ticker := time.NewTicker(time.Minute)

	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				logger.Info("example worker stopped")
				return
			case <-ticker.C:
				logger.Debug("example worker heartbeat")
			}
		}
	}()
}
