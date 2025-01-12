package conf

import "time"

// constants

const (
	EnvLocal = "local"
	EnvDev   = "dev"
	EnvTest  = "test"
	EnvPre   = "pre"
	EnvProd  = "prod"

	CrawlInterval = 60 // crawl interval, unit: minute
)

var (
	NullJsonArray = make([]interface{}, 0)

	TwitterAPIToken = ""
	BaseScope       = ""

	LLMSavePath = ""

	TwMetricRefreshInterval int64 = 0

	TimeZone       = time.FixedZone("UTC", 0)
	NewTimeZone, _ = time.LoadLocation("Greenwich")

	TwitterOAuth2JumpFrontUrl = ""
)
