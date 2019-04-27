package rolebinding

import (
	"github.com/dimss/skipper/pkg/sankey"
	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	rbacV1 "k8s.io/client-go/kubernetes/typed/rbac/v1"

	"k8s.io/client-go/tools/clientcmd"
)

func NewRolesBindingsSankeyGraph(ns string) *RolesBindingSankeyGraph {
	rb := RolesBindingSankeyGraph{namespace: ns}
	return &rb
}

func (rb *RolesBindingSankeyGraph) LoadK8SObjects() {

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
	roles, err := rbacV1Client.RoleBindings(rb.namespace).List(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	logrus.Infof("Total roles for namespace %s in OCP cluster: %d", rb.namespace, len(roles.Items))
	for _, ocpRole := range roles.Items {
		rb.rolesBindings = append(rb.rolesBindings, ocpRole)
	}
	return
}

func (rb *RolesBindingSankeyGraph) CreateGraphData() {
	rb.nodesIndexMap = rb.createNodesIndexMap()
	rb.Nodes = rb.createNodes()
	rb.Links = rb.createLinks()
}

func (d *RolesBindingSankeyGraph) createNodesIndexMap() (nodes map[string]int) {
	i := 0
	nodes = make(map[string]int)
	for _, roleBinding := range d.rolesBindings {
		if _, ok := nodes[roleBinding.Name]; !ok {
			nodes[roleBinding.Name] = i
			i++
		}
		if _, ok := nodes[roleBinding.RoleRef.Name]; !ok {
			nodes[roleBinding.RoleRef.Name] = i
			i++
		}

		for _, subject := range roleBinding.Subjects {
			if _, ok := nodes[subject.Name]; !ok {
				nodes[subject.Name] = i
				i++
			}
		}

	}
	return
}

func (rb *RolesBindingSankeyGraph) createNodes() (nodes []sankey.Node) {
	nodes = make([]sankey.Node, len(rb.nodesIndexMap))
	for node, idx := range rb.nodesIndexMap {
		nodes[idx] = sankey.Node{node, node}
	}
	return
}

func (rb *RolesBindingSankeyGraph) createLinks() (links []sankey.Link) {

	for _, roleBinding := range rb.rolesBindings {
		links = append(links, sankey.Link{
			Source: roleBinding.RoleRef.Name,
			Target: roleBinding.Name,
			Value:  1,
		})
		for _, subject := range roleBinding.Subjects {
			links = append(links, sankey.Link{
				Source: subject.Name,
				Target: roleBinding.Name,
				Value:  1,
			})
		}
	}
	return
}
