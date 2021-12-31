package dummy

import (
	"time"

	"github.com/sirupsen/logrus"
)

// Resolver is a dummy resolver for static peers.
type Resolver struct {
	peers    []string
	logger   logrus.FieldLogger
	interval time.Duration
	output   chan string
	stop     chan struct{}
}

// New returns a new dummy resolver.
func New(peers []string, opts ...Option) *Resolver {
	const (
		defaultInterval = 2 * time.Second
	)

	var (
		defaultLogger = logrus.StandardLogger()
	)

	r := &Resolver{
		peers:    peers,
		logger:   defaultLogger,
		interval: defaultInterval,
		output:   make(chan string),
		stop:     make(chan struct{}),
	}

	for _, opt := range opts {
		opt(r)
	}

	return r
}

// Start implements resolver.Resolver.
func (d *Resolver) Start() (chan string, error) {
	ticker := time.NewTicker(d.interval)

	go func() {
		for {
			select {
			case <-d.stop:
				ticker.Stop()
				return
			case <-ticker.C:
				for _, peer := range d.peers {
					d.output <- peer
				}
			}
		}
	}()

	return d.output, nil
}

// Stop implements resolver.Resolver.
func (d *Resolver) Stop() {
	d.stop <- struct{}{}
	close(d.output)
}
