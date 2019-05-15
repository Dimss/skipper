package main

import (
	"fmt"
	"github.com/dimss/skipper/pkg/clusterrole"
	"github.com/dimss/skipper/pkg/clusterrolebinding"
	"github.com/dimss/skipper/pkg/namespaces"
	"github.com/dimss/skipper/pkg/rolebinding"
	"github.com/dimss/skipper/pkg/roles"
	"github.com/dimss/skipper/pkg/uiserver"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var rootCmd = &cobra.Command{
	Use:   "skipper",
	Short: "View and manage K8S/OKD/OpenShift RBAC",
}

var webServerCmd = &cobra.Command{
	Use:   "server",
	Short: "Start API and UI server ",
	Run: func(cmd *cobra.Command, args []string) {

		watchRolesChan := make(chan roles.WatchRoleRequest, 100)
		// Start role watcher
		go roles.StartRolesWatcher(watchRolesChan)
		// Start web server
		logrus.Info("Starting up web server...")

		r := chi.NewRouter()

		cors := cors.New(cors.Options{
			AllowedOrigins:   []string{"*"},
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
			ExposedHeaders:   []string{"Link"},
			AllowCredentials: true,
			MaxAge:           300,
		})
		r.Use(cors.Handler)
		r.Get("/api/roles", roles.GetRolesHandler)
		r.Post("/api/roles/watch", func(w http.ResponseWriter, r *http.Request) {
			roles.WatchHandler(w, r, watchRolesChan)
		})
		r.Get("/api/bindings", rolebinding.GetRolesBindingsHandler)
		r.Get("/api/namespaces", namespaces.GetNamespacesHandler)
		r.Get("/api/clusterroles", clusterrole.GetClusterRolesHandler)
		r.Get("/api/clusterrolesbindings", clusterrolebinding.GetClusterRolesBindingsHandler)

		uiserver.UiServer(r, "/ui")

		http.ListenAndServe(viper.GetString("app.bind"), r)

	},
}

var dumpRuntimeConfigCmd = &cobra.Command{
	Use:   "dumpconfig",
	Short: "Dump all runtime configs",
	Run: func(cmd *cobra.Command, args []string) {
		logrus.Info("Dumping runtime configs")
		logrus.Infof("app.bind: %s", viper.GetString("app.bind"))
		logrus.Infof("redis.conn: %s", viper.GetString("redis.conn"))
	},
}
var watcherCmd = &cobra.Command{
	Use:   "watch",
	Short: "Watch for RBAC resource",
	Run: func(cmd *cobra.Command, args []string) {
		//roles.StartRolesWatcher()
	},
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringP("kubeconfig", "k", "", "Path to kubeconfig file, default to $home/.kube/config")
	rootCmd.PersistentFlags().StringP("configpath", "c", "", "Path to config directory with config.json file, default to . ")
	if err := viper.BindPFlag("kubeconfig", rootCmd.PersistentFlags().Lookup("kubeconfig")); err != nil {
		panic(err)
	}
	if err := viper.BindPFlag("configpath", rootCmd.PersistentFlags().Lookup("configpath")); err != nil {
		panic(err)
	}
	rootCmd.AddCommand(webServerCmd)
	rootCmd.AddCommand(dumpRuntimeConfigCmd)
	rootCmd.AddCommand(watcherCmd)
	// Init log
	logrus.SetOutput(os.Stdout)
	logrus.SetReportCaller(true)
	logrus.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
}

func initConfig() {
	viper.SetConfigType("json")
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetEnvPrefix("SKIPPER")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
	configPath := viper.GetString("configpath")
	logrus.Info(configPath)
	// If config flag is empty, assume config.json located in current directory
	if configPath != "" {
		viper.AddConfigPath(configPath)
	}
	// look for kubeconfig file, if not found, assume running inside OPC cluster
	kubeconfig := viper.GetString("kubeconfig")
	if kubeconfig == "" {
		// Check if kubeconfig file exists in user's HOME
		kubeconfig = filepath.Join(os.Getenv("HOME"), ".kube", "config")
		_, err := os.Stat(kubeconfig)
		if os.IsNotExist(err) {
			// The kubeconfig wasn't passed in and not found under user's home directory, assuming inClusterConfig mode
			logrus.Info("Unable to find kubeconfig, assuming running inside K8S cluster, gonna use inClusterConfig")
			viper.Set("kubeconfig", "useInClusterConfig")
		} else {
			// Use kubeconfig from user's home directory
			logrus.Info("Gonna use kubeconfig from user's HOME directory")
			viper.Set("kubeconfig", kubeconfig)
		}
	}
	err := viper.ReadInConfig()
	if err != nil {
		logrus.Errorf("Unable to read config.json file, err: %s", err)
		os.Exit(1)
	}
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
