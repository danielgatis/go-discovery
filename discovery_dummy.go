package discovery

import (
	"time"

	"github.com/sirupsen/logrus"
)

// DummyDiscovery is a dummy resolver for static peers.
type DummyDiscovery struct {
	peers    []string
	interval time.Duration
	logger   logrus.FieldLogger
	output   chan []string
	stop     chan struct{}
}

// NewDummyDiscovery returns a new dummy resolver.
func NewDummyDiscovery(peers []string, interval time.Duration, logger logrus.FieldLogger) *DummyDiscovery {
	return &DummyDiscovery{
		peers:    peers,
		interval: interval,
		logger:   logger,
		output:   make(chan []string),
		stop:     make(chan struct{}),
	}
}

// Start implements resolver.Resolver.
func (d *DummyDiscovery) Start() (chan []string, error) {
	ticker := time.NewTicker(d.interval)

	go func() {
		for {
			select {
			case <-d.stop:
				ticker.Stop()
				return
			case <-ticker.C:
				d.output <- d.peers
			}
		}
	}()

	return d.output, nil
}

// Stop implements resolver.Resolver.
func (d *DummyDiscovery) Stop() {
	d.stop <- struct{}{}
	close(d.output)
}
