package rbac

import (
	userv1 "github.com/openshift/client-go/user/clientset/versioned/typed/user/v1"
	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
	"os"
)

func init() {
	logrus.SetOutput(os.Stdout)
	logrus.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
}

func GetUsers() (ocpUsers []string) {
	logrus.Info("Getting users")
	conf := "/Users/dima/.kube/config"
	config, err := clientcmd.BuildConfigFromFlags("", conf)
	if err != nil {
		panic(err.Error())
	}

	userv1Client, err := userv1.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	users, err := userv1Client.Users().List(metav1.ListOptions{})
	logrus.Infof("Total users in OCP cluster: %d", len(users.Items))
	for _, user := range users.Items {
		ocpUsers = append(ocpUsers, user.Name)
	}
	return
}

func GetServiceAccounts() {

	//for _, project := range projects.Items {
	//	projectsNames = append(projectsNames, project.Name)
	//}

}
