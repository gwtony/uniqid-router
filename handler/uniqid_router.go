package handler

import (
	"github.com/gwtony/gapi/log"
	"github.com/gwtony/gapi/api"
	//"github.com/gwtony/gapi/hserver"
	"github.com/gwtony/gapi/config"
)

// InitContext inits uniqid router context
func InitContext(conf *config.Config, log log.Log) error {
	cf := &URouterConfig{}
	err := cf.ParseConfig(conf)
	if err != nil {
		log.Error("Macedon parse config failed")
		return err
	}

	rh := InitRedisHandler(cf.raddr, log)

	//apiLoc := cf.apiLoc
	//token  := cf.token


	api.AddTcpHandler(&URouterHandler{rh: rh, log: log})
	//api.AddHttpHandler(apiLoc + MACEDON_ADD_LOC, &AddHandler{h: h, domain: domain, pc: pc, token: token, log: log})
	//api.AddHttpHandler(apiLoc + MACEDON_DELETE_LOC, &DeleteHandler{h: h, domain: domain, pc: pc, token: token, log: log})
	//api.AddHttpHandler(apiLoc + MACEDON_READ_LOC, &ReadHandler{h: h, domain: domain, token: token, log: log})
	//api.AddHttpHandler(apiLoc + MACEDON_SCAN_LOC, &ScanHandler{h: h, domain: domain, token: token, log: log})
	//api.AddHttpHandler(apiLoc + MACEDON_ADD_SERVER_LOC, &AddServerHandler{h: h, pc: pc, token: token, log: log})
	//api.AddHttpHandler(apiLoc + MACEDON_DELETE_SERVER_LOC, &DeleteServerHandler{h: h, pc: pc, token: token, log: log})
	//api.AddHttpHandler(apiLoc + MACEDON_READ_SERVER_LOC, &ReadServerHandler{h: h, token: token, log: log})

	return nil
}


