package discovery

import (
	"context"
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
func NewDummyDiscovery(peers []string, logger logrus.FieldLogger) *DummyDiscovery {
	return &DummyDiscovery{
		peers:  peers,
		logger: logger,
	}
}

// Lookup implements discovery.Lookup.
func (d *DummyDiscovery) Lookup(ctx context.Context) ([]string, error) {
	return d.peers, nil
}
