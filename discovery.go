package discovery

import "context"

// Discovery represents a base interface for all resolvers.
type Discovery interface {
	// Lookup returns the found peers. A peer is a ip:port.
	Lookup() ([]string, error)

	// Register registers self as a peer.
	Register(ctx context.Context) error
}
