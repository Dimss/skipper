package clusterrole

import (
	"encoding/json"
	"github.com/dimss/skipper/pkg/sankey"
	"net/http"
)

func GetClusterRolesHandler(w http.ResponseWriter, r *http.Request) {
	var sankeyGraph sankey.Sankey
	sankeyGraph = &ClusterRoleSankeyGraph{}
	sankeyGraph.LoadK8SObjects()
	sankeyGraph.CreateGraphData()
	if response, err := json.Marshal(sankeyGraph); err == nil {
		if _, err := w.Write(response); err != nil {
			panic(err.Error())
		}
	}
}
