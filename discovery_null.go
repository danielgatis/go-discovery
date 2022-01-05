package discovery

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
