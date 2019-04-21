package rbac

import (
	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	rbacV1 "k8s.io/client-go/kubernetes/typed/rbac/v1"
	"sync"
	rbacApiV1 "k8s.io/api/rbac/v1"
	"k8s.io/client-go/tools/clientcmd"
)

func getRolesForNs(rbacV1Client *rbacV1.RbacV1Client, nsChan chan string, rolesChan chan map[string][]rbacApiV1.Role, wg *sync.WaitGroup) {
	roles2Namespace := map[string][]rbacApiV1.Role{}
	for ns := range nsChan {

		roles, err := rbacV1Client.Roles("").List(metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}
		logrus.Infof("Total roles for namespace %s in OCP cluster: %d", ns, len(roles.Items))
		for _, ocpRole := range roles.Items {
			roles2Namespace[ocpRole.Namespace] = append(roles2Namespace[ocpRole.Namespace], ocpRole)
		}
		rolesChan <- roles2Namespace
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

func GetRoles() (ocpRoles map[string][]rbacApiV1.Role) {
	ocpRoles = make(map[string][]rbacApiV1.Role)
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
	// Get Roles in all namespaces - empty namespace "" = --all-namespaces
	namespaces := []string{""}
	var wg sync.WaitGroup
	nsChan := make(chan string, len(namespaces))
	rolesChan := make(chan map[string][]rbacApiV1.Role, len(namespaces))
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
	//ocpRoles["clusterrole"] = getClusterRoles(rbacV1Client)
	return
}
