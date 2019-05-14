package roles

import (
	"encoding/json"
	"github.com/dimss/skipper/pkg/sankey"
	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	rbacApiV1 "k8s.io/api/rbac/v1"

	"time"
)

var watchRoleRedisKey = "watchRoles"
var redisClient = redis.NewClient(&redis.Options{Addr: viper.GetString("redis.conn")})

func getRedisKey() (key string) {
	timestamp := time.Now()
	return "roles-" + timestamp.String()
}

func storeRoles() {
	var sankeyGraph sankey.Sankey
	sankeyGraph = NewRolesSankeyGraph("")
	sankeyGraph.LoadK8SObjects()
	sankeyGraph.CreateGraphData()

	roles := sankeyGraph.GetK8SObjects()

	if response, err := json.Marshal(roles); err == nil {
		err := redisClient.Set("key", response, time.Second*60).Err()
		if err != nil {
			panic(err)
		}
	}
	skg := RoleSankeyGraph{}
	res := redisClient.Get("key")
	res1 := res.Val()
	json.Unmarshal([]byte(res1), &skg.roles)
	logrus.Info(res1)

	logrus.Info("Roles sankey data loaded")
}

func getWatchRolesRequest() (watchRoles []WatchRoleRequest) {
	watchRoles = []WatchRoleRequest{}
	redisCmd := redisClient.Get("watchRoles")
	if redisCmd.Val() == "" {
		return
	}
	if err := json.Unmarshal([]byte(redisCmd.Val()), &watchRoles); err != nil {
		logrus.Error("Enable to unmarshal JSON string from redis")
		panic(err)
	}
	return
}

func setWatchRolesRequest(role WatchRoleRequest) {
	watchRoles := getWatchRolesRequest()
	watchRoles = append(watchRoles, role)
	if jsonString, err := json.Marshal(watchRoles); err == nil {
		err := redisClient.Set(watchRoleRedisKey, jsonString, 0).Err()
		if err != nil {
			panic(err)
		}
	}
}

func storeRole(role rbacApiV1.Role) {
	roles := []rbacApiV1.Role{}
	redisCmd := redisClient.Get("roles")
	if redisCmd.Val() == "" {
		roles = append(roles, role)
		jsonString, err := json.Marshal(roles)
		if err != nil {
			panic(err)
		}
		err = redisClient.Set("roles", jsonString, 0).Err()
		if err != nil {
			panic(err)
		}
	} else {
		if err := json.Unmarshal([]byte(redisCmd.Val()), &roles); err != nil {
			logrus.Error("Enable to unmarshal JSON string from redis")
			panic(err)
		}
		roles = append(roles, role)
	}
}

func watch() {

	for {
		watchRolesRequests := getWatchRolesRequest()
		if len(watchRolesRequests) == 0 {
			continue
		}
		for _, watchRequest := range watchRolesRequests {
			logrus.Info(watchRequest)
		}
		time.Sleep(time.Second * 2)
	}
}

func StartRolesWatcher(watchRolesChan chan WatchRoleRequest) {
	go watch()
	for role := range watchRolesChan {
		setWatchRolesRequest(role)
	}
}
