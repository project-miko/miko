package models

import (
	"time"

	"github.com/project-miko/miko/tools"
)

const (
	ScheduleExecStatusSuccess = 1
	ScheduleExecStatusFail    = 2
)

type ScheduleLog struct {
	Id           int64  `json:"id"`
	JobId        string `json:"job_id"`
	ExecDuration int    `json:"exec_duration"`
	Status       int    `json:"status"`
	ErrorMsg     string `json:"error_msg"`
	ExecAt       int64  `json:"exec_at"`
	CreatedAt    int64  `json:"created_at"`
	UpdatedAt    int64  `json:"updated_at"`
}

func (r *ScheduleLog) TableName() string {
	return "schedule_logs"
}

func (r *ScheduleLog) Save() error {
	return GetDbInst().Save(r).Error
}

func (r *ScheduleLog) Update() error {
	r.UpdatedAt = tools.GetMillisecond(time.Now())
	return r.Save()
}

func (r *ScheduleLog) Del() error {
	return GetDbInst().Delete(r).Error
}

func SaveScheduleLog(jobId string, execAt int64, err error) error {
	reqLog := &ScheduleLog{
		JobId:  jobId,
		ExecAt: execAt,
		Status: ScheduleExecStatusSuccess,
	}

	if err != nil {
		reqLog.ErrorMsg = err.Error()
		reqLog.Status = ScheduleExecStatusFail
	}

	reqLog.ExecDuration = int(tools.GetMillisecond(time.Now()) - execAt)
	reqLog.CreatedAt = tools.GetMillisecond(time.Now())
	if e := reqLog.Save(); e != nil {
		return e
	}

	return err
}
