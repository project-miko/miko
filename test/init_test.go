package main

import (
	"fmt"
	"os"

	"github.com/project-miko/miko/conf"
	"github.com/project-miko/miko/core"
	"github.com/project-miko/miko/models"
	"github.com/project-miko/miko/sdk/twitterapi"
	"github.com/project-miko/miko/tools/log"
	"github.com/project-miko/miko/tools/logger"
)

// used to init configs
func initTester() {
	// initialize configuration
	err := conf.ParseConfigINI("../build/conf.ini")
	if err != nil {
		fmt.Println("err : parse config failed", err.Error())
		os.Exit(1)
	}

	logPath := "../storage/logs"

	logger.InitLogger(logPath, "20060102") // log
	el, err := conf.GetConfigInt("log", "level")
	if err != nil {
		el = 0 // if log path error, modify config file log_path to absolute path
	}
	log.SetLogErrorLevel(int(el))

	// initialize database
	models.InitModel()

	_loginKey := conf.GetConfigString("app", "login_key")
	if _loginKey == "" {
		panic(fmt.Errorf("login_key did not configured"))
	}
	core.LoginKey = []byte(_loginKey)

	//core.GetTwitterAPIToken()
	twitterapi.InitConfig()
}
