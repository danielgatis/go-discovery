package dummy

import (
	"time"

	"github.com/sirupsen/logrus"
)

// Discovery is a dummy resolver for static peers.
type Discovery struct {
	peers    []string
	logger   logrus.FieldLogger
	interval time.Duration
	output   chan string
	stop     chan struct{}
}

// New returns a new dummy resolver.
func New(peers []string, opts ...Option) *Discovery {
	const (
		defaultInterval = 2 * time.Second
	)

	var (
		defaultLogger = logrus.StandardLogger()
	)

	d := &Discovery{
		peers:    peers,
		logger:   defaultLogger,
		interval: defaultInterval,
		output:   make(chan string),
		stop:     make(chan struct{}),
	}

	for _, opt := range opts {
		opt(d)
	}

	return d
}

// Start implements resolver.Resolver.
func (d *Discovery) Start() (chan string, error) {
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
func (d *Discovery) Stop() {
	d.stop <- struct{}{}
	close(d.output)
}
