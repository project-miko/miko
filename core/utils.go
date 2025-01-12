package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"text/template"
	"time"

	"github.com/g8rswimmer/go-twitter/v2"
	"github.com/gomodule/redigo/redis"
	"github.com/project-miko/miko/conf"
	"github.com/project-miko/miko/models"
	"github.com/project-miko/miko/tools/log"
	"github.com/shopspring/decimal"
)

const (
	red   = "red"
	green = "green"
	down  = "↓"
	up    = "↑"
)

func NewDefaultRateLimit() *twitter.RateLimit {
	limit := new(twitter.RateLimit)
	limit.Limit = 100
	limit.Remaining = 100
	limit.Reset = twitter.Epoch(time.Now().Unix())

	return limit
}

func SaveUserRateLimitCreateTweet(userId string, resp *twitter.CreateTweetResponse, dstErr error) {
	var err error
	var limit *twitter.RateLimit

	if resp != nil && resp.RateLimit != nil {
		limit = resp.RateLimit
	} else if dstErr != nil {
		limit, _ = twitter.RateLimitFromError(dstErr)
	}

	if limit != nil {
		if err = SaveUserRateLimit2Redis(userId, limit); err != nil {
			log.Error("", "Error saving rate limit to Redis: %v", err)
		}
	}
}

func GetUserRateLimit(userId string) (*twitter.RateLimit, error) {
	client := models.GetRdbInst()

	key := fmt.Sprintf(conf.AISERAppRateLimitCreateTweet, userId)
	jsonStr, err := client.GetString(key)
	if err != nil {
		if err == redis.ErrNil {
			return NewDefaultRateLimit(), nil
		}

		return nil, err
	}

	limit := new(twitter.RateLimit)
	err = json.Unmarshal([]byte(jsonStr), limit)
	if err != nil {
		return nil, err
	}

	if time.Now().Unix() > int64(limit.Reset) {
		return NewDefaultRateLimit(), nil
	}

	return limit, nil
}

func SaveUserRateLimit2Redis(userId string, limit *twitter.RateLimit) error {
	client := models.GetRdbInst()
	b, err := json.Marshal(limit)
	if err != nil {
		return err
	}
	key := fmt.Sprintf(conf.AISERAppRateLimitCreateTweet, userId)
	err = client.SetString(key, string(b), int64(1*24*time.Hour.Seconds()))
	if err != nil {
		return err
	}

	return nil
}

func WrapColorForInteger(dst int) string {
	color := ""
	flag := ""

	if dst > 0 {
		color = green
		flag = up
	} else if dst < 0 {
		flag = down
		color = red
	}

	return fmt.Sprintf(`<font color='%s'>%s%d</font>`, color, flag, dst)
}

func WrapColorForDecimal(dst decimal.Decimal) string {
	color := ""
	flag := ""

	if dst.Cmp(decimal.Zero) > 0 {
		color = green
		flag = up
	} else if dst.Cmp(decimal.Zero) < 0 {
		flag = down
		color = red
	}

	dst = dst.Round(1)

	return fmt.Sprintf("<font color='%s'>%s%s</font>", color, flag, dst.String()+"%")
}

func ParseTemplate2String(fileName string, obj interface{}) (string, error) {
	content, err := ReadFileContent(fileName)
	if err != nil {
		return "", err
	}

	// parse template
	tmpl, err := template.New("template").Parse(content)
	if err != nil {
		return "", err
	}

	buf := bytes.NewBuffer(nil)
	// execute template and write result to string variable
	err = tmpl.Execute(buf, obj)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func ReadFileContent(filename string) (string, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

type ResponseDecodeError struct {
	Name string
	Err  error
}

func (r *ResponseDecodeError) Error() string {
	return fmt.Sprintf("%s decode error: %v", r.Name, r.Err)
}

type HTTPError struct {
	Status     string `json:"status"`
	StatusCode int    `json:"status_code"`
	URL        string `json:"url"`
}

func (h HTTPError) Error() string {
	return fmt.Sprintf("code: %d status: %s", h.StatusCode, h.Status)
}
