package clusterrole

import (
	"github.com/dimss/skipper/pkg/sankey"
	rbacApiV1 "k8s.io/api/rbac/v1"
)

type ClusterRoleSankeyGraph struct {
	Nodes         []sankey.Node `json:"nodes"`
	Links         []sankey.Link `json:"links"`
	nodesIndexMap map[string]int
	clusterRoles  map[string][]rbacApiV1.ClusterRole
}
