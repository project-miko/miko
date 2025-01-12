package models

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/project-miko/miko/tools"
)

type Admin2FAGA struct {
	Id                int64  `json:"id"`
	AdminId           int64  `json:"admin_id"`
	Secret            string `json:"secret"`
	RecoverCode       string `json:"recover_code"`
	RecoverCodeStatus int    `json:"recover_code_status"`
	CreatedAt         int64  `json:"created_at"`
	UpdatedAt         int64  `json:"updated_at"`
}

func (*Admin2FAGA) TableName() string {
	return "admin_2fa_ga"
}

func (m *Admin2FAGA) Save() error {
	return GetDbInst().Save(m).Error
}

func (m *Admin2FAGA) Update() error {
	m.UpdatedAt = tools.GetMillisecond(time.Now())
	return m.Save()
}

func GetAdmin2FAGAByAdminId(adminId int64) (*Admin2FAGA, error) {
	result := new(Admin2FAGA)
	err := GetDbInst().Where("admin_id=?", adminId).Find(result).Error
	if gorm.IsRecordNotFoundError(err) {
		return nil, nil
	}

	return result, err
}
