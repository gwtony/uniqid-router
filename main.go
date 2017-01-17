package main

import (
	"fmt"
	"time"
	"github.com/gwtony/gapi/api"
	"github.com/gwtony/uniqid_router/handler"
)

func main() {
	err := api.Init()
	if err != nil {
		fmt.Println("Init api failed")
		return
	}
	api.SetConfig("uniqid_router.conf")
	config := api.GetConfig()
	log := api.GetLog()

	err = handler.InitContext(config, log)
	if err != nil {
		fmt.Println("Init uniqid router failed")
		return
	}

	api.Run()
	time.Sleep(time.Second)
}
