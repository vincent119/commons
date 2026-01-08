package graceful

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"time"
)

// Option defines a function to configure the graceful shutdown manager.
type Option func(*options)

type options struct {
	shutdownTimeout time.Duration
	logger          *slog.Logger
	cleaners        []Cleaner
}

// defaultOptions returns the default options.
func defaultOptions() *options {
	return &options{
		shutdownTimeout: 30 * time.Second,
		logger:          slog.Default(),
		cleaners:        make([]Cleaner, 0),
	}
}

// WithTimeout sets the timeout for the shutdown process.
// If cleanup tasks take longer than this duration, they may be cancelled.
// Default is 30 seconds.
func WithTimeout(d time.Duration) Option {
	return func(o *options) {
		o.shutdownTimeout = d
	}
}

// WithLogger sets the logger used by the manager.
// Accepts *slog.Logger.
func WithLogger(l *slog.Logger) Option {
	return func(o *options) {
		if l != nil {
			o.logger = l
		}
	}
}

// WithCleanup adds a cleanup function to be executed during shutdown.
// Cleanup functions are executed in LIFO order.
func WithCleanup(c Cleaner) Option {
	return func(o *options) {
		if c != nil {
			o.cleaners = append(o.cleaners, c)
		}
	}
}

// WithCloser adds an io.Closer to be closed during shutdown.
// The Close method will be called within a Cleaner wrapper.
// Note: Since io.Closer does not accept context, if Close blocks beyond the shutdown timeout,
// the manager will give up waiting and return a timeout error, but the underlying Close
// operation will continue running in the background until it completes or the process exits.
func WithCloser(c io.Closer) Option {
	return func(o *options) {
		if c != nil {
			o.cleaners = append(o.cleaners, func(ctx context.Context) error {
				done := make(chan error, 1)
				go func() {
					done <- c.Close()
				}()

				select {
				case err := <-done:
					return err
				case <-ctx.Done():
					return fmt.Errorf("closer (%T) timed out: %w", c, ctx.Err())
				}
			})
		}
	}
}

// WithClosers adds multiple io.Closer instances to be closed during shutdown.
// This is a convenience version of WithCloser for registering multiple closers at once.
// Example:
//
//	graceful.Run(task,
//	    graceful.WithClosers(db, redisClient, kafkaProducer),
//	)
//
// Equivalent to:
//
//	graceful.Run(task,
//	    graceful.WithCloser(db),
//	    graceful.WithCloser(redisClient),
//	    graceful.WithCloser(kafkaProducer),
//	)
func WithClosers(closers ...io.Closer) Option {
	return func(o *options) {
		for _, c := range closers {
			if c != nil {
				// Copy variable to avoid closure capture issue
				closer := c
				o.cleaners = append(o.cleaners, func(ctx context.Context) error {
					done := make(chan error, 1)
					go func() {
						done <- closer.Close()
					}()

					select {
					case err := <-done:
						return err
					case <-ctx.Done():
						return fmt.Errorf("closer (%T) timed out: %w", closer, ctx.Err())
					}
				})
			}
		}
	}
}
