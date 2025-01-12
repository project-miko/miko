package models

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/project-miko/miko/tools"
	"github.com/project-miko/miko/tools/strutils"
)

type TwScheduleLib struct {
	Id        int64  `json:"id"`
	Content   string `json:"content"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}

func (tsl *TwScheduleLib) TableName() string {
	return "tw_schedule_lib"
}

func (tsl *TwScheduleLib) Save() error {
	return GetDbInst().Save(tsl).Error
}

func (tsl *TwScheduleLib) Update() error {
	tsl.UpdatedAt = tools.GetMillisecond(time.Now())
	return tsl.Save()
}

func (tsl *TwScheduleLib) Del() error {
	return GetDbInst().Delete(tsl).Error
}

func GetTwScheduleLibById(id int64) (*TwScheduleLib, error) {
	result := new(TwScheduleLib)
	err := GetDbInst().Where("id=?", id).Find(result).Error
	if gorm.IsRecordNotFoundError(err) {
		return nil, nil
	}
	return result, err
}

func GetTwScheduleListByIds(ids []int64) ([]*TwScheduleLib, error) {
	results := make([]*TwScheduleLib, 0)
	db := GetDbInst()
	db = db.Where(fmt.Sprintf("id in (%s)", strutils.IdsToInString(ids)))

	err := db.Find(&results).Error
	if gorm.IsRecordNotFoundError(err) {
		return results, nil
	}
	return results, err
}
