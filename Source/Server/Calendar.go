package server

import (
	"bufio"
	"fmt"
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

// Advance the in-game time by one hour if the real-world time has passed the threshold
func AdvanceTime() {
	var NowSec int64

	DEBUGIT(5)
	NowSec = time.Now().Unix()
	if NowSec < pCalendar.TimeToAdvanceHour {
		// Not time to advance the hour
		return
	}
	pCalendar.TimeToAdvanceHour = NowSec + REAL_MINUTES_PER_HOUR*60
	pCalendar.Hour++
	if pCalendar.Hour > HOURS_PER_DAY {
		pCalendar.Day++
		pCalendar.Hour = 1
		pCalendar.DayOfWeek++
		if pCalendar.DayOfWeek > DAYS_PER_WEEK {
			pCalendar.DayOfWeek = 1
		}
	}
	if pCalendar.Day > DAYS_PER_MONTH {
		pCalendar.Month++
		pCalendar.Day = 1
	}
	if pCalendar.Month > MONTHS_PER_YEAR {
		pCalendar.Year++
		pCalendar.Month = 1
	}
	SaveTime()
}

// CloseCalendarFile closes the calendar file.
func CloseCalendarFile() {
	DEBUGIT(1)
	CalendarFile.Close()
}

// Retrieve the start time from the calendar file
func GetStartTime() {
	DEBUGIT(1)
	Stuff := ""
	scanner := bufio.NewScanner(CalendarFile)
	if scanner.Scan() {
		Stuff = scanner.Text()
		LogBuf = "Read line from calendar file: " + Stuff
		LogIt(LogBuf)
	} else if err := scanner.Err(); err != nil {
		LogBuf = "Error reading calendar file: " + err.Error()
		LogIt(LogBuf)
	}
	CloseCalendarFile()

	pCalendar.Year      = StrToInt(StrGetWord(Stuff, 1))
	pCalendar.Month     = StrToInt(StrGetWord(Stuff, 2))
	pCalendar.Day       = StrToInt(StrGetWord(Stuff, 3))
	pCalendar.Hour      = StrToInt(StrGetWord(Stuff, 4))
	pCalendar.DayOfWeek = StrToInt(StrGetWord(Stuff, 5))

	if pCalendar.Year <= 0 {
		// Invalid year
		pCalendar.Year = 1
		LogBuf = "Calendar::GetStartTime - Year forced to 1"
		LogIt(LogBuf)
	}
	if pCalendar.Month <= 0 {
		// Invalid month
		pCalendar.Month = 1
		LogBuf = "Calendar::GetStartTime - Month forced to 1"
		LogIt(LogBuf)
	}
	if pCalendar.Day <= 0 {
		// Invalid day
		pCalendar.Day = 1
		LogBuf = "Calendar::GetStartTime - Day forced to 1"
		LogIt(LogBuf)
	}
	if pCalendar.Hour <= 0 {
		// Invalid hour
		pCalendar.Hour = 1
		LogBuf = "Calendar::GetStartTime - Hour forced to 1"
		LogIt(LogBuf)
	}
	if pCalendar.DayOfWeek <= 0 {
		// Invalid day of week
		pCalendar.DayOfWeek = 1
		LogBuf = "Calendar::GetStartTime - Day of Week forced to 1"
		LogIt(LogBuf)
	}

	LogBuf = "Start date and time is: "
	Buffer := fmt.Sprintf("Year: %d Month: %d Day: %d Hour: %d Day of Week: %d", pCalendar.Year, pCalendar.Month, pCalendar.Day, pCalendar.Hour, pCalendar.DayOfWeek)
	Stuff = Buffer
	LogBuf += Stuff
	LogIt(LogBuf)
}

// Format and return the current in-game date and time
func GetTime() string {
	DEBUGIT(1)
	FormattedDateTime := ""

	Stuff := pCalendar.DayNames[pCalendar.DayOfWeek-1]
	FormattedDateTime += Stuff + ", "

	Stuff = pCalendar.MonthNames[pCalendar.Month-1]
	FormattedDateTime += Stuff + " "

	Stuff = pCalendar.DayOfMonth[pCalendar.Day-1]
	FormattedDateTime += Stuff + ", "

	Buffer := fmt.Sprintf("%d ", pCalendar.Year)
	Stuff = Buffer
	FormattedDateTime += Stuff

	Stuff = pCalendar.HourNames[pCalendar.Hour-1]
	FormattedDateTime += Stuff

	return FormattedDateTime
}

// Load the day names from a file into the DayNames array
func LoadDayNamesArray() {
	DEBUGIT(1)
	DayNamesFileName := DAY_NAMES_DIR + "DayNames.txt"

	DayNamesFile, err := os.Open(DayNamesFileName)
	if err != nil {
		// Open failed
		LogBuf = "Calendar::LoadDayNamesArray - Open Day Names file failed (read)"
		LogIt(LogBuf)
		return
	}
	defer DayNamesFile.Close()

	pCalendar.DayNames = []string{}
	scanner := bufio.NewScanner(DayNamesFile)
	for scanner.Scan() {
		pCalendar.DayNames = append(pCalendar.DayNames, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		LogBuf = "Calendar::LoadDayNamesArray - Error reading Day Names file: " + err.Error()
		LogIt(LogBuf)
		return
	}

	LogBuf = "DayNames array loaded"
	LogIt(LogBuf)
}

// Load the day of month names from a file into the DayOfMonth array
func LoadDayOfMonthArray() {
	DEBUGIT(1)
	DayOfMonthFileName := DAY_NAMES_DIR + "DayOfMonth.txt"

	DayOfMonthFile, err := os.Open(DayOfMonthFileName)
	if err != nil {
		// Open failed
		LogBuf = "Calendar::LoadDayOfMonthArray - Open Day Of Month file failed (read)"
		LogIt(LogBuf)
		return
	}
	defer DayOfMonthFile.Close()

	pCalendar.DayOfMonth = []string{}
	scanner := bufio.NewScanner(DayOfMonthFile)
	for scanner.Scan() {
		pCalendar.DayOfMonth = append(pCalendar.DayOfMonth, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		LogBuf = "Calendar::LoadDayOfMonthArray - Error reading Day Of Month file: " + err.Error()
		LogIt(LogBuf)
		return
	}

	LogBuf = "DayOfMonth array loaded"
	LogIt(LogBuf)
}

// Load the hour names from a file into the HourNames array
func LoadHourNamesArray() {
	DEBUGIT(1)
	HourNamesFileName := HOUR_NAMES_DIR + "HourNames.txt"

	HourNamesFile, err := os.Open(HourNamesFileName)
	if err != nil {
		// Open failed
		LogBuf = "Calendar::LoadHourNamesArray - Open Hour Names file failed (read)"
		LogIt(LogBuf)
		return
	}
	defer HourNamesFile.Close()

	pCalendar.HourNames = []string{}
	scanner := bufio.NewScanner(HourNamesFile)
	for scanner.Scan() {
		pCalendar.HourNames = append(pCalendar.HourNames, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		LogBuf = "Calendar::LoadHourNamesArray - Error reading Hour Names file: " + err.Error()
		LogIt(LogBuf)
		return
	}

	LogBuf = "HourNames array loaded"
	LogIt(LogBuf)
}

// Load the month names from a file into the MonthNames array
func LoadMonthNamesArray() {
	DEBUGIT(1)
	MonthNamesFileName := MONTH_NAMES_DIR + "MonthNames.txt"

	MonthNamesFile, err := os.Open(MonthNamesFileName)
	if err != nil {
		// Open failed
		LogBuf = "Calendar::LoadMonthNamesArray - Open Month Names file failed (read)"
		LogIt(LogBuf)
		return
	}
	defer MonthNamesFile.Close()

	pCalendar.MonthNames = []string{}
	scanner := bufio.NewScanner(MonthNamesFile)
	for scanner.Scan() {
		pCalendar.MonthNames = append(pCalendar.MonthNames, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		LogBuf = "Calendar::LoadMonthNamesArray - Error reading Month Names file: " + err.Error()
		LogIt(LogBuf)
		return
	}

	LogBuf = "MonthNames array loaded"
	LogIt(LogBuf)
}

// Open the calendar file for reading
func OpenCalendarFile() {
	DEBUGIT(1)
	CalendarFileIsOpen = false
	CalendarFileName := CONTROL_DIR + "Calendar.txt"

	CalendarFileInp, err := os.Open(CalendarFileName)
	if err != nil {
		// Calendar file does not exist
		LogBuf = "Calendar file not found."
		LogIt(LogBuf)
		LogBuf = "Forcing start date to Year: 1 Month: 1 Day: 1 Hour: 1 Day of Week: 1"
		LogIt(LogBuf)
		return
	}

	// Open was successful
	CalendarFile = CalendarFileInp
	CalendarFileIsOpen = true
}

// Save the current in-game time to the calendar file
func SaveTime() {
  DEBUGIT(1)
  CalendarFileName := CONTROL_DIR + "Calendar.txt"

  CalendarFileOut, err := os.Create(CalendarFileName)
  if err != nil {
    // Open failed
    LogBuf = "Calendar::SaveTime - Open calendar file - Failed"
    LogIt(LogBuf)
    return
  }
  defer CalendarFileOut.Close()

  Buffer := fmt.Sprintf("%d %d %d %d %d", pCalendar.Year, pCalendar.Month, pCalendar.Day, pCalendar.Hour, pCalendar.DayOfWeek)
  _, err = CalendarFileOut.WriteString(Buffer)
  if err != nil {
    LogBuf = "Calendar::SaveTime - Error writing to calendar file: " + err.Error()
    LogIt(LogBuf)
    return
  }

  LogBuf = "Calendar time saved successfully"
  LogIt(LogBuf)
}