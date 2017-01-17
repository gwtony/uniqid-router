package handler

import (
	//"time"
	"gopkg.in/redis.v5"
	//"github.com/gwtony/gapi/variable"
	"github.com/gwtony/gapi/log"
	//"github.com/gwtony/gapi/errors"
)

// Handler Redis handler
type RedisHandler struct {
	raddr      []string
	raddrSize  int
	loc        string
	log        log.Log
	rclient    *redis.Client
}

// InitHandler inits redis handler
func InitRedisHandler(raddr []string, log log.Log) *RedisHandler {
	h := &RedisHandler{}
	h.raddr = raddr
	h.raddrSize = len(raddr)
	h.log = log
	h.rclient = redis.NewClient(&redis.Options{
		Addr:     h.raddr[0],
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return h
}

func (rh *RedisHandler) Set(key string, value []byte, ttl int) error {
	//TODO: long connection
	rh.log.Debug("key is %s", key)
	//rh.log.Debug("value is %s", value)
	err := rh.rclient.Set(key, value, 0).Err()
	//err := client.Set(key, string(value), time.Duration(ttl) * time.Second).Err()
	if err != nil {
		rh.log.Error("Set to redis failed", err)
	}
	//val, err := client.Get(key).Result()
	//if err != nil {
	//	rh.log.Error("Get error", err)
	//} else {
	//	rh.log.Debug("Get value is ", string(val))
	//}
	//client.Close()
	return err
}

//TODO:
func (rh *RedisHandler) Get(key, value []byte, ttl int) error {
	return nil
}

