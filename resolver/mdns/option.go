package mdns

import (
	"time"

	"github.com/sirupsen/logrus"
)

// Option is a constructor option.
type Option func(*Resolver)

// WithInstance sets an instance name.
func WithInstance(instance string) Option {
	return func(r *Resolver) {
		r.instance = instance
	}
}

// WithService sets a service name.
func WithService(service string) Option {
	return func(r *Resolver) {
		r.service = service
	}
}

// WithDomain sets the domain.
func WithDomain(domain string) Option {
	return func(r *Resolver) {
		r.domain = domain
	}
}

// WithPort sets the port.
func WithPort(port int) Option {
	return func(r *Resolver) {
		r.port = port
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
