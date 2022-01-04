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
	interval time.Duration
	logger   logrus.FieldLogger
	output   chan []string
	stop     chan struct{}
}

// NewMdnsDiscovery returns a new mDNS resolver.
func NewMdnsDiscovery(instance, service, domain string, port int, interval time.Duration, logger logrus.FieldLogger) *MdnsDiscovery {
	return &MdnsDiscovery{
		instance: instance,
		service:  service,
		domain:   domain,
		port:     port,
		interval: interval,
		logger:   logger,
		output:   make(chan []string),
		stop:     make(chan struct{}),
	}
}

// Start implements resolver.Resolver.
func (d *MdnsDiscovery) Start() (chan []string, error) {
	server, err := zeroconf.Register(d.instance, d.service, d.domain, d.port, []string{"txtv=0", "lo=1", "la=2"}, nil)
	if err != nil {
		return nil, fmt.Errorf("zeroconf.Register(...): %w", err)
	}

	resolver, err := zeroconf.NewResolver(nil)
	if err != nil {
		return d.output, fmt.Errorf("zeroconf.NewResolver(...): %w", err)
	}

	ticker := time.NewTicker(d.interval)

	go func() {
		for {
			select {
			case <-d.stop:
				ticker.Stop()
				server.Shutdown()
				return
			case <-ticker.C:
				peers := make([]string, 0)

				ctx := context.Background()
				entries := make(chan *zeroconf.ServiceEntry)

				go func() {
					for entry := range entries {
						peers = append(peers, fmt.Sprintf("%s:%d", entry.AddrIPv4[0], entry.Port))
					}
				}()

				if err := resolver.Browse(ctx, d.service, d.domain, entries); err != nil {
					d.logger.Errorf("Error during mDNS lookup: %v\n", err)
				}

				<-ctx.Done()
				d.output <- peers
			}
		}
	}()

	return d.output, nil
}

// Stop implements resolver.Resolver.
func (d *MdnsDiscovery) Stop() {
	d.stop <- struct{}{}
	close(d.output)
}
