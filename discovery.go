package discovery

// Discovery represents a base interface for all providers.
type Discovery interface {
	// Start returns a channel that outputs the found peers. A peer is a ip:port.
	Start() (chan string, error)

	// Stop stops the resolver.
	Stop()
}
