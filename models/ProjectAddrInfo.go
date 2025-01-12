package models

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/project-miko/miko/tools"
)

type ProjectAddrInfo struct {
	Id         uint64 `json:"id"`
	Name       string `json:"name"`
	TestUrl    string `json:"test_url"`
	PreviewUrl string `json:"preview_url"`
	MainUrl    string `json:"main_url"`
	CreatedAt  int64  `json:"created_at"`
	UpdatedAt  int64  `json:"updated_at"`
}

func (m *ProjectAddrInfo) TableName() string {
	return "project_addr_info"
}

func (m *ProjectAddrInfo) Save() error {
	return GetDbInst().Save(m).Error
}

func (m *ProjectAddrInfo) Update() error {
	m.UpdatedAt = tools.GetMillisecond(time.Now())
	return m.Save()
}

func (m *ProjectAddrInfo) Del() error {
	return GetDbInst().Delete(m).Error
}

func GetProjectAddrInfoByName(name string) (*ProjectAddrInfo, error) {
	result := new(ProjectAddrInfo)
	err := GetDbInst().Where("name=?", name).Find(result).Error
	if gorm.IsRecordNotFoundError(err) {
		return nil, nil
	}

	return result, err
}

func GetProjectAddrInfoById(id int64) (*ProjectAddrInfo, error) {
	result := new(ProjectAddrInfo)
	err := GetDbInst().Where("id=?", id).Find(result).Error
	if gorm.IsRecordNotFoundError(err) {
		return nil, nil
	}

	return result, err
}

func GetProjectAddrInfoList(name string, page int64, limit int64) (int64, []*ProjectAddrInfo, error) {
	var amount int64
	results := make([]*ProjectAddrInfo, 0)
	db := GetDbInst()
	if len(name) > 0 {
		db = db.Where("name like ?", "%"+name+"%")
	}

	err := db.Model(ProjectAddrInfo{}).Count(&amount).Error
	if err != nil {
		return 0, nil, err
	}
	if amount == 0 {
		return 0, results, nil
	}

	offset := (page - 1) * limit
	err = db.Offset(offset).Limit(limit).Order("id").Find(&results).Error
	if gorm.IsRecordNotFoundError(err) {
		return 0, results, nil
	}

	return amount, results, err
}
