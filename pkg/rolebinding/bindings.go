package rolebinding

import (
	"github.com/sirupsen/logrus"
	"k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	rbacV1 "k8s.io/client-go/kubernetes/typed/rbac/v1"

	"k8s.io/client-go/tools/clientcmd"
)



func GetBindings(ns string) (ocpRoles []v1.RoleBinding) {

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
	roles, err := rbacV1Client.RoleBindings(ns).List(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	logrus.Infof("Total roles for namespace %s in OCP cluster: %d", ns, len(roles.Items))
	for _, ocpRole := range roles.Items {
		ocpRoles = append(ocpRoles, ocpRole)
	}
	return
}
