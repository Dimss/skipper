package roles

import (
	"encoding/json"
	"github.com/dimss/skipper/pkg/clientcmdconfigs"
	"github.com/dimss/skipper/pkg/sankey"
	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	rbacV1 "k8s.io/client-go/kubernetes/typed/rbac/v1"
	"time"
)

func NewRolesSankeyGraph(ns string) *RoleSankeyGraph {
	d := RoleSankeyGraph{namespace: ns}
	d.loadRoleSankeyGraph()
	return &d
}

func (d *RoleSankeyGraph) LoadK8SObjects() () {
	if len(d.Nodes) != 0 && len(d.Links) != 0 {
		return
	}
	logrus.Info("Getting roles")
	rbacV1Client, err := rbacV1.NewForConfig(clientcmdconfigs.GetClientcmdConfigs())
	if err != nil {
		panic(err.Error())
	}
	roles, err := rbacV1Client.Roles(d.namespace).List(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	logrus.Infof("Total roles for namespace %s in OCP cluster: %d", d.namespace, len(roles.Items))
	for _, ocpRole := range roles.Items {
		d.roles = append(d.roles, ocpRole)
	}
	return
}

func (d *RoleSankeyGraph) GetK8SObjects() interface{} {
	return d.roles
}

func (d *RoleSankeyGraph) CreateGraphData() {
	if len(d.Nodes) != 0 && len(d.Links) != 0 {
		return
	}
	if (len(d.roles) == 0) {
		logrus.Warn("Empty roles list, gonna return empty results")
		return
	}
	d.nodesIndexMap = d.createNodesIndexMap()
	d.Nodes = d.createNodes()
	d.Links = d.createLinks()
	d.cacheRoleSankeyGraph()

}

func (d *RoleSankeyGraph) createNodesIndexMap() (nodes map[string]int) {

	nodes = make(map[string]int)
	for _, role := range d.roles {
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
	return
}

func (d *RoleSankeyGraph) createNodes() (nodes []sankey.Node) {
	i := 0
	nodes = make([]sankey.Node, len(d.nodesIndexMap))
	for node := range d.nodesIndexMap {
		nodes[i] = sankey.Node{node, node}
		i++
	}
	return
}

func (d *RoleSankeyGraph) createLinks() (links []sankey.Link) {

	for _, role := range d.roles {
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
	return
}

func (d *RoleSankeyGraph) cacheRoleSankeyGraph() {
	rolesAsJsonString, err := json.Marshal(d)
	if err != nil {
		panic(err)
	}
	redisClient = redis.NewClient(&redis.Options{Addr: viper.GetString("redis.conn")})
	err = redisClient.Set("roles", rolesAsJsonString, time.Second * viper.GetDuration("redis.ttl")).Err()
	if err != nil {
		panic(err)
	}
}

func (d *RoleSankeyGraph) loadRoleSankeyGraph() {
	redisClient = redis.NewClient(&redis.Options{Addr: viper.GetString("redis.conn")})
	redisCmd := redisClient.Get("roles")
	if redisCmd.Val() != "" {
		logrus.Info("Reading roles from cache")
		err := json.Unmarshal([]byte(redisCmd.Val()), d)
		if err != nil {
			panic(err)
		}
	}
}

func getRole(roleName string, ns string) {
	rbacV1Client, err := rbacV1.NewForConfig(clientcmdconfigs.GetClientcmdConfigs())
	if err != nil {
		panic(err.Error())
	}
	if role, err := rbacV1Client.Roles(ns).Get(roleName, metav1.GetOptions{}); err != nil {
		logrus.Info(role)
	}

}
