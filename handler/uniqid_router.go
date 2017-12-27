package handler

import (
	"github.com/gwtony/logger"
	"github.com/gwtony/gapi/api"
	//"github.com/gwtony/gapi/hserver"
	"github.com/gwtony/gapi/config"
)

// InitContext inits uniqid router context
func InitContext(conf *config.Config, log logger.Log) error {
	cf := &URouterConfig{}
	err := cf.ParseConfig(conf, log)
	if err != nil {
		log.Error("uniqid router parse config failed")
		return err
	}

	rh := InitRedisHandler(cf.raddr, log)

	api.AddTcpHandler(&URouterHandler{rh: rh, log: log})
	//TODO: add some monitor handler
	//api.AddHttpHandler(apiLoc + READ_LOC, &ReadHandler{h: h, domain: domain, token: token, log: log})

	return nil
}


