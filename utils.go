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