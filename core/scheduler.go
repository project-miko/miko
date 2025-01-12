package core

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/project-miko/miko/conf"
)

var scheduler *Scheduler

var (
	ErrScheduleJobExists = fmt.Errorf("schedule job exists")
)

type Scheduler struct {
	Scheduler *gocron.Scheduler
	Tag       string
	Event     interface{}
	Params    []interface{}
	JobFun    interface{}
}

func InitScheduler() {
	s := gocron.NewScheduler(conf.NewTimeZone)
	s.WaitForScheduleAll()
	s.TagsUnique()
	s.StartAsync()

	scheduler = new(Scheduler)
	scheduler.Scheduler = s

	if err := InitTwCreateTweetJobs(); err != nil {
		panic(err)
	}
}

func (s *Scheduler) SetJobFuncAndParams(jobFun interface{}, params ...interface{}) {
	s.JobFun = jobFun
	s.Params = params
}

func GetScheduler() *Scheduler {
	return scheduler
}

func (s *Scheduler) Add(cronExp string, tag string, limit int, nextRunAt int64) (*gocron.Job, error) {
	err := s.checkJobExists(tag)
	if err != nil {
		return nil, err
	}

	var j *gocron.Job
	if strings.Contains(cronExp, "/") {
		j, err = s.AddWithTime(cronExp, tag, limit, nextRunAt)
	} else {
		j, err = s.Scheduler.CronWithSeconds(cronExp).Tag(tag).LimitRunsTo(limit).Do(s.JobFun, s.Params...)
	}

	if err != nil {
		return nil, err
	}

	return j, nil
}

func (s *Scheduler) AddWithTime(cronExp string, tag string, limit int, nextRunAt int64) (*gocron.Job, error) {
	params, err := ParseCronExpInterval(cronExp)
	if err != nil {
		return nil, err
	}

	_nextRunAt := time.UnixMilli(nextRunAt).In(conf.NewTimeZone)
	if nextRunAt == -1 {
		now := time.Now().In(conf.NewTimeZone)
		weekDay := params.WeekDay
		if weekDay == 0 {
			weekDay = 7
		}

		nextRunDay := weekDay - int(now.Weekday())
		if isBeforeNow(params.WeekDay, params.Hour, params.Minute) {
			nextRunDay = nextRunDay + 7
		}

		_nextRunAt = time.Date(now.Year(), now.Month(), now.Day()+nextRunDay, params.Hour, params.Minute, 0, 0, now.Location())
	}

	interval := time.Duration(params.LoopUnit) * 7 * 24 * time.Hour
	j, err := s.Scheduler.StartAt(_nextRunAt).Every(interval).Tag(tag).LimitRunsTo(limit).Do(s.JobFun, s.Params...)
	if err != nil {
		return nil, err
	}

	return j, nil
}

func (s *Scheduler) checkJobExists(tag string) error {
	jobs, err := s.Scheduler.FindJobsByTag(tag)
	if err != nil {
		if err != gocron.ErrJobNotFoundWithTag {
			return err
		}
	}
	if len(jobs) != 0 {
		return ErrScheduleJobExists
	}

	return nil
}

func (s *Scheduler) Remove(tag string) error {
	err := s.Scheduler.RemoveByTag(tag)
	if err != nil {
		return err
	}

	return nil
}

type CronExpParam struct {
	Minute   int
	Hour     int
	WeekDay  int
	LoopUnit int
}

func ParseCronExpInterval(cronExp string) (*CronExpParam, error) {
	var err error

	fields := strings.Fields(cronExp)
	if len(fields) != 6 {
		return nil, fmt.Errorf("invalid cron expression")
	}

	minute := 0
	hour := 0
	weekday := 0
	weekInterval := 0

	if fields[1] != "*" {
		minute, err = strconv.Atoi(fields[1])
		if err != nil {
			return nil, err
		}
	}

	if fields[2] != "*" {
		hour, err = strconv.Atoi(fields[2])
		if err != nil {
			return nil, err
		}
	}

	weekdayField := fields[5]
	if strings.Contains(weekdayField, "/") {
		weekday, err = strconv.Atoi(weekdayField[:1])
		if err != nil {
			return nil, err
		}
		weekInterval, err = strconv.Atoi(weekdayField[2:])
		if err != nil {
			return nil, err
		}
		weekInterval = weekInterval / 7
	} else {
		weekday, err = strconv.Atoi(weekdayField)
		if err != nil {
			return nil, err
		}
	}

	result := &CronExpParam{
		Minute:   minute,
		Hour:     hour,
		WeekDay:  weekday,
		LoopUnit: weekInterval,
	}

	return result, nil
}

func isBeforeNow(weekDay, hour, minute int) bool {
	now := time.Now().In(conf.TimeZone)

	nowWeekDay, nowHour, nowMinute := int(now.Weekday()), now.Hour(), now.Minute()

	if weekDay == 0 {
		weekDay = 7
	}
	if nowWeekDay == 0 {
		nowWeekDay = 7
	}

	// first check week
	if weekDay < nowWeekDay {
		return true // given week is before current week
	} else if weekDay > nowWeekDay {
		return false
	}

	// if week is same, then check hour
	if hour < nowHour {
		return true // given hour is before current hour
	} else if hour > nowHour {
		return false
	}

	// if hour is same, then check minute
	if minute < nowMinute {
		return true // given minute is before current minute
	} else {
		return false
	}
}
