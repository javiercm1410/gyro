package utils

import (
	"time"

	"github.com/charmbracelet/log"
)

var timeNow = time.Now

func DateSince(daysSince int) time.Time {
	timezone, err := time.LoadLocation("America/Santo_Domingo")
	if err != nil {
		log.Errorf("Failed to load time zone: %v", err)
		return time.Time{}
	}

	currentTime := timeNow().In(timezone)

	if daysSince <= 0 {
		return time.Time{}
	}

	lastDate := currentTime.AddDate(0, 0, -daysSince)
	log.Infof("Last login reference date: %v", lastDate.Format("2006-01-02 15:04:05"))

	return lastDate
}
