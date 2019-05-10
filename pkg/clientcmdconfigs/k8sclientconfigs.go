package clientcmdconfigs

import (
	"github.com/spf13/viper"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func GetClientcmdConfigs() *rest.Config {
	conf := viper.GetString("kubeconfig")
	if conf == "useInClusterConfig" {
		config, err := rest.InClusterConfig()
		if err != nil {
			panic(err.Error())
		}
		return config
	} else {
		config, err := clientcmd.BuildConfigFromFlags("", conf)
		if err != nil {
			panic(err.Error())
		}
		return config
	}
}
