package clusterrolebinding

import (
	"github.com/dimss/skipper/pkg/clientcmdconfigs"
	"github.com/dimss/skipper/pkg/sankey"
	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	rbacV1 "k8s.io/client-go/kubernetes/typed/rbac/v1"
)

func (crb *ClusterRolesBindingSankeyGraph) LoadK8SObjects() {

	logrus.Info("Getting roles")

	rbacV1Client, err := rbacV1.NewForConfig(clientcmdconfigs.GetClientcmdConfigs())
	if err != nil {
		panic(err.Error())
	}
	roles, err := rbacV1Client.ClusterRoleBindings().List(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	logrus.Infof("Total cluster roles for in OCP cluster: %d", len(roles.Items))
	for _, ocpRole := range roles.Items {
		crb.clusterRolesBindings = append(crb.clusterRolesBindings, ocpRole)
	}
	return
}

func (crb *ClusterRolesBindingSankeyGraph) GetK8SObjects() interface{} {
	return crb.clusterRolesBindings
}

func (crb *ClusterRolesBindingSankeyGraph) CreateGraphData() {
	if (len(crb.clusterRolesBindings) == 0) {
		logrus.Warn("Empty rolesbinding list, gonna return empty results")
		return
	}
	crb.nodesIndexMap = crb.createNodesIndexMap()
	crb.Nodes = crb.createNodes()
	crb.Links = crb.createLinks()
}

func (d *ClusterRolesBindingSankeyGraph) createNodesIndexMap() (nodes map[string]int) {
	nodes = make(map[string]int)
	for _, roleBinding := range d.clusterRolesBindings {
		nodes["rb-"+roleBinding.Name] = 1
		nodes["r-"+roleBinding.RoleRef.Name] = 1
		for _, subject := range roleBinding.Subjects {
			nodes["s-"+subject.Name] = 1
		}
	}
	return
}

func (rb *ClusterRolesBindingSankeyGraph) createNodes() (nodes []sankey.Node) {
	i := 0
	nodes = make([]sankey.Node, len(rb.nodesIndexMap))
	for node := range rb.nodesIndexMap {
		nodes[i] = sankey.Node{node, node}
		i++
	}
	return
}

func (rb *ClusterRolesBindingSankeyGraph) createLinks() (links []sankey.Link) {

	for _, roleBinding := range rb.clusterRolesBindings {
		links = append(links, sankey.Link{
			Source: "r-" + roleBinding.RoleRef.Name,
			Target: "rb-" + roleBinding.Name,
			Value:  1,
		})
		for _, subject := range roleBinding.Subjects {
			links = append(links, sankey.Link{
				Source: "rb-" + roleBinding.Name,
				Target: "s-" + subject.Name,
				Value:  1,
			})
		}
	}
	return
}
