package roles

import (
	"encoding/json"
	"github.com/dimss/skipper/pkg/sankey"
	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
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

func getWatchRoles() (watchRoles []WatchRoleRequest) {
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

func setWatchRoles(role WatchRoleRequest) {
	watchRoles := getWatchRoles()
	watchRoles = append(watchRoles, role)
	if jsonString, err := json.Marshal(watchRoles); err == nil {
		err := redisClient.Set(watchRoleRedisKey, jsonString, 0).Err()
		if err != nil {
			panic(err)
		}
	}
}

func watch() {
	for {
		watchRoles := getWatchRoles()
		if len(watchRoles) == 0 {
			continue
		}
		for _, role := range watchRoles {
			logrus.Info(role)
		}
		time.Sleep(time.Second * 2)
	}
}

func StartRolesWatcher(watchRolesChan chan WatchRoleRequest) {
	go watch()
	for role := range watchRolesChan {
		setWatchRoles(role)
	}
}
