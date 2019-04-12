package rbac

import (
	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	rbacV1 "k8s.io/client-go/kubernetes/typed/rbac/v1"

	"k8s.io/client-go/tools/clientcmd"
)

func getRolesForNs(nsChan chan []string, roles chan map[string][]string) {

}

func GetRoles() (ocpRoles map[string][]string) {
	ocpRoles = make(map[string][]string)
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
	// Get Roles
	namespaces := GetNamespaces()
	for _, ns := range namespaces {
		roles, err := rbacV1Client.Roles(ns).List(metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}
		logrus.Infof("Total roles for namespace %s in OCP cluster: %d", ns, len(roles.Items))
		var nsRoles []string
		for _, ocpRole := range roles.Items {
			nsRoles = append(nsRoles, ocpRole.Name)
		}
		ocpRoles[ns] = nsRoles
	}
	// Get ClusterRoles
	clusterRoles, err := rbacV1Client.ClusterRoles().List(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	var nsRoles []string
	for _, ocpRole := range clusterRoles.Items {
		nsRoles = append(nsRoles, ocpRole.Name)
	}
	ocpRoles["clusterrole"] = nsRoles
	return
}
