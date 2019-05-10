package roles

import (
	"github.com/dimss/skipper/pkg/sankey"
	rbacApiV1 "k8s.io/api/rbac/v1"
)

type RoleSankeyGraph struct {
	Nodes         []sankey.Node `json:"nodes"`
	Links         []sankey.Link `json:"links"`
	nodesIndexMap map[string]int
	namespace     string
	roles         []rbacApiV1.Role
}

type WatchRoleRequest struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}
