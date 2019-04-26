package sankey

import rbacApiV1 "k8s.io/api/rbac/v1"

func NewSankey(k8sObjects interface{}) *GraphData {
	d := GraphData{k8sObjects: k8sObjects}
	return &d
}

func (d *GraphData) CreateGraphData() {
	ocpRoles := d.k8sObjects.(map[string][]rbacApiV1.Role)
	d.nodesIndexMap = d.createNodesIndexMap(ocpRoles)
	d.Nodes = d.createNodes()
	d.Links = d.createLinks(ocpRoles)
}

func (d *GraphData) createNodesIndexMap(ocpRoles map[string][]rbacApiV1.Role) (nodes map[string]int) {
	i := 0
	nodes = make(map[string]int)
	for _, ns := range ocpRoles {
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

func (d *GraphData) createNodes() (nodes []Node) {
	nodes = make([]Node, len(d.nodesIndexMap))
	for node, idx := range d.nodesIndexMap {
		nodes[idx] = Node{node, node}
	}
	return
}

func (d *GraphData) createLinks(ocpRoles map[string][]rbacApiV1.Role) (links []Link) {

	for _, ns := range ocpRoles {
		for _, role := range ns {
			for _, rule := range role.Rules {
				for _, verb := range rule.Verbs {
					links = append(links, Link{
						Source: d.nodesIndexMap[verb],
						Target: d.nodesIndexMap[role.Name],
						Value:  1,
					})
				}
				for _, resource := range rule.Resources {
					links = append(links, Link{
						Source: d.nodesIndexMap[role.Name],
						Target: d.nodesIndexMap[resource],
						Value:  1,
					})
				}
			}
		}
	}
	return
}

func GetNodes(ocpRoles map[string][]rbacApiV1.Role) (nodes map[string]int, sunkeyNodes []Node) {
	i := 0
	nodes = make(map[string]int)
	for _, ns := range ocpRoles {
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
	sunkeyNodes = make([]Node, len(nodes))
	for node, idx := range nodes {
		sunkeyNodes[idx] = Node{node, node}
	}
	return
}

func GetLinks(ocpRoles map[string][]rbacApiV1.Role, nodes map[string]int) (links []Link) {

	for _, ns := range ocpRoles {
		for _, role := range ns {
			for _, rule := range role.Rules {
				for _, verb := range rule.Verbs {
					links = append(links, Link{
						Source: nodes[verb],
						Target: nodes[role.Name],
						Value:  1,
					})
				}
				for _, resource := range rule.Resources {
					links = append(links, Link{
						Source: nodes[role.Name],
						Target: nodes[resource],
						Value:  1,
					})
				}
			}
		}
	}
	return
}
