package k8s

import (
	"time"

	"github.com/sirupsen/logrus"
)

// Option is a constructor option.
type Option func(*Discovery)

// WithNamespace sets the k8s namespace.
func WithNamespace(namespace string) Option {
	return func(d *Discovery) {
		d.namespace = namespace
	}
}

// WithLabels sets the labels to lookup.
func WithLabels(labels map[string]string) Option {
	return func(d *Discovery) {
		d.labels = labels
	}
}

// WithInterval sets the lookup interval.
func WithInterval(interval time.Duration) Option {
	return func(d *Discovery) {
		d.interval = interval
	}
}

// WithLogger sets the logger.
func WithLogger(logger logrus.FieldLogger) Option {
	return func(d *Discovery) {
		d.logger = logger
	}
}
