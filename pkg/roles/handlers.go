package roles

import (
	"encoding/json"
	"github.com/dimss/skipper/pkg/sankey"
	"net/http"
)

func GetRolesHandler(w http.ResponseWriter, r *http.Request) {
	var sankeyGraph sankey.Sankey
	ns := r.URL.Query().Get("ns")
	sankeyGraph = NewRolesSankeyGraph(ns)
	sankeyGraph.LoadK8SObjects()
	sankeyGraph.CreateGraphData()
	if response, err := json.Marshal(sankeyGraph); err == nil {
		if _, err := w.Write(response); err != nil {
			panic(err.Error())
		}
	}
}

