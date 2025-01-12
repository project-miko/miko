package models

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/project-miko/miko/tools"
)

const (
	UploadTypeLLM = 1
)

type UploadFileLog struct {
	Id         int64  `json:"id"`
	UploadUser string `json:"upload_user"`
	FilePath   string `json:"file_path"`
	FileHash   string `json:"file_hash"`
	Type       int    `json:"type"`
	CreatedAt  int64  `json:"created_at"`
	UpdatedAt  int64  `json:"updated_at"`
}

func (u *UploadFileLog) TableName() string {
	return "upload_file_log"
}

func (u *UploadFileLog) Save() error {
	return GetDbInst().Save(u).Error
}

func (u *UploadFileLog) Update() error {
	u.UpdatedAt = tools.GetMillisecond(time.Now())
	return u.Save()
}

func GetUploadFileLogList(fileName string, createdAt int64, page int64, limit int64) (int64, []*UploadFileLog, error) {
	var amount int64
	results := make([]*UploadFileLog, 0)
	db := GetDbInst()
	if len(fileName) > 0 {
		db = db.Where("file_path like ?", "%"+fileName+"%")
	}
	if createdAt > 0 {
		db = db.Where("created_at >= ?", createdAt)
	}

	err := db.Model(UploadFileLog{}).Count(&amount).Error
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

func GetUploadFileLogByHash(hash string) (*UploadFileLog, error) {
	result := new(UploadFileLog)
	err := GetDbInst().Where("file_hash=?", hash).Find(result).Error
	if gorm.IsRecordNotFoundError(err) {
		return nil, nil
	}
	return result, err
}
