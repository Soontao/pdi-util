package types

import "reflect"

// GrammerNode type
type GrammerNode map[string]interface{}

// GetType type
func (n *GrammerNode) GetType() string {
	return (*n)["Type"].(string)
}

// GetNode by name
func (n *GrammerNode) GetNode(nodeName string) *GrammerNode {
	targetNode := (*n)[nodeName]
	if targetNode != nil {
		if reflect.TypeOf(targetNode) == reflect.TypeOf(&GrammerNode{}) {
			return targetNode.(*GrammerNode)
		}
		return nil
	}
	return nil
}

// GetNodeList by name
func (n *GrammerNode) GetNodeList(field string) []*GrammerNode {
	target := (*n)[field]
	if target != nil {
		if reflect.TypeOf(target) == reflect.TypeOf(&[]*GrammerNode{}) {
			return *(target.(*[]*GrammerNode))
		}
		return nil
	}
	return nil
}

// GetString value
func (n *GrammerNode) GetString(field string) string {
	return (*n)[field].(string)
}
