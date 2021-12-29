package k8s

import (
	"time"

	"github.com/sirupsen/logrus"
)

// Option is a constructor option.
type Option func(*Resolver)

// WithNamespace sets the k8s namespace.
func WithNamespace(namespace string) Option {
	return func(r *Resolver) {
		r.namespace = namespace
	}
}

// WithLabels sets the labels to lookup.
func WithLabels(labels map[string]string) Option {
	return func(r *Resolver) {
		r.labels = labels
	}
}

// WithInterval sets the lookup interval.
func WithInterval(interval time.Duration) Option {
	return func(r *Resolver) {
		r.interval = interval
	}
}

// WithLogger sets the logger.
func WithLogger(logger logrus.FieldLogger) Option {
	return func(r *Resolver) {
		r.logger = logger
	}
}
