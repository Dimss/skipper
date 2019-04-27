package rolebinding

import (
	"encoding/json"
	"github.com/dimss/skipper/pkg/sankey"
	"net/http"
)

func GetRolesBindingsHandler(w http.ResponseWriter, r *http.Request) {
	var sankeyGraph sankey.Sankey
	ns := r.URL.Query().Get("ns")
	sankeyGraph = NewRolesBindingsSankeyGraph(ns)
	sankeyGraph.LoadK8SObjects()
	sankeyGraph.CreateGraphData()
	if response, err := json.Marshal(sankeyGraph); err == nil {
		w.Write(response)
	}
}
