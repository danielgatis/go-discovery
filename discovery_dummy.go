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
	running  bool
}

// NewDummyDiscovery returns a new dummy resolver.
func NewDummyDiscovery(peers []string, interval time.Duration, logger logrus.FieldLogger) *DummyDiscovery {
	return &DummyDiscovery{
		peers:    peers,
		interval: interval,
		logger:   logger,
		output:   make(chan []string),
		stop:     make(chan struct{}),
		running:  false,
	}
}

// Start implements resolver.Resolver.
func (d *DummyDiscovery) Start() (chan []string, error) {
	if d.running {
		return d.output, nil
	}

	ticker := time.NewTicker(d.interval)

	f := func() {
		d.output <- d.peers
	}

	go func() {
		f()

		for {
			select {
			case <-d.stop:
				ticker.Stop()
				return
			case <-ticker.C:
				f()
			}
		}
	}()

	d.running = true
	return d.output, nil
}

// Stop implements resolver.Resolver.
func (d *DummyDiscovery) Stop() {
	if !d.running {
		return
	}

	d.stop <- struct{}{}
	close(d.output)
	d.running = false
}
