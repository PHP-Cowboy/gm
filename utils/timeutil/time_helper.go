package timeutil

import (
	"fmt"
	"math"
	"sort"
	"time"
)

func GetDateTime() string {
	return time.Now().Format(TimeFormat)
}

func GetDate() string {
	return time.Now().Format(DateNumberFormat)
}

// 获取传入的时间所在月份的第一天，即某月第一天的0点。如传入time.Now(), 返回当前月份的第一天0点时间。
func GetFirstDateOfMonth(d time.Time) time.Time {
	d = d.AddDate(0, 0, -d.Day()+1)
	return GetZeroTime(d)
}

// 获取传入的时间所在月份的最后一天，即某月最后一天的0点。如传入time.Now(), 返回当前月份的最后一天0点时间。
func GetLastDateOfMonth(d time.Time) time.Time {
	return GetFirstDateOfMonth(d).AddDate(0, 1, -1)
}

// 获取某一天的0点时间
func GetZeroTime(d time.Time) time.Time {
	return time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, d.Location())
}

// 获取某一天的最后一秒时间
func GetLastTime(d time.Time) time.Time {
	return time.Date(d.Year(), d.Month(), d.Day(), 23, 59, 59, 0, d.Location())
}

// 格式化为一天最小时间格式
func FormatTimeToMinDateTime(d time.Time, format string) string {
	d = GetZeroTime(d)
	return d.Format(format)
}

// 格式化为一天最大时间格式
func FormatTimeToMaxDateTime(d time.Time, format string) string {
	d = GetLastTime(d)
	return d.Format(format)
}

// 时间字符串转换为时间戳
// timeStr 时间字符串
// timeStringFormat 时间字符串的格式 不传递默认使用 2006-01-02 15:04:05
func StandardStr2Time(timeStr string, timeStringFormat ...string) int64 {
	timeFormat := TimeFormat
	if len(timeStringFormat) > 0 {
		timeFormat = timeStringFormat[0]
	}
	t, err := time.Parse(timeFormat, timeStr)
	if err != nil {
		return int64(0)
	}
	return t.Unix()
}

// 获取当前时间还有多少秒到下一天
func DayLeftSeconds() int64 {
	now := time.Now()
	return now.Unix() - GetLastTime(now).Unix()
}

// 日期格式的字符串转当天 00:00:00 的日期时间格式字符串
// 例如 2020-10-26 转 2020-10-26 00:00:00
// dateTimeSting 时间
// formFormat 时间格式
// toFormat 目标格式
func GetDateTimeStartByDate(dateTimeString string, formFormat, toFormat string) (string, error) {
	t, err := time.Parse(formFormat, dateTimeString)
	if err != nil {
		return "", err
	}
	return FormatTimeToMinDateTime(t, toFormat), nil
}

func GetDateTimeByDateTimeString(dateTimeString string, formFormat, toFormat string) (string, error) {
	t, err := time.Parse(formFormat, dateTimeString)
	if err != nil {
		return "", err
	}
	return t.Format(toFormat), nil
}

// 日期格式的字符串转当天 23:59:59 的日期时间格式字符串
// 例如 2020-10-26 转 2020-10-26 23:59:59
// dateTimeSting 时间
// formFormat 时间格式
// toFormat 目标格式
func GetDateTimeEndByDate(dateTimeString, formFormat, toFormat string) (string, error) {
	t, err := time.Parse(formFormat, dateTimeString)
	if err != nil {
		return "", err
	}
	return FormatTimeToMaxDateTime(t, toFormat), nil
}

// 计算两个日期之间差多少天
func DiffDays(startDateTime, endDateTime time.Time) int64 {
	return int64(math.Abs(endDateTime.Sub(startDateTime).Hours() / 24))
}

// 获取当前日期前后 days 天 months月 years年 日期
func GetTimeAroundByNum(days, months, years int) string {
	t := time.Now().AddDate(years, months, days)
	return GetZeroTime(t).Format(TimeFormat)
}

// 获取当前日期前后 days 天 日期
func GetTimeAroundByDays(days int) string {
	t := time.Now().AddDate(0, 0, days)
	return GetZeroTime(t).Format(TimeFormat)
}

// 获取当前日期前后 days 天 日期
func GetTimeAroundByMonths(months int) string {
	t := time.Now().AddDate(0, months, 0)
	return GetZeroTime(t).Format(TimeFormat)
}

func GetCurrentMonth() string {
	return time.Now().Format("01")
}

// 获取当前时间前后时间段时间
func GetDurationDateTime(s string) (string, error) {

	duration, err := time.ParseDuration(s) //"-30m"

	if err != nil {
		return "", err
	}

	dateTime := time.Now().Add(duration).Format("2006-01-02 15:04:05")

	return dateTime, nil
}

func DateRangeToZeroAndLastTime(start, end string) (startTime, endTime time.Time, err error) {
	// 尝试按完整时间格式解析开始时间和结束时间
	startTime, err = time.ParseInLocation(TimeFormat, start, time.Local)
	if err != nil {
		// 如果按完整时间格式解析失败，则尝试仅按日期格式解析
		startTime, err = time.ParseInLocation(DateFormat, start, time.Local)
		if err != nil {
			return time.Time{}, time.Time{}, fmt.Errorf("failed to parse start time: %w", err)
		}
		// 如果只解析了日期，则将其转换为0点时间
		startTime = GetZeroTime(startTime)
	}

	endTime, err = time.ParseInLocation(TimeFormat, end, time.Local)

	if err != nil {
		// 如果按完整时间格式解析失败，则尝试仅按日期格式解析
		endTime, err = time.ParseInLocation(DateFormat, end, time.Local)
		if err != nil {
			return time.Time{}, time.Time{}, fmt.Errorf("failed to parse end time: %w", err)
		}
		// 如果只解析了日期，则将其转换为最后一秒时间
		endTime = GetLastTime(endTime)
	}

	return
}

func DateRangeToZeroAndLastTimeFormat(start, end string, format string) (startFormat, endFormat string, err error) {
	var startTime, endTime time.Time

	// 尝试按完整时间格式解析开始时间和结束时间
	startTime, err = time.ParseInLocation(format, start, time.Local)
	if err != nil {
		return
	}

	endTime, err = time.ParseInLocation(format, end, time.Local)

	if err != nil {
		return
	}

	// 如果只解析了日期，则将其转换为0点时间
	startFormat = GetZeroTime(startTime).Format(format)
	// 如果只解析了日期，则将其转换为最后一秒时间
	endFormat = GetLastTime(endTime).Format(format)

	return
}

// getMonthsInRange 计算开始日期到结束日期之间每个月的第一天
func GetMonthsInRange(startDateStr, endDateStr, format string) ([]time.Time, error) {
	var months []time.Time

	// 解析开始日期
	startDate, err := time.Parse(format, startDateStr)
	if err != nil {
		return nil, err
	}

	//解析结束日期
	endDate, err := time.Parse(format, endDateStr)
	if err != nil {
		return nil, err
	}

	// 确保开始日期是月初的第一天
	startDate = time.Date(startDate.Year(), startDate.Month(), 1, 0, 0, 0, 0, startDate.Location())

	// 循环直到endDate所在的月份（包括它）
	currentMonth := startDate
	for !currentMonth.After(endDate) {
		months = append(months, currentMonth)

		// 计算下一个月的第一天
		nextMonth := currentMonth.AddDate(0, 1, 0)
		// 如果下一个月的第一天超过了endDate，则退出循环
		if nextMonth.Month() != endDate.Month() || nextMonth.Year() != endDate.Year() {
			break
		}
		currentMonth = nextMonth
	}

	// 使用sort.SliceStable进行降序排序
	sort.SliceStable(months, func(i, j int) bool {
		// 降序排序，所以返回months[i].After(months[j])
		return months[i].After(months[j])
	})

	return months, nil
}
