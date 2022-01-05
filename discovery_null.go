package discovery

import "context"

var _ Discovery = (*NullDiscovery)(nil)

// NullDiscovery is a dummy resolver for static peers.
type NullDiscovery struct {
}

// NewNullDiscovery returns a new null resolver.
func NewNullDiscovery() *NullDiscovery {
	return &NullDiscovery{}
}

// Lookup implements discovery.Discovery.
func (d *NullDiscovery) Lookup() ([]string, error) {
	return []string{}, nil
}

// Register implements discovery.Discovery.
func (d *NullDiscovery) Register(ctx context.Context) error {
	return nil
}
