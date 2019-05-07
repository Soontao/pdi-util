package types

// GrammerNode type
type GrammerNode map[string]interface{}

// GetType type
func (n *GrammerNode) GetType() string {
	return (*n)["Type"].(string)
}
