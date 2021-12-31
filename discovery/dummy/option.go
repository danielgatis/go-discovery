package dummy

import (
	"time"

	"github.com/sirupsen/logrus"
)

// Option is a constructor option.
type Option func(*Resolver)

// WithPeers sets an array of peers.
func WithPeers(peers []string) Option {
	return func(r *Resolver) {
		r.peers = peers
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
