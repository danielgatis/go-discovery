package discovery

// NullDiscovery is a dummy resolver for static peers.
type NullDiscovery struct {
	output  chan []string
	stop    chan struct{}
	running bool
}

// NewNullDiscovery returns a new null resolver.
func NewNullDiscovery() *NullDiscovery {
	return &NullDiscovery{
		output:  make(chan []string),
		stop:    make(chan struct{}),
		running: false,
	}
}

// Start implements resolver.Resolver.
func (d *NullDiscovery) Start() (chan []string, error) {
	if d.running {
		return d.output, nil
	}

	d.output = make(chan []string)
	d.running = true
	return d.output, nil
}

// Stop implements resolver.Resolver.
func (d *NullDiscovery) Stop() {
	if !d.running {
		return
	}

	d.stop <- struct{}{}
	close(d.output)
	d.running = false
}
