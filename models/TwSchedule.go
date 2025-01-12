package models

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/project-miko/miko/tools"
	"github.com/project-miko/miko/tools/strutils"
)

const (
	TwScheduleStatusUnFinished = 1
	TwScheduleStatusFinished   = 2
	TwScheduleStatusDeleted    = 3
	TwScheduleStatusError      = 4
)

type TwSchedule struct {
	Id              int64  `json:"id"`
	UserId          string `json:"user_id"`
	TwScheduleLibId int64  `json:"tw_schedule_lib_id"`
	CronExpression  string `json:"cron_expression"`
	TotalCount      int    `json:"total_count"`
	RemainCount     int    `json:"remain_count"`
	SourceType      int    `json:"source_type"`
	Status          int    `json:"status"`
	NextRunAt       int64  `json:"next_run_at"`
	CreatedAt       int64  `json:"created_at"`
	UpdatedAt       int64  `json:"updated_at"`
}

func (ts *TwSchedule) TableName() string {
	return "tw_schedule"
}

func (ts *TwSchedule) Save() error {
	return GetDbInst().Save(ts).Error
}

func (ts *TwSchedule) Update() error {
	ts.UpdatedAt = tools.GetMillisecond(time.Now())
	return ts.Save()
}

func (ts *TwSchedule) Del() error {
	return GetDbInst().Delete(ts).Error
}

func GetTwScheduleById(id int64, status int) (*TwSchedule, error) {
	result := new(TwSchedule)
	db := GetDbInst()
	db = db.Where("id=?", id)
	if status > -1 {
		db = db.Where("status=?", status)
	}

	err := db.Find(result).Error
	if gorm.IsRecordNotFoundError(err) {
		return nil, nil
	}
	return result, err
}

func GetTwScheduleByUserIdAndTwLibId(userId string, twScheduleLibId int64) (*TwSchedule, error) {
	result := new(TwSchedule)
	db := GetDbInst()
	db = db.Where("user_id=?", userId)
	db = db.Where("tw_schedule_lib_id=?", twScheduleLibId)
	db = db.Where("status=?", TwScheduleStatusUnFinished)

	err := db.Find(result).Error
	if gorm.IsRecordNotFoundError(err) {
		return nil, nil
	}
	return result, err
}

func GetAllTwScheduleList(status int) ([]*TwSchedule, error) {
	results := make([]*TwSchedule, 0)
	db := GetDbInst()
	if status > 0 {
		db = db.Where("status = ?", status)
	}

	err := db.Find(&results).Error
	if gorm.IsRecordNotFoundError(err) {
		return results, nil
	}
	return results, err
}

func GetTwScheduleList(userId string, status []int64, sourceType int, twScheduleLibIdList []int64, page, limit int64) (int64, []*TwSchedule, error) {
	var amount int64 = 0
	results := make([]*TwSchedule, 0)

	db := GetDbInst()
	if len(userId) > 0 {
		db = db.Where("user_id = ?", userId)
	}
	if len(status) != 0 {
		db = db.Where(fmt.Sprintf("status in (%s)", strutils.IdsToInString(status)))
	}
	if sourceType > 0 {
		db = db.Where("source_type = ?", sourceType)
	}
	if len(twScheduleLibIdList) > 0 {
		db = db.Where(fmt.Sprintf("tw_schedule_lib_id in (%s)", strutils.IdsToInString(twScheduleLibIdList)))
	}

	err := db.Model(TwSchedule{}).Count(&amount).Error
	if err != nil {
		return 0, nil, err
	}
	if amount == 0 {
		return 0, results, nil
	}

	orderStr := "created_at desc"
	if len(status) != 0 {
		for _, v := range status {
			if v == TwScheduleStatusUnFinished {
				orderStr = "next_run_at asc"
				break
			}
		}
	}

	offset := (page - 1) * limit
	err = db.Offset(offset).Limit(limit).Order(orderStr).Find(&results).Error
	if gorm.IsRecordNotFoundError(err) {
		return 0, results, nil
	}

	return amount, results, err
}
