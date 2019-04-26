package namespaces

import (
	"github.com/sirupsen/logrus"
	coreV1 "k8s.io/client-go/kubernetes/typed/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
)

func GetNamespaces() (ocpNamespaces []string) {
	logrus.Info("Getting namespaces")
	conf := "/Users/dima/.kube/config"
	config, err := clientcmd.BuildConfigFromFlags("", conf)
	if err != nil {
		panic(err.Error())
	}
	coreV1Client, err := coreV1.NewForConfig(config)
	namespaces, err := coreV1Client.Namespaces().List(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	for _, ocpNs := range namespaces.Items {
		ocpNamespaces = append(ocpNamespaces, ocpNs.Name)
	}
	// Append clusterrole manually
	ocpNamespaces = append(ocpNamespaces, "clusterrole")
	return
}
