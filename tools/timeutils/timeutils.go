package timeutils

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

var (
	location *time.Location
)

func init() {
	location = time.FixedZone("UTC+8", 8*60*60)
}

func TimeSectionOfMonth(t time.Time) (start, end int64, err error) {
	year, month, _, err := GetDateNumber(t)
	if err != nil {
		return 0, 0, err
	}

	// month
	monthBegin, err := time.ParseInLocation("2006-1-2", fmt.Sprintf("%d-%d-1", year, month), location)
	if err != nil {
		return 0, 0, err
	}

	monthEnd := monthBegin.AddDate(0, 1, 0)
	return monthBegin.Unix(), monthEnd.Unix(), nil
}

func TimeSectionOfWeek(t time.Time) (start, end int64, err error) {
	// week
	offset := int(time.Monday - t.Weekday())
	if offset > 0 {
		offset = -6
	}

	weekBegin := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, location).AddDate(0, 0, offset)
	weekEnd := weekBegin.Unix() + 604800
	return weekBegin.Unix(), weekEnd, nil
}

func TimeSectionOfDay(t time.Time) (start, end int64, err error) {

	year, month, day, err := GetDateNumber(t)
	if err != nil {
		return 0, 0, err
	}

	// today
	todayBegin, err := time.ParseInLocation("2006-1-2", fmt.Sprintf("%d-%d-%d", year, month, day), location)
	if err != nil {
		return 0, 0, err
	}

	todayEnd := todayBegin.AddDate(0, 0, 1)

	return todayBegin.Unix(), todayEnd.Unix(), nil
}

func GetDateNumber(t time.Time) (year, month, day int64, err error) {

	nowStr := t.Format("2006-1-2")

	nowArr := strings.Split(nowStr, "-")

	if len(nowArr) < 3 {
		return 0, 0, 0, fmt.Errorf("wrong time")
	}

	year, err = strconv.ParseInt(nowArr[0], 10, 64)
	if err != nil {
		return 0, 0, 0, err
	}

	month, err = strconv.ParseInt(nowArr[1], 10, 64)
	if err != nil {
		return 0, 0, 0, err
	}

	day, err = strconv.ParseInt(nowArr[2], 10, 64)
	if err != nil {
		return 0, 0, 0, err
	}

	return year, month, day, nil
}

func FormatShortData(target time.Time) (time.Time, error) {
	layout := "2006-1-2"
	timeStr := target.Format(layout)
	result, err := time.ParseInLocation(layout, timeStr, location)
	return result, err
}

func FormatShortDataTime(target time.Time, interval int) (time.Time, error) {
	// calculate the time difference between the previous n minutes
	previousOffset := target.Minute() % interval
	previous := target.Add(time.Duration(-previousOffset) * time.Minute)

	layout := "2006-1-2 15:04"
	timeStr := previous.Format(layout)
	result, err := time.ParseInLocation(layout, timeStr, location)
	return result, err
}
