package rolebinding

import (
	"github.com/dimss/skipper/pkg/sankey"
	"k8s.io/api/rbac/v1"
)

type RolesBindingSankeyGraph struct {
	Nodes         []sankey.Node `json:"nodes"`
	Links         []sankey.Link `json:"links"`
	nodesIndexMap map[string]int
	namespace     string
	rolesBindings []v1.RoleBinding
}
