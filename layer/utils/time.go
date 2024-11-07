package utils

import "time"

const STR_MONTH_FORMAT = "0001/01"

func GetMonth(dateTime time.Time) string {
	return dateTime.Format(STR_MONTH_FORMAT)
}

func ParseMonthTime(strMonth string) (time.Time, error) {
	parsedTime, err := time.Parse(STR_MONTH_FORMAT, strMonth)
	if err != nil {
		return time.Time{}, err
	}
	return parsedTime, nil
}
