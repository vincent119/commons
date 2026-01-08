// Package graceful 提供優雅關機（Graceful Shutdown）的管理器與工具。
package graceful

import (
	"context"
	"errors"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

// Task represents a long-running task that should listen for context cancellation.
// It should return when the context is done or when a fatal error occurs.
type Task func(ctx context.Context) error

// Cleaner defines a function that releases resources during shutdown.
// Important: Cleaner must respect the ctx timeout. If ctx.Done() is received,
// it should return ctx.Err() immediately, otherwise it will block subsequent cleanup.
type Cleaner func(ctx context.Context) error

// Run executes the given task and handles graceful shutdown on system signals.
// It listens for SIGINT and SIGTERM.
func Run(task Task, opts ...Option) error {
	o := defaultOptions()
	for _, opt := range opts {
		opt(o)
	}

	// 1. Setup signal context
	// NotifyContext returns a copy of the parent context that is marked done
	// (its Done channel is closed) when one of the listed signals arrives,
	// when the returned stop function is called, or when the parent context's
	// Done channel is closed, whichever happens first.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// 2. Run the task
	o.logger.Info("starting task")
	startTime := time.Now()

	// Execute the task. Expected to block until done or ctx is done.
	err := task(ctx)

	// Log task exit
	duration := time.Since(startTime)
	if err != nil {
		o.logger.Error("task exited with error", "error", err, "duration", duration)
	} else {
		o.logger.Info("task exited successfully", "duration", duration)
	}

	// 3. Run cleanup
	// We create a new context for cleanup since the signal context is already done.
	o.logger.Info("starting shutdown cleanup", "timeout", o.shutdownTimeout)

	shutdownCtx, cancel := context.WithTimeout(context.Background(), o.shutdownTimeout)
	defer cancel()

	var cleanupErrors []error
	// Execute cleaners in LIFO order (Last-In-First-Out)
	// Ensures resources with higher dependencies (usually registered later) are released first
	for i := len(o.cleaners) - 1; i >= 0; i-- {
		c := o.cleaners[i]
		if cErr := c(shutdownCtx); cErr != nil {
			o.logger.Error("cleanup failed", "error", cErr)
			cleanupErrors = append(cleanupErrors, cErr)
		}
	}

	if len(cleanupErrors) > 0 {
		// Include task error if present, to avoid swallowing errors
		if err != nil {
			cleanupErrors = append([]error{err}, cleanupErrors...)
		}
		return errors.Join(cleanupErrors...)
	}

	o.logger.Info("shutdown complete")
	return err
}

// HTTPTask wraps an http.Server as a graceful.Task.
// It starts the server in a goroutine and waits for context cancellation.
// Note: You typically need to register a cleanup function to shutdown the server, e.g.:
// graceful.WithCleanup(func(ctx context.Context) error { return srv.Shutdown(ctx) })
// or use a helper function that handles this logic if preferred.
//
// However, to keep it simple and composable, this helper only handles the "run" part.
// The user is responsible for ensuring srv.Shutdown is called during cleanup.
func HTTPTask(srv *http.Server) Task {
	return func(ctx context.Context) error {
		errCh := make(chan error, 1)
		go func() {
			if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				errCh <- err
			}
			close(errCh)
		}()

		select {
		case err := <-errCh:
			// Server failed or stopped
			return err
		case <-ctx.Done():
			// Signal received. Return to allow cleanup phase to begin.
			// The server is still running at this point!
			// The cleanup phase (which should include srv.Shutdown) will stop it.
			return nil
		}
	}
}
