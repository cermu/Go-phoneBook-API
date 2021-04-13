package utils

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v7"
	"github.com/spf13/viper"
	"log"
	"net/http"
)

var vp = viper.New()
var redisClient *redis.Client

/*
init function is called only once, we use it to read the configs file
and initialize a redis client before the applications comes up fully
*/
func init() {
	vp.SetConfigName("dev_env")                       // config file name without extension
	vp.SetConfigType("yaml")                          // config file type
	vp.AddConfigPath("./conf/")                       // . is the root dir of the app
	vp.AddConfigPath("/opt/goApps/conf/phoneBookAPI") // you can have multiple config paths
	vp.AutomaticEnv()                                 // read values from ENV variable

	err := vp.ReadInConfig()
	if err != nil {
		log.Fatalf("ERROR | Reading application's config file failed with message: %v\n", err.Error())
	}

	// initialize a redis client
	redisHost := vp.GetString("REDIS.HOST")
	redisPort := vp.GetInt("REDIS.PORT")
	redisDNS := fmt.Sprintf("%s:%d", redisHost, redisPort)
	redisClient = redis.NewClient(&redis.Options{
		Addr: redisDNS,
	})

	_, redisErr := redisClient.Ping().Result()
	if redisErr != nil {
		log.Fatalf("ERROR | Redis client initialization failed with message: %v\n", redisErr.Error())
	}
}

// ReadConfigs public function that returns a pointer to a viper object
func ReadConfigs() *viper.Viper {
	return vp
}

// RedisClient public function that returns a pointer to a redis client
func RedisClient() *redis.Client {
	return redisClient
}

// Message public function builds json messages
func Message(code int32, description string) map[string]interface{} {
	return map[string]interface{}{"response_code": code, "response_description": description}
}

// Respond public function responds with json message
func Respond(w http.ResponseWriter, response map[string]interface{}) {
	w.Header().Add("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(response)
}
