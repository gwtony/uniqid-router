package handler

import (
	"fmt"
	"os"
	"time"
	"strings"
	"github.com/gwtony/gapi/config"
	"github.com/gwtony/gapi/errors"
)

// URouterConfig URouter config
type URouterConfig struct {
	raddr      []string /* redis addr */

	apiLoc     string   /* urouter api location */

	timeout    time.Duration

	token      string   /* access token */
}

// ParseConfig parses config
func (conf *URouterConfig) ParseConfig(cf *config.Config) error {
	var err error
	if cf.C == nil {
		return errors.BadConfigError
	}
	raddrStr, err := cf.C.GetString("urouter", "redis_addr")
	if err != nil {
		fmt.Fprintln(os.Stderr, "[Error] [urouter] Read conf: No redis_addr")
		return err
	}
	if raddrStr == "" {
		fmt.Fprintln(os.Stderr, "[Error] [urouter] Empty redis server address")
		return errors.BadConfigError
	}
	raddr := strings.Split(raddrStr, ",")
	for i := 0; i < len(raddr); i++ {
		if raddr[i] != "" {
			if !strings.Contains(raddr[i], ":") {
				conf.raddr = append(conf.raddr, raddr[i] + ":" + DEFAULT_REDIS_PORT)
			} else {
				conf.raddr = append(conf.raddr, raddr[i])
			}
		}
	}

	conf.apiLoc, err = cf.C.GetString("urouter", "api_location")
	if err != nil {
		fmt.Fprintln(os.Stderr, "[Info] [URouter] Read conf: No api_location, use default location", UROUTER_DEFAULT_LOC)
		conf.apiLoc = UROUTER_DEFAULT_LOC
	}

	timeout, err := cf.C.GetInt64("urouter", "timeout")
	if err != nil || timeout <= 0 {
		fmt.Fprintln(os.Stderr, "[Info] [URouter] Read conf: use default purge_timeout: ", UROUTER_DEFAULT_TIMEOUT)
		timeout = UROUTER_DEFAULT_TIMEOUT
	}
	conf.timeout =  time.Duration(timeout) * time.Second

	conf.token, err = cf.C.GetString("urouter", "token")
	if err != nil {
		fmt.Fprintln(os.Stderr, "[Info] [URouter] Read conf: No token")
		conf.token = ""
	}

	return nil
}
