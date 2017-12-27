package handler

import (
	"time"
	"gopkg.in/redis.v5"
	"github.com/gwtony/logger"
)

// Handler Redis handler
type RedisHandler struct {
	raddr      []string
	raddrSize  int
	loc        string
	log        logger.Log
	rclient    *redis.Client
	ch         chan *RedisMessage
}

// InitHandler inits redis handler
func InitRedisHandler(raddr []string, log logger.Log) *RedisHandler {
	h := &RedisHandler{}
	h.raddr = raddr
	h.raddrSize = len(raddr)
	h.log = log
	h.rclient = redis.NewClient(&redis.Options{
		Addr:     h.raddr[0],
		Password: "", // no password set
		DB:       0,  // use default DB
		//use default dial timeout 5s, read and write timeout 3s
		//PoolTimeout: 10,
		//PoolSize: 10, default pool size is 10
	})
	h.ch = make(chan *RedisMessage, 1000)

	i := 0
	for {
		i++
		go h.Run()
		if i >= REDIS_POOL_SIZE {
			break
		}
	}

	return h
}

func (rh *RedisHandler) Run() {
	var rm *RedisMessage
	cur := time.Now()
	tsb := cur.UnixNano() / 1000000
	tse := cur.UnixNano() / 1000000
	for {
		rm = <-rh.ch
		cur = time.Now()
		tsb = cur.UnixNano() / 1000000
		rh.Set(rm.Key, rm.Value, UROUTER_DEFAULT_TTL)
		cur = time.Now()
		tse = cur.UnixNano() / 1000000
		rh.log.Debug("Redis set cost in ms", logger.Int64("cost_ms", tse - tsb))
	}
}

func (rh *RedisHandler) RedisSet(key string, value []byte) {
	var rm RedisMessage
	rm.Key = key
	rm.Value = value
	rh.ch<-&rm
}

func (rh *RedisHandler) Set(key string, value []byte, ttl int) error {
	rh.log.Debug("Set key", logger.String("key", key))

	//No expire, just for test
	//err := rh.rclient.Set(key, value, 0).Err()
	err := rh.rclient.Set(key, string(value), time.Duration(ttl) * time.Second).Err()
	if err != nil {
		rh.log.Error("Set to redis failed", logger.Err(err))
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

