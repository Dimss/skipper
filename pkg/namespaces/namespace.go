package namespaces

import (
	"github.com/dimss/skipper/pkg/clientcmdconfigs"
	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	coreV1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

func GetNamespaces() (ocpNamespaces []string) {
	logrus.Info("Getting namespaces")
	coreV1Client, err := coreV1.NewForConfig(clientcmdconfigs.GetClientcmdConfigs())
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
