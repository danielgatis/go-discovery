package mdns

import (
	"time"

	"github.com/sirupsen/logrus"
)

// Option is a constructor option.
type Option func(*Discovery)

// WithInstance sets an instance name.
func WithInstance(instance string) Option {
	return func(d *Discovery) {
		d.instance = instance
	}
}

// WithService sets a service name.
func WithService(service string) Option {
	return func(d *Discovery) {
		d.service = service
	}
}

// WithDomain sets the domain.
func WithDomain(domain string) Option {
	return func(d *Discovery) {
		d.domain = domain
	}
}

// WithPort sets the port.
func WithPort(port int) Option {
	return func(d *Discovery) {
		d.port = port
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
