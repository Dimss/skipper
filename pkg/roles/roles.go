package roles

import (
	"github.com/dimss/skipper/pkg/sankey"
	"github.com/sirupsen/logrus"
	rbacApiV1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	rbacV1 "k8s.io/client-go/kubernetes/typed/rbac/v1"
	"k8s.io/client-go/tools/clientcmd"
)

func NewRolesSankeyGraph(ns string) *RoleSankeyGraph {
	d := RoleSankeyGraph{namespace: ns}
	return &d
}

func (d *RoleSankeyGraph) LoadK8SObjects() () {
	d.roles = make(map[string][]rbacApiV1.Role)
	logrus.Info("Getting roles")
	conf := "/Users/dima/.kube/config"
	config, err := clientcmd.BuildConfigFromFlags("", conf)

	if err != nil {
		panic(err.Error())
	}

	rbacV1Client, err := rbacV1.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	roles, err := rbacV1Client.Roles(d.namespace).List(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	logrus.Infof("Total roles for namespace %s in OCP cluster: %d", d.namespace, len(roles.Items))
	for _, ocpRole := range roles.Items {
		d.roles[ocpRole.Namespace] = append(d.roles[ocpRole.Namespace], ocpRole)
	}

	return
}

func (d *RoleSankeyGraph) CreateGraphData() {
	if (len(d.roles) == 0) {
		logrus.Warn("Empty roles list, gonna return empty results")
		return
	}
	d.nodesIndexMap = d.createNodesIndexMap()
	d.Nodes = d.createNodes()
	d.Links = d.createLinks()
}

func (d *RoleSankeyGraph) createNodesIndexMap() (nodes map[string]int) {
	i := 0
	nodes = make(map[string]int)
	for _, ns := range d.roles {
		for _, role := range ns {
			if _, ok := nodes[role.Name]; !ok {
				nodes[role.Name] = i
				i++
			}
			for _, rule := range role.Rules {
				for _, verb := range rule.Verbs {
					if _, ok := nodes[verb]; !ok {
						nodes[verb] = i
						i++
					}
				}
				for _, resource := range rule.Resources {
					if _, ok := nodes[resource]; !ok {
						nodes[resource] = i
						i++
					}
				}
			}
		}
	}
	return
}

func (d *RoleSankeyGraph) createNodes() (nodes []sankey.Node) {
	nodes = make([]sankey.Node, len(d.nodesIndexMap))
	for node, idx := range d.nodesIndexMap {
		nodes[idx] = sankey.Node{node, node}
	}
	return
}

func (d *RoleSankeyGraph) createLinks() (links []sankey.Link) {

	for _, ns := range d.roles {
		for _, role := range ns {
			for _, rule := range role.Rules {
				for _, verb := range rule.Verbs {
					links = append(links, sankey.Link{
						Source: verb,
						Target: role.Name,
						Value:  1,
					})
				}
				for _, resource := range rule.Resources {
					links = append(links, sankey.Link{
						Source: role.Name,
						Target: resource,
						Value:  1,
					})
				}
			}
		}
	}
	return
}
