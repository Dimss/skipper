package rbac

import (
	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	rbacV1 "k8s.io/client-go/kubernetes/typed/rbac/v1"
	"sync"

	"k8s.io/client-go/tools/clientcmd"
)

func getRolesForNs(rbacV1Client *rbacV1.RbacV1Client, nsChan chan string, rolesChan chan map[string][]string, wg *sync.WaitGroup) {
	for ns := range nsChan {
		roles, err := rbacV1Client.Roles(ns).List(metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}
		logrus.Infof("Total roles for namespace %s in OCP cluster: %d", ns, len(roles.Items))
		var nsRoles []string
		for _, ocpRole := range roles.Items {
			nsRoles = append(nsRoles, ocpRole.Name)
		}
		rolesChan <- map[string][]string{ns: nsRoles}
	}
	wg.Done()
}

func getClusterRoles(rbacV1Client *rbacV1.RbacV1Client) (nsRoles []string) {
	clusterRoles, err := rbacV1Client.ClusterRoles().List(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	for _, ocpRole := range clusterRoles.Items {
		nsRoles = append(nsRoles, ocpRole.Name)
	}
	return
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
	var wg sync.WaitGroup
	nsChan := make(chan string, len(namespaces))
	rolesChan := make(chan map[string][]string, len(namespaces))
	for _, ns := range namespaces {
		wg.Add(1)
		go getRolesForNs(rbacV1Client, nsChan, rolesChan, &wg)
		nsChan <- ns
	}
	close(nsChan)
	wg.Wait()
	close(rolesChan)
	for role := range rolesChan {
		for ns, roles := range role {
			ocpRoles[ns] = roles
		}
	}
	// Get cluster roles
	ocpRoles["clusterrole"] = getClusterRoles(rbacV1Client)
	return
}
