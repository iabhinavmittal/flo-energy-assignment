package main

import (
	"time"
)

func ConvertStringToDate(dateString string) (time.Time, error) {
	dateFormat := "20060102"
	date, err := time.Parse(dateFormat, dateString)

	if err != nil {
		return time.Now(), err
	}

	return date, nil
}
