package models

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/project-miko/miko/tools"
)

const (
	Enable2FAGA  = 1
	Disable2FAGA = 2
)

type Admin struct {
	Id          int64  `json:"id"`
	Username    string `json:"username"`
	Userpwd     string `json:"userpwd"`
	Salt        string `json:"salt"`
	Enable2FAGA int    `json:"enable_2fa_ga" gorm:"column:enable_2fa_ga"`
	CreatedAt   int64  `json:"created_at"`
	UpdatedAt   int64  `json:"updated_at"`
}

func (*Admin) TableName() string {
	return "admin"
}

func (m *Admin) Save() error {
	return GetDbInst().Save(m).Error
}

func (m *Admin) Update() error {
	m.UpdatedAt = tools.GetMillisecond(time.Now())
	return m.Save()
}

func FindAdminByUsername(username string) (*Admin, error) {
	admin := new(Admin)

	err := GetDbInst().Where("username=?", username).Find(admin).Error

	if gorm.IsRecordNotFoundError(err) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return admin, nil
}

func GetAdminById(id int64) (*Admin, error) {
	admin := new(Admin)
	err := GetDbInst().Where("id=?", id).Find(admin).Error
	if gorm.IsRecordNotFoundError(err) {
		return nil, nil
	}

	return admin, err
}
