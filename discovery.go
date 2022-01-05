package discovery

import "context"

// Lookup represents a base interface for all lookups.
type Lookup interface {
	// Lookup returns the found peers. A peer is a ip:port.
	Lookup() ([]string, error)
}

// Register represents a base interface for all registers.
type Register interface {
	// Register registers self as a peer.
	Register(ctx context.Context) error
}
