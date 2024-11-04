package util

import (
	"context"
	"log/slog"
	"sync"
	"time"
)

var wg sync.WaitGroup

func SpawnBackgroundTask(ctx context.Context, taskName string, bgTask func(context.Context)) {
	_, span := Tracer.Start(ctx, taskName)
	wg.Add(1)
	go func(ctx context.Context) {
		defer wg.Done()
		defer span.End()
		defer slog.InfoContext(ctx, "background task completed", "taskName", taskName) // TODO is defer even useful here?
		bgTask(ctx)
	}(ctx)
}

func WaitForBackgroundTasks(timeout time.Duration) {
	slog.Info("waiting for background tasks to complete", "timeout", timeout.String())
	c := make(chan struct{})
	go func() {
		defer close(c)
		wg.Wait()
	}()
	select {
	case <-c:
		slog.Info("all background tasks completed")
	case <-time.After(timeout):
		slog.Warn("some background tasks timed out")
	}
}
