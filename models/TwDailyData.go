package models

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/project-miko/miko/tools"
	"github.com/project-miko/miko/tools/strutils"
)

type TwDailyData struct {
	Id            int64  `json:"id"`
	UserId        string `json:"user_id"`
	TweetIds      string `json:"tweet_ids"`
	LikeCount     int    `json:"like_count"`
	ReplyCount    int    `json:"reply_count"`
	RetweetCount  int    `json:"retweet_count"`
	TweetCount    int    `json:"tweet_count"`
	FollowerCount int    `json:"follower_count"`
	StatisticAt   int64  `json:"statistic_at"`
	CreatedAt     int64  `json:"created_at"`
	UpdatedAt     int64  `json:"updated_at"`
}

func (d *TwDailyData) TableName() string {
	return "tw_daily_data"
}

func (d *TwDailyData) Save() error {
	return GetDbInst().Save(d).Error
}

func (d *TwDailyData) Update() error {
	d.UpdatedAt = tools.GetMillisecond(time.Now())
	return d.Save()
}

func (d *TwDailyData) Del() error {
	return GetDbInst().Delete(d).Error
}

func GetListByUserIdsAndTime(userIds []string, dstTime int64) ([]*TwDailyData, error) {
	results := make([]*TwDailyData, 0)
	db := GetDbInst()
	db = db.Where("statistic_at=?", dstTime)
	db = db.Where(fmt.Sprintf("user_id in (%s)", strutils.StringSliceToInString(userIds)))
	err := db.Find(&results).Error
	if gorm.IsRecordNotFoundError(err) {
		return results, nil
	}

	return results, err
}

func GetListByUserIds(userIds []string, start, end int64) ([]*TwDailyData, error) {
	results := make([]*TwDailyData, 0)
	db := GetDbInst()
	db = db.Where("statistic_at between ? and ?", start, end)
	db = db.Where(fmt.Sprintf("user_id in (%s)", strutils.StringSliceToInString(userIds)))
	err := db.Find(&results).Error
	if gorm.IsRecordNotFoundError(err) {
		return results, nil
	}

	return results, err
}
