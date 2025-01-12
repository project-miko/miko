package models

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/project-miko/miko/tools"
)

const (
	TweetLikeType      = 1
	ReTweetType        = 2
	ReTweetAndLikeType = 3

	TweetLikeScope      = "like.write"
	ReTweetScope        = "tweet.write"
	ReTweetAndLikeScope = "like.write tweet.write"
)

var (
	AllowAuthType = map[int64]string{
		TweetLikeType:      TweetLikeScope,
		ReTweetType:        ReTweetScope,
		ReTweetAndLikeType: ReTweetAndLikeScope,
	}
)

type TwAccount struct {
	Id              int64  `json:"id"`
	UserId          string `json:"user_id"`
	Name            string `json:"name"`
	Account         string `json:"account"`
	ProfileImageUrl string `json:"profile_image_url"`
	AuthType        int    `json:"auth_type"`
	Scope           string `json:"scope"`
	AccessToken     string `json:"access_token"`
	RefreshToken    string `json:"refresh_token"`
	ExpiredAt       int64  `json:"expired_at"`
	CreatedAt       int64  `json:"created_at"`
	UpdatedAt       int64  `json:"updated_at"`
}

func (ta *TwAccount) TableName() string {
	return "tw_account"
}

func (ta *TwAccount) Save() error {
	return GetDbInst().Save(ta).Error
}

func (ta *TwAccount) Update() error {
	ta.UpdatedAt = tools.GetMillisecond(time.Now())
	return ta.Save()
}

func (ta *TwAccount) Del() error {
	return GetDbInst().Delete(ta).Error
}

func (ta *TwAccount) SaveOrUpdateWithLog() error {
	tx := GetDbInst().Begin()
	if tx.Error != nil {
		return tx.Error
	}

	defer func() {
		tx.Rollback()
	}()

	err := tx.Save(ta).Error
	if err != nil {
		return fmt.Errorf("save TwAccount error %s", err.Error())
	}

	tLog := &TwAccountLog{
		UserId:       ta.UserId,
		Name:         ta.Name,
		Account:      ta.Account,
		Scope:        ta.Scope,
		AccessToken:  ta.AccessToken,
		RefreshToken: ta.RefreshToken,
		ExpiredAt:    ta.ExpiredAt,
		CreatedAt:    tools.GetMillisecond(time.Now()),
	}

	err = tx.Save(tLog).Error
	if err != nil {
		return fmt.Errorf("save TwAccountLog error %s", err.Error())
	}

	tx.Commit()

	return nil
}

func GetReTweetScopeAccount() ([]*TwAccount, error) {
	results := make([]*TwAccount, 0)
	err := GetDbInst().Where("auth_type=? or auth_type=?", ReTweetType, ReTweetAndLikeType).Find(&results).Error
	if gorm.IsRecordNotFoundError(err) {
		return results, nil
	}
	return results, err
}

func GetLikeScopeAccount() ([]*TwAccount, error) {
	results := make([]*TwAccount, 0)
	err := GetDbInst().Where("auth_type=? or auth_type=?", TweetLikeType, ReTweetAndLikeType).Find(&results).Error
	if gorm.IsRecordNotFoundError(err) {
		return results, nil
	}
	return results, err
}

func GetTwAccountByUserId(userId string) (*TwAccount, error) {
	result := new(TwAccount)
	err := GetDbInst().Where("user_id=?", userId).Find(result).Error
	if gorm.IsRecordNotFoundError(err) {
		return nil, nil
	}
	return result, err
}

func GetAccountList(page, limit int64) (int64, []*TwAccount, error) {
	var amount int64 = 0
	results := make([]*TwAccount, 0)

	db := GetDbInst()
	err := db.Model(TwAccount{}).Count(&amount).Error
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

func GetAllAccountList(name, account string) ([]*TwAccount, error) {
	results := make([]*TwAccount, 0)
	db := GetDbInst()
	if len(name) > 0 {
		db = db.Where("name like ?", "%"+name+"%")
	}
	if len(account) > 0 {
		db = db.Where("account like ?", "%"+account+"%")
	}

	err := db.Find(&results).Error
	if gorm.IsRecordNotFoundError(err) {
		return results, nil
	}
	return results, err
}
