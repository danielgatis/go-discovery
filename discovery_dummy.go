package discovery

import (
	"context"

	"github.com/sirupsen/logrus"
)

var _ Discovery = (*DummyDiscovery)(nil)

// DummyDiscovery is a dummy resolver for static peers.
type DummyDiscovery struct {
	peers   []string
	logger  logrus.FieldLogger
	output  chan []string
	stop    chan struct{}
	running bool
}

// NewDummyDiscovery returns a new dummy resolver.
func NewDummyDiscovery(peers []string, logger logrus.FieldLogger) *DummyDiscovery {
	return &DummyDiscovery{
		peers:  peers,
		logger: logger,
	}
}

// Register implements discovery.Discovery.
func (d *DummyDiscovery) Register(ctx context.Context) error {
	return nil
}

// Lookup implements discovery.Discovery.
func (d *DummyDiscovery) Lookup() ([]string, error) {
	return d.peers, nil
}
