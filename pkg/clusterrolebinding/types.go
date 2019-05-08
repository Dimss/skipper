package clusterrolebinding

import (
	"github.com/dimss/skipper/pkg/sankey"
	"k8s.io/api/rbac/v1"
)

type ClusterRolesBindingSankeyGraph struct {
	Nodes         []sankey.Node `json:"nodes"`
	Links         []sankey.Link `json:"links"`
	nodesIndexMap map[string]int
	namespace     string
	clusterRolesBindings []v1.ClusterRoleBinding
}
