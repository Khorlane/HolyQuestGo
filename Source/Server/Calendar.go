package server

import (
	"os"
	"time"
)

var CalendarFile      *os.File
var CalendarFileIsOpen bool

type Calendar struct {
	Year              int
	Month             int
	Day               int
	Hour              int
	DayOfWeek         int
	TimeToAdvanceHour int64

	DayNames          []string
	DayOfMonth        []string
	HourNames         []string
	MonthNames        []string
}

var pCalendar =	 &Calendar{
	Year:              1,
	Month:             1,
	Day:               1,
	Hour:              1,
	DayOfWeek:         1,
	TimeToAdvanceHour: 60,
}

func CreateCalendar() *Calendar {
	DEBUGIT(1)
	timer := time.Now().Unix()
	pCalendar.TimeToAdvanceHour = timer + REAL_MINUTES_PER_HOUR*60
	OpenCalendarFile()
	if CalendarFileIsOpen {
		GetStartTime()
	}
	LoadDayNamesArray()
	LoadDayOfMonthArray()
	LoadHourNamesArray()
	LoadMonthNamesArray()
	return pCalendar
}

// Advances the in-game time by one hour
func AdvanceTime() {
}

// Closes the calendar file
func CloseCalendarFile() {
}

// Retrieves the start time from the calendar file
func GetStartTime() {
}

// Gets the current in-game date and time
func GetTime() string {
	return ""
}

// Loads the day names from a file into the DayNames array
func LoadDayNamesArray() {
}

// Loads the day of month names from a file into the DayOfMonth array
func LoadDayOfMonthArray() {
}

// Loads the hour names from a file into the HourNames array
func LoadHourNamesArray() {
}

// Loads the month names from a file into the MonthNames array
func LoadMonthNamesArray() {
}

// Opens the calendar file for reading
func OpenCalendarFile() {
	DEBUGIT(1)
	CalendarFileIsOpen = false
	CalendarFileName := CONTROL_DIR + "Calendar.txt"
	CalendarFileInp, err := os.Open(CalendarFileName)
	if err != nil {
		LogBuf = "Calendar file not found."
		LogIt(LogBuf)
		LogBuf = "Forcing start date to Year: 1 Month: 1 Day: 1 Hour: 1 Day of Week: 1"
		LogIt(LogBuf)
		return
	}
	CalendarFile = CalendarFileInp
	CalendarFileIsOpen = true
}

// Saves the current in-game time to the calendar file
func SaveTime() {
}
