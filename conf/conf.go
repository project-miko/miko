package conf

import (
	"fmt"
	"math/rand"
	"os"

	"gopkg.in/ini.v1"
)

var (
	configFp *ini.File

	ServerID   int64
	ServerName string

	Env = EnvDev
)

func GetConfigInt1(section, name string) (int, error) {
	return configFp.Section(section).Key(name).Int()
}

func GetConfigInt(section, name string) (int64, error) {
	return configFp.Section(section).Key(name).Int64()
}

func GetConfigString(section, name string) string {
	return configFp.Section(section).Key(name).String()
}

func ParseConfigINI(cpath string) (err error) {

	configFp, err = ini.Load(cpath)
	if err != nil {
		return err
	}

	ServerID = rand.Int63()

	hostName, e := os.Hostname()

	if e != nil {
		hostName = fmt.Sprintf("no-hostname-%d", ServerID)
	}

	ServerName = fmt.Sprintf("%s:%s",
		hostName,
		GetConfigString("ws", "port"))

	if appEnv := GetConfigString("app", "environment"); appEnv != "" {
		Env = appEnv
	}

	// TwitterAPIToken = GetConfigString("twitter", "api_token")
	// if TwitterAPIToken == "" {
	// 	return fmt.Errorf("twitter api_token is not config")
	// }
	// BaseScope = GetConfigString("twitter", "base_scope")
	// if BaseScope == "" {
	// 	return fmt.Errorf("twitter base_scope is not config")
	// }

	// LLMSavePath = GetConfigString("llm", "save_path")
	// if len(LLMSavePath) == 0 {
	// 	panic("config error")
	// }

	// TwMetricRefreshInterval, err = GetConfigInt("twitter", "tw_metric_refresh_interval")
	// if err != nil {
	// 	panic("config error")
	// }

	// TwitterOAuth2JumpFrontUrl = GetConfigString("twitter", "jump_front_url")
	// if len(TwitterOAuth2JumpFrontUrl) == 0 {
	// 	panic("config error")
	// }

	return nil
}
