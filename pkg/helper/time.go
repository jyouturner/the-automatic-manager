package helper

import (
	"fmt"
	"time"
)

//AddBusinessDays add days to the given time
func AddBusinessDays(date time.Time, noOfDays int) time.Time {
	actualDateAfterNoOfDays := date.Add(time.Duration(noOfDays) * 24 * time.Hour)
	noOfWorkingDaysInBetween := GetWorkingDaysInBetween(date, actualDateAfterNoOfDays)

	additionalDaysNeeded := noOfDays - noOfWorkingDaysInBetween

	fmt.Println(date, ",", actualDateAfterNoOfDays, ",", noOfWorkingDaysInBetween)
	if additionalDaysNeeded <= 0 {
		return actualDateAfterNoOfDays
	}

	return actualDateAfterNoOfDays.Add(
		AddBusinessDays(actualDateAfterNoOfDays, additionalDaysNeeded).Sub(actualDateAfterNoOfDays),
	)
}

// GetWorkingDaysInBetween find the business days between two dates
func GetWorkingDaysInBetween(s, e time.Time) int {
	sW := int(s.Weekday())
	n := GetDaysInBetween(s, e)
	saturdays := DayOfWeekCount(sW, 6, n)
	sundays := DayOfWeekCount(sW, 0, n)

	// formula excludes end date thus must add one if the date is a weekday
	additional := 1
	if e.Weekday() == time.Saturday || e.Weekday() == time.Sunday {
		additional = 0
	}
	return n - (saturdays + sundays) + additional
}

func GetDaysInBetween(s, e time.Time) int {
	return int(e.Sub(s).Hours() / 24.0)
}

func DayOfWeekCount(startDayOfWeek, targetDayOfWeek, noOfDays int) int {
	// Gives number of a weekday given; a starting weekday and number of days
	// startDayOfWeek, targetDayOfWeek are indexed 0 ([0..6] => [Sunday ... Saturday])
	x := CircularSequenceDistance(startDayOfWeek, targetDayOfWeek, 7)
	y := int(float64(noOfDays) / 7)
	z := 0
	z1 := CircularSequenceDistance(startDayOfWeek, ((startDayOfWeek+noOfDays)%7)-1, 7)
	if z1 >= x {
		z = 1
	}

	return y + z
}

// CircularSequenceDistance
func CircularSequenceDistance(from, to, length int) int {
	return ((to - from) + length) % length
}
