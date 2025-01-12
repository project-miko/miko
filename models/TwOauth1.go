package models

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/project-miko/miko/tools"
)

type TwOAuth1 struct {
	Id           int64  `json:"id"`
	UserId       string `json:"user_id"`
	Name         string `json:"name"`
	Account      string `json:"account"`
	AccessToken  string `json:"access_token"`
	AccessSecret string `json:"access_secret"`
	CreatedAt    int64  `json:"created_at"`
	UpdatedAt    int64  `json:"updated_at"`
}

func (r *TwOAuth1) TableName() string {
	return "tw_oauth1"
}

func (r *TwOAuth1) Save() error {
	return GetDbInst().Save(r).Error
}

func (r *TwOAuth1) Update() error {
	r.UpdatedAt = tools.GetMillisecond(time.Now())
	return r.Save()
}

func (r *TwOAuth1) Del() error {
	return GetDbInst().Delete(r).Error
}

func GetTwOAuth1ByUserId(userId string) (*TwOAuth1, error) {
	result := new(TwOAuth1)
	err := GetDbInst().Where("user_id=?", userId).Find(result).Error
	if gorm.IsRecordNotFoundError(err) {
		return nil, nil
	}
	return result, err
}

func GetAllTwOAuth1List() ([]*TwOAuth1, error) {
	results := make([]*TwOAuth1, 0)
	err := GetDbInst().Find(&results).Error
	if gorm.IsRecordNotFoundError(err) {
		return results, nil
	}

	return results, err
}
