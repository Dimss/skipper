package clusterrolebinding

import (
	"encoding/json"
	"github.com/dimss/skipper/pkg/sankey"
	"net/http"
)

func GetClusterRolesBindingsHandler(w http.ResponseWriter, r *http.Request) {
	var sankeyGraph sankey.Sankey
	sankeyGraph = &ClusterRolesBindingSankeyGraph{}
	sankeyGraph.LoadK8SObjects()
	sankeyGraph.CreateGraphData()
	if response, err := json.Marshal(sankeyGraph); err == nil {
		w.Write(response)
	}
}
