package models

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/project-miko/miko/tools"
)

const (
	// 1 internal account 2 third-party account
	Internal   = 1
	ThirdParty = 2
)

var (
	AllowUserType = map[int64]interface{}{
		Internal:   nil,
		ThirdParty: nil,
	}
)

type TwUserInfo struct {
	Id              int64  `json:"id"`
	UserId          string `json:"user_id"`
	Name            string `json:"name"`
	Account         string `json:"account"`
	Type            int    `json:"type"`
	ProfileImageUrl string `json:"profile_image_url"`
	CreatedAt       int64  `json:"created_at"`
	UpdatedAt       int64  `json:"updated_at"`
}

func (info *TwUserInfo) TableName() string {
	return "tw_user_info"
}

func (info *TwUserInfo) Save() error {
	return GetDbInst().Save(info).Error
}

func (info *TwUserInfo) Update() error {
	info.UpdatedAt = tools.GetMillisecond(time.Now())
	return info.Save()
}

func (info *TwUserInfo) Del() error {
	return GetDbInst().Delete(info).Error
}

func GetUserInfoById(id int64) (*TwUserInfo, error) {
	userInfo := new(TwUserInfo)
	err := GetDbInst().Where("id=?", id).Find(userInfo).Error
	if gorm.IsRecordNotFoundError(err) {
		return nil, nil
	}
	return userInfo, err
}

func GetUserInfoByUserId(id string) (*TwUserInfo, error) {
	userInfo := new(TwUserInfo)
	err := GetDbInst().Where("user_id=?", id).Find(userInfo).Error
	if gorm.IsRecordNotFoundError(err) {
		return nil, nil
	}
	return userInfo, err
}

func GetTwUserInfoCount() (int64, error) {
	var amount int64 = 0
	err := GetDbInst().Model(TwUserInfo{}).Count(&amount).Error
	if err != nil {
		return 0, err
	}
	return amount, nil
}

func GetTwUserInfoListByType(page, limit, userType, start int64) (int64, []*TwUserInfo, error) {
	var amount int64
	results := make([]*TwUserInfo, 0)
	db := GetDbInst()
	db = db.Where("type=?", userType)
	//db = db.Where("created_at<=?", start)

	err := db.Model(TwUserInfo{}).Count(&amount).Error
	if err != nil {
		return 0, nil, err
	}
	if amount == 0 {
		return 0, results, nil
	}

	offset := (page - 1) * limit
	err = db.Offset(offset).Limit(limit).Find(&results).Error
	if gorm.IsRecordNotFoundError(err) {
		return 0, results, nil
	}

	return amount, results, err
}

func GetTwUserInfoList(page, limit int64) ([]*TwUserInfo, error) {
	results := make([]*TwUserInfo, 0)
	offset := page * limit
	err := GetDbInst().Offset(offset).Limit(limit).Find(&results).Error
	if gorm.IsRecordNotFoundError(err) {
		return results, nil
	}

	return results, err
}

func GetAllTwUserInfoList(userType int64) ([]*TwUserInfo, error) {
	results := make([]*TwUserInfo, 0)
	db := GetDbInst()
	if userType > 0 {
		db = db.Where("type=?", userType)
	}
	err := db.Find(&results).Error
	if gorm.IsRecordNotFoundError(err) {
		return results, nil
	}

	return results, err
}
