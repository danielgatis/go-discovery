package discovery

import "context"

// NullDiscovery is a dummy resolver for static peers.
type NullDiscovery struct {
}

// NewNullDiscovery returns a new null resolver.
func NewNullDiscovery() *NullDiscovery {
	return &NullDiscovery{}
}

// Lookup implements discovery.Lookup.
func (d *NullDiscovery) Lookup() ([]string, error) {
	return []string{}, nil
}

// Register implements discovery.Register.
func (d *NullDiscovery) Register(ctx context.Context) error {
	return nil
}
