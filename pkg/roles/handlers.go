package roles

import (
	"encoding/json"
	"github.com/dimss/skipper/pkg/sankey"
	"io/ioutil"
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

func WatchHandler(w http.ResponseWriter, r *http.Request, watchRolesChan chan WatchRoleRequest) {
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	watchRoleRequest := WatchRoleRequest{}
	err = json.Unmarshal(b, &watchRoleRequest)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	watchRolesChan <- watchRoleRequest
}
