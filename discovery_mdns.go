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

// Register implements discovery.Discovery.
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

// Lookup implements discovery.Discovery.
func (d *MdnsDiscovery) Lookup() ([]string, error) {
	peers := make([]string, 0)
	entries := make(chan *zeroconf.ServiceEntry)

	resolver, err := zeroconf.NewResolver(nil)
	if err != nil {
		return peers, err
	}

	go func() {
		ctxt, _ := context.WithTimeout(context.Background(), 1*time.Second)
		if err := resolver.Browse(ctxt, d.service, d.domain, entries); err != nil {
			d.logger.Errorf("Error during mDNS lookup: %v\n", err)
		}

		<-ctxt.Done()
	}()

	for entry := range entries {
		peers = append(peers, fmt.Sprintf("%s:%d", entry.AddrIPv4[0], entry.Port))
	}

	return peers, nil
}
