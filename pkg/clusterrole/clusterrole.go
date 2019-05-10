package clusterrole

import (
	"github.com/dimss/skipper/pkg/clientcmdconfigs"
	"github.com/dimss/skipper/pkg/sankey"
	"github.com/sirupsen/logrus"
	rbacApiV1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	rbacV1 "k8s.io/client-go/kubernetes/typed/rbac/v1"
)

func (cr *ClusterRoleSankeyGraph) LoadK8SObjects() () {
	cr.clusterRoles = make(map[string][]rbacApiV1.ClusterRole)
	logrus.Info("Getting roles")

	rbacV1Client, err := rbacV1.NewForConfig(clientcmdconfigs.GetClientcmdConfigs())
	if err != nil {
		panic(err.Error())
	}

	roles, err := rbacV1Client.ClusterRoles().List(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	logrus.Infof("Total cluster roles in OCP cluster: %d", len(roles.Items))
	for _, ocpRole := range roles.Items {
		cr.clusterRoles[ocpRole.Namespace] = append(cr.clusterRoles[ocpRole.Namespace], ocpRole)
	}
	return
}

func (cr *ClusterRoleSankeyGraph) GetK8SObjects() interface{} {
	return cr.clusterRoles
}

func (cr *ClusterRoleSankeyGraph) CreateGraphData() {
	if (len(cr.clusterRoles) == 0) {
		logrus.Warn("Empty cluster roles list, gonna return empty results")
		return
	}
	cr.nodesIndexMap = cr.createNodesIndexMap()
	cr.Nodes = cr.createNodes()
	cr.Links = cr.createLinks()
}

func (cr *ClusterRoleSankeyGraph) createNodesIndexMap() (nodes map[string]int) {
	nodes = make(map[string]int)
	for _, ns := range cr.clusterRoles {
		for _, role := range ns {
			nodes["r-"+role.Name] = 1
			for _, rule := range role.Rules {
				for _, verb := range rule.Verbs {
					nodes["v-"+verb] = 1
				}
				for _, resource := range rule.Resources {
					nodes["re-"+resource] = 1
				}
			}
		}
	}
	return
}

func (cr *ClusterRoleSankeyGraph) createNodes() (nodes []sankey.Node) {
	i := 0
	nodes = make([]sankey.Node, len(cr.nodesIndexMap))
	for node := range cr.nodesIndexMap {
		nodes[i] = sankey.Node{node, node}
		i++
	}
	return
}

func (cr *ClusterRoleSankeyGraph) createLinks() (links []sankey.Link) {
	for _, ns := range cr.clusterRoles {
		for _, role := range ns {
			for _, rule := range role.Rules {
				for _, verb := range rule.Verbs {
					links = append(links, sankey.Link{
						Source: "v-" + verb,
						Target: "r-" + role.Name,
						Value:  1,
					})
				}
				for _, resource := range rule.Resources {
					links = append(links, sankey.Link{
						Source: "r-" + role.Name,
						Target: "re-" + resource,
						Value:  1,
					})
				}
			}
		}
	}
	return
}
