package models

import (
	"time"

	"github.com/project-miko/miko/tools"
)

type AdminLog struct {
	Id        int64
	AdminId   int64
	Content   string
	CreatedAt int64
}

func (*AdminLog) TableName() string {
	return "admin_log"
}

func InsertAdminLog(adminId int64, content string) error {
	adminLog := &AdminLog{
		AdminId:   adminId,
		Content:   content,
		CreatedAt: tools.GetMillisecond(time.Now()),
	}

	return GetDbInst().Save(adminLog).Error
}
