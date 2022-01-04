package discovery

import (
	"context"
	"fmt"
	"time"

	"github.com/grandcat/zeroconf"
	"github.com/sirupsen/logrus"
)

// MdnsDiscovery is a mDNS resolver.
type MdnsDiscovery struct {
	instance string
	service  string
	domain   string
	port     int
	timeout  time.Duration
	logger   logrus.FieldLogger
}

// NewMdnsDiscovery returns a new mDNS resolver.
func NewMdnsDiscovery(instance, service, domain string, port int, logger logrus.FieldLogger) *MdnsDiscovery {
	return &MdnsDiscovery{
		instance: instance,
		service:  service,
		domain:   domain,
		port:     port,
		logger:   logger,
	}
}

// Register implements discovery.Register.
func (d *MdnsDiscovery) Register(ctx context.Context) error {
	server, err := zeroconf.Register(d.instance, d.service, d.domain, d.port, []string{"txtv=0", "lo=1", "la=2"}, nil)
	if err != nil {
		return fmt.Errorf("zeroconf.Register(...): %w", err)
	}

	go func() {
		<-ctx.Done()
		server.Shutdown()
	}()

	return nil
}

// Lookup implements discovery.Lookup.
func (d *MdnsDiscovery) Lookup(ctx context.Context) ([]string, error) {
	peers := make([]string, 0)
	entries := make(chan *zeroconf.ServiceEntry)

	go func() {
		for entry := range entries {
			peers = append(peers, fmt.Sprintf("%s:%d", entry.AddrIPv4[0], entry.Port))
		}
	}()

	resolver, err := zeroconf.NewResolver(nil)
	if err != nil {
		return peers, err
	}

	if err := resolver.Browse(ctx, d.service, d.domain, entries); err != nil {
		d.logger.Errorf("Error during mDNS lookup: %v\n", err)
	}

	<-ctx.Done()
	return peers, nil
}
