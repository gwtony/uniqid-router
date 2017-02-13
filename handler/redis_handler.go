package handler

import (
	"time"
	"gopkg.in/redis.v5"
	"github.com/gwtony/gapi/log"
)

// Handler Redis handler
type RedisHandler struct {
	raddr      []string
	raddrSize  int
	loc        string
	log        log.Log
	rclient    *redis.Client
	ch         chan *RedisMessage
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
		//PoolTimeout: 10,
		PoolSize: 40,
	})
	h.ch = make(chan *RedisMessage, 100000)

	i := 0
	for {
		i++
		go h.Run()
		if i >= 40 {
			break
		}
	}

	return h
}

func (rh *RedisHandler) Run() {
	var rm *RedisMessage
	for {
		rm = <-rh.ch
		rh.Set(rm.Key, rm.Value, UROUTER_DEFAULT_TTL)
	}
}

func (rh *RedisHandler) RedisSet(key string, value []byte) {
	var rm RedisMessage
	rm.Key = key
	rm.Value = value
	rh.ch<-&rm
}

func (rh *RedisHandler) Set(key string, value []byte, ttl int) error {
	rh.log.Debug("key is %s", key)

	//No expire, just for test
	//err := rh.rclient.Set(key, value, 0).Err()
	err := rh.rclient.Set(key, string(value), time.Duration(ttl) * time.Second).Err()
	if err != nil {
		rh.log.Error("Set to redis failed", err)
	}

	//For debug
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

