//*******************************************************************
// HolyQuest: Powered by the Online Multi-User Game Server (OMugs)  *
// File:      Calendar.go                                           *
// Usage:     Maintains the game calendar                           *
// This file is part of the HolyQuestGo project.                    *
// It is licensed under the Unlicense.                              *
// See the LICENSE file in the project root for details.            *
//*******************************************************************

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

var pCalendar =  &Calendar{
  Year:              1,
  Month:             1,
  Day:               1,
  Hour:              1,
  DayOfWeek:         1,
  TimeToAdvanceHour: 60,
}

// Initializes the calendar system
func CalendarConstructor() {
  DEBUGIT(1)
  timer := time.Now().Unix()
  pCalendar.TimeToAdvanceHour = timer + int64(REAL_MINUTES_PER_HOUR*60)
  pCalendar.Year      = 1
  pCalendar.Month     = 1
  pCalendar.Day       = 1
  pCalendar.Hour      = 1
  pCalendar.DayOfWeek = 1
  OpenCalendarFile()
  if CalendarFileIsOpen {
    // Get start time
    GetStartTime()
  }
  LoadDayNamesArray()
  LoadDayOfMonthArray()
  LoadHourNamesArray()
  LoadMonthNamesArray()
}

// Cleans up the calendar system
func CalendarDestructor() {
  DEBUGIT(1)
  SaveTime()
}

// Advances the in-game time by one hour
func AdvanceTime() {
  DEBUGIT(5)
  NowSec := time.Now().Unix()
  if NowSec < pCalendar.TimeToAdvanceHour {
    // Not time to advance the hour
    return
  }
  pCalendar.TimeToAdvanceHour = NowSec + int64(REAL_MINUTES_PER_HOUR*60)
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

// Format and return the current in-game date and time
func GetTime() string {
  DEBUGIT(1)
  FormattedDateTime := ""

  Stuff := pCalendar.DayNames[pCalendar.DayOfWeek-1]
  FormattedDateTime += Stuff
  FormattedDateTime += ", "

  Stuff = pCalendar.MonthNames[pCalendar.Month-1]
  FormattedDateTime += Stuff
  FormattedDateTime += " "

  Stuff = pCalendar.DayOfMonth[pCalendar.Day-1]
  FormattedDateTime += Stuff
  FormattedDateTime += ", "

  Buffer := fmt.Sprintf("%d ", pCalendar.Year)
  Stuff = Buffer
  FormattedDateTime += Stuff

  Stuff = pCalendar.HourNames[pCalendar.Hour-1]
  FormattedDateTime += Stuff

  return FormattedDateTime
}
// Closes the calendar file
func CloseCalendarFile() {
  DEBUGIT(1)
  if CalendarFile != nil {
    CalendarFile.Close()
  }
}

// Retrieves the start time from the calendar file
func GetStartTime() {
  DEBUGIT(1)
  scanner := bufio.NewScanner(CalendarFile)
  stuff := ""
  scanner.Scan()
  stuff = scanner.Text()
  CloseCalendarFile()
  pCalendar.Year      = StrToInt(StrGetWord(stuff, 1))
  pCalendar.Month     = StrToInt(StrGetWord(stuff, 2))
  pCalendar.Day       = StrToInt(StrGetWord(stuff, 3))
  pCalendar.Hour      = StrToInt(StrGetWord(stuff, 4))
  pCalendar.DayOfWeek = StrToInt(StrGetWord(stuff, 5))
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
    pCalendar.Hour = 1
    LogBuf = "Calendar::GetStartTime - Day of Week forced to 1"
    LogIt(LogBuf)
  }
  LogBuf = "Start date and time is: "
  buf := fmt.Sprintf(
    "Year: %d Month: %d Day: %d Hour: %d Day of Week: %d",
    pCalendar.Year, pCalendar.Month, pCalendar.Day, pCalendar.Hour, pCalendar.DayOfWeek,
  )
  LogBuf += buf
  LogIt(LogBuf)
}

// Opens the calendar file for reading
func OpenCalendarFile() {
  var CalendarFileName string

  DEBUGIT(1)
  CalendarFileIsOpen = false
  CalendarFileName = CONTROL_DIR
  CalendarFileName += "Calendar.txt"
  CalendarFileInp, err := os.Open(CalendarFileName)
  if err != nil {
    // Calendar file does not exist
    LogBuf = "Calendar file not found."
    LogIt(LogBuf)
    LogBuf = "Forcing start date to Year: 1 Month: 1 Day: 1 Hour: 1 Day of Week: 1"
    LogIt(LogBuf)
    return
  }
  CalendarFile = CalendarFileInp
  // Open was successful
  CalendarFileIsOpen = true
}

// Loads the day names from a file into the DayNames array
func LoadDayNamesArray() {
  var DayNamesFileName string

  DEBUGIT(1)
  DayNamesFileName = DAY_NAMES_DIR
  DayNamesFileName += "DayNames.txt"
  file, err := os.Open(DayNamesFileName)
  if err != nil {
    // Open failed
    LogBuf = "Calendar::LoadDayNamesArray - Open Day Names file failed (read)"
    LogIt(LogBuf)
    return
  }
  pCalendar.DayNames = nil
  scanner := bufio.NewScanner(file)
  scanner.Scan()
  Stuff = scanner.Text()
  for scanner.Scan() {
    // Read all day names
    pCalendar.DayNames = append(pCalendar.DayNames, Stuff)
    Stuff = scanner.Text()
  }
  pCalendar.DayNames = append(pCalendar.DayNames, Stuff)
  file.Close()
  LogBuf = "DayNames array loaded"
  LogIt(LogBuf)
}

// Loads the day of month names from a file into the DayOfMonth array
func LoadDayOfMonthArray() {
  var DayOfMonthFileName string

  DEBUGIT(1)
  DayOfMonthFileName = DAY_OF_MONTH_DIR
  DayOfMonthFileName += "DayOfMonth.txt"
  file, err := os.Open(DayOfMonthFileName)
  if err != nil {
    // Open failed
    LogBuf = "Calendar::LoadDayOfMonthArray - Open Day Of Month file failed (read)"
    LogIt(LogBuf)
    return
  }
  pCalendar.DayOfMonth = nil
  scanner := bufio.NewScanner(file)
  scanner.Scan()
  Stuff = scanner.Text()
  for scanner.Scan() {
    // Read all day of month
    pCalendar.DayOfMonth = append(pCalendar.DayOfMonth, Stuff)
    Stuff = scanner.Text()
  }
  pCalendar.DayOfMonth = append(pCalendar.DayOfMonth, Stuff)
  file.Close()
  LogBuf = "DayOfMonth array loaded"
  LogIt(LogBuf)
}

// Loads the hour names from a file into the HourNames array
func LoadHourNamesArray() {
  var HourNamesFileName string

  DEBUGIT(1)
  HourNamesFileName = HOUR_NAMES_DIR
  HourNamesFileName += "HourNames.txt"
  file, err := os.Open(HourNamesFileName)
  if err != nil {
    // Open failed
    LogBuf = "Calendar::LoadHourNamesArray - Open Hour Names file failed (read)"
    LogIt(LogBuf)
    return
  }
  pCalendar.HourNames = nil
  scanner := bufio.NewScanner(file)
  scanner.Scan()
  Stuff = scanner.Text()
  for scanner.Scan() {
    // Read all hour names
    pCalendar.HourNames = append(pCalendar.HourNames, Stuff)
    Stuff = scanner.Text()
  }
  pCalendar.HourNames = append(pCalendar.HourNames, Stuff)
  file.Close()
  LogBuf = "HourNames array loaded"
  LogIt(LogBuf)
}

// Loads the month names from a file into the MonthNames array
func LoadMonthNamesArray() {
  var MonthNamesFileName string

  DEBUGIT(1)
  MonthNamesFileName = MONTH_NAMES_DIR
  MonthNamesFileName += "MonthNames.txt"
  file, err := os.Open(MonthNamesFileName)
  if err != nil {
    // Open failed
    LogBuf = "Calendar::LoadMonthNamesArray - Open Month Names file failed (read)"
    LogIt(LogBuf)
    return
  }
  pCalendar.MonthNames = nil
  scanner := bufio.NewScanner(file)
  scanner.Scan()
  Stuff = scanner.Text()
  for scanner.Scan() {
    // Read all month names
    pCalendar.MonthNames = append(pCalendar.MonthNames, Stuff)
    Stuff = scanner.Text()
  }
  pCalendar.MonthNames = append(pCalendar.MonthNames, Stuff)
  file.Close()
  LogBuf = "MonthNames array loaded"
  LogIt(LogBuf)
}

// Saves the current in-game time to the calendar file
func SaveTime() {
  var CalendarFileName string

  DEBUGIT(1)
  CalendarFileName = CONTROL_DIR
  CalendarFileName += "Calendar.txt"
  file, err := os.Create(CalendarFileName)
  if err != nil {
    // Open failed
    LogBuf = "Calendar::SaveTime - Open calendar file - Failed"
    LogIt(LogBuf)
    return
  }
  buffer := fmt.Sprintf("%d %d %d %d %d", pCalendar.Year, pCalendar.Month, pCalendar.Day, pCalendar.Hour, pCalendar.DayOfWeek)
  file.WriteString(buffer)
  file.Close()
}
