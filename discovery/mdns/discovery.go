package mdns

import (
	"context"
	"fmt"
	"time"

	"github.com/grandcat/zeroconf"
	"github.com/sirupsen/logrus"
)

// Discovery is a mDNS resolver.
type Discovery struct {
	instance string
	service  string
	domain   string
	port     int
	logger   logrus.FieldLogger
	interval time.Duration
	output   chan string
	stop     chan struct{}
}

// New returns a new mDNS resolver.
func New(instance, service string, port int, opts ...Option) *Discovery {
	const (
		defaultDomain   = "local."
		defaultInterval = 5 * time.Second
	)

	var (
		defaultLogger = logrus.StandardLogger()
	)

	d := &Discovery{
		instance: instance,
		service:  service,
		domain:   defaultDomain,
		port:     port,
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
	server, err := zeroconf.Register(d.instance, d.service, d.domain, d.port, []string{"txtv=0", "lo=1", "la=2"}, nil)
	if err != nil {
		return nil, fmt.Errorf("zeroconf.Register(...): %w", err)
	}

	resolver, err := zeroconf.NewResolver(nil)
	if err != nil {
		return d.output, fmt.Errorf("zeroconf.NewResolver(...): %w", err)
	}

	entries := make(chan *zeroconf.ServiceEntry)

	go func() {
		for entry := range entries {
			d.output <- fmt.Sprintf("%s:%d", entry.AddrIPv4[0], entry.Port)
		}
	}()

	ticker := time.NewTicker(d.interval)
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		for {
			select {
			case <-d.stop:
				ticker.Stop()
				cancel()
				server.Shutdown()
				return
			case <-ticker.C:
				if err := resolver.Browse(ctx, d.service, d.domain, entries); err != nil {
					d.logger.Errorf("Error during mDNS lookup: %v\n", err)
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
