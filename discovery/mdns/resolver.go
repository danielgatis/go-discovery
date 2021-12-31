package mdns

import (
	"context"
	"fmt"
	"time"

	"github.com/grandcat/zeroconf"
	"github.com/sirupsen/logrus"
)

// Resolver is a mDNS resolver.
type Resolver struct {
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
func New(instance, service string, port int, opts ...Option) *Resolver {
	const (
		defaultDomain   = "local."
		defaultInterval = 5 * time.Second
	)

	var (
		defaultLogger = logrus.StandardLogger()
	)

	r := &Resolver{
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
		opt(r)
	}

	return r
}

// Start implements resolver.Resolver.
func (r *Resolver) Start() (chan string, error) {
	server, err := zeroconf.Register(r.instance, r.service, r.domain, r.port, []string{"txtv=0", "lo=1", "la=2"}, nil)
	if err != nil {
		return nil, fmt.Errorf("zeroconf.Register(...): %w", err)
	}

	resolver, err := zeroconf.NewResolver(nil)
	if err != nil {
		return r.output, fmt.Errorf("zeroconf.NewResolver(...): %w", err)
	}

	entries := make(chan *zeroconf.ServiceEntry)

	go func() {
		for entry := range entries {
			r.output <- fmt.Sprintf("%s:%d", entry.AddrIPv4[0], entry.Port)
		}
	}()

	ticker := time.NewTicker(r.interval)
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		for {
			select {
			case <-r.stop:
				ticker.Stop()
				cancel()
				server.Shutdown()
				return
			case <-ticker.C:
				if err := resolver.Browse(ctx, r.service, r.domain, entries); err != nil {
					r.logger.Errorf("Error during mDNS lookup: %v\n", err)
				}
			}
		}
	}()

	return r.output, nil
}

// Stop implements resolver.Resolver.
func (r *Resolver) Stop() {
	r.stop <- struct{}{}
	close(r.output)
}
