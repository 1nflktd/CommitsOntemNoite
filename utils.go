package main

import (
	"time"
)

func getActualDate() time.Time {
	date := time.Now().UTC() // get day before
	local := date
	location, err := time.LoadLocation("America/Sao_Paulo")
	if err == nil {
		local = local.In(location)
	}
	return date
}

func getDate(format string, y, m, d int) string {
	return getActualDate().AddDate(y, m, d).Format(format)
}

func getDateTime(y, m, d int) time.Time {
	return getActualDate().AddDate(y, m, d)
}

func getOnlyDate(y, m, d int) time.Time {
	date := getActualDate().AddDate(y, m, d)
	year, month, day := date.Date()
	location, err := time.LoadLocation("America/Sao_Paulo")
	if err != nil {
		return date
	}
	return time.Date(year, month, day, 0, 0, 0, 0, location)
}
