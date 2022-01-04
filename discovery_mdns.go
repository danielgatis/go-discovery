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
	running  bool
}

// NewMdnsDiscovery returns a new mDNS resolver.
func NewMdnsDiscovery(interval time.Duration, logger logrus.FieldLogger, configFn func() (instance, service, domain string, port int)) *MdnsDiscovery {
	instance, service, domain, port := configFn()

	return &MdnsDiscovery{
		instance: instance,
		service:  service,
		domain:   domain,
		port:     port,
		interval: interval,
		logger:   logger,
		output:   make(chan []string),
		stop:     make(chan struct{}),
		running:  false,
	}
}

// Start implements resolver.Resolver.
func (d *MdnsDiscovery) Start() (chan []string, error) {
	if d.running {
		return d.output, nil
	}

	server, err := zeroconf.Register(d.instance, d.service, d.domain, d.port, []string{"txtv=0", "lo=1", "la=2"}, nil)
	if err != nil {
		return d.output, fmt.Errorf("zeroconf.Register(...): %w", err)
	}

	f := func() error {
		peers := make([]string, 0)
		entries := make(chan *zeroconf.ServiceEntry)

		go func() {
			for entry := range entries {
				peers = append(peers, fmt.Sprintf("%s:%d", entry.AddrIPv4[0], entry.Port))
			}
		}()

		resolver, err := zeroconf.NewResolver(nil)
		if err != nil {
			return err
		}

		ctx, _ := context.WithTimeout(context.Background(), d.interval)
		if err := resolver.Browse(ctx, d.service, d.domain, entries); err != nil {
			d.logger.Errorf("Error during mDNS lookup: %v\n", err)
		}

		<-ctx.Done()
		d.output <- peers

		return nil
	}

	go func() {
		if err := f(); err != nil {
			d.logger.Errorf("Error during mdns service lookup: %v\n", err)
		}

		for {
			select {
			case <-d.stop:
				server.Shutdown()
				return
			default:
				if err := f(); err != nil {
					d.logger.Errorf("Error during mdns service lookup: %v\n", err)
				}
			}
		}
	}()

	d.running = true
	return d.output, nil
}

// Stop implements resolver.Resolver.
func (d *MdnsDiscovery) Stop() {
	if !d.running {
		return
	}

	d.stop <- struct{}{}
	close(d.output)
	d.running = false
}
