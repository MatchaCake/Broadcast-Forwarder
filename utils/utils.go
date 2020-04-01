package utils

import (
	"log"
	"time"
)

func ParseTimeWithLocation(t string, loc *time.Location) string {
	result, err := time.ParseInLocation(time.RFC3339, t, loc)
	if err != nil {
		log.Fatalf("error in parsing time to location:%v time:%v", loc.String(), err)
	}
	return result.String()
}

func ParseToBeijingTime(t string) string {
	secondsEastOfUTC := int((8 * time.Hour).Seconds())
	beijing := time.FixedZone("Beijing Time", secondsEastOfUTC)
	return ParseTimeWithLocation(t, beijing)
}