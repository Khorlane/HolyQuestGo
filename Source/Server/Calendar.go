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
  OpenCalendarFile()
  if CalendarFileIsOpen {
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
  if CalendarFile == nil {
    LogBuf = "Calendar::GetStartTime - Calendar file handle is nil"
    LogIt(LogBuf)
    return
  }
  scanner := bufio.NewScanner(CalendarFile)
  stuff := ""
  if scanner.Scan() {
    stuff = scanner.Text()
  }
  CloseCalendarFile()
  pCalendar.Year      = StrToInt(StrGetWord(stuff, 1))
  pCalendar.Month     = StrToInt(StrGetWord(stuff, 2))
  pCalendar.Day       = StrToInt(StrGetWord(stuff, 3))
  pCalendar.Hour      = StrToInt(StrGetWord(stuff, 4))
  pCalendar.DayOfWeek = StrToInt(StrGetWord(stuff, 5))
  if pCalendar.Year <= 0 {
    pCalendar.Year = 1
    LogBuf = "Calendar::GetStartTime - Year forced to 1"
    LogIt(LogBuf)
  }
  if pCalendar.Month <= 0 {
    pCalendar.Month = 1
    LogBuf = "Calendar::GetStartTime - Month forced to 1"
    LogIt(LogBuf)
  }
  if pCalendar.Day <= 0 {
    pCalendar.Day = 1
    LogBuf = "Calendar::GetStartTime - Day forced to 1"
    LogIt(LogBuf)
  }
  if pCalendar.Hour <= 0 {
    pCalendar.Hour = 1
    LogBuf = "Calendar::GetStartTime - Hour forced to 1"
    LogIt(LogBuf)
  }
  if pCalendar.DayOfWeek <= 0 {
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

// Loads the day names from a file into the DayNames array
func LoadDayNamesArray() {
  DEBUGIT(1)
  filename := DAY_NAMES_DIR + "DayNames.txt"
  file, err := os.Open(filename)
  if err != nil {
    LogBuf = "LoadDayNamesArray - DayNames file not found."
    LogIt(LogBuf)
    return
  }
  defer file.Close()
  pCalendar.DayNames = nil
  scanner := bufio.NewScanner(file)
  for scanner.Scan() {
    line := scanner.Text()
    if line != "" {
      pCalendar.DayNames = append(pCalendar.DayNames, line)
    }
  }
}

// Loads the day of month names from a file into the DayOfMonth array
func LoadDayOfMonthArray() {
  DEBUGIT(1)
  filename := DAY_OF_MONTH_DIR + "DayOfMonth.txt"
  file, err := os.Open(filename)
  if err != nil {
    LogBuf = "LoadDayOfMonthArray - DayOfMonth file not found."
    LogIt(LogBuf)
    return
  }
  defer file.Close()
  pCalendar.DayOfMonth = nil
  scanner := bufio.NewScanner(file)
  for scanner.Scan() {
    line := scanner.Text()
    if line != "" {
      pCalendar.DayOfMonth = append(pCalendar.DayOfMonth, line)
    }
  }
}

// Loads the hour names from a file into the HourNames array
func LoadHourNamesArray() {
  DEBUGIT(1)
  filename := HOUR_NAMES_DIR + "HourNames.txt"
  file, err := os.Open(filename)
  if err != nil {
    LogBuf = "LoadHourNamesArray - HourNames file not found."
    LogIt(LogBuf)
    return
  }
  defer file.Close()
  pCalendar.HourNames = nil
  scanner := bufio.NewScanner(file)
  for scanner.Scan() {
    line := scanner.Text()
    if line != "" {
      pCalendar.HourNames = append(pCalendar.HourNames, line)
    }
  }
}

// Loads the month names from a file into the MonthNames array
func LoadMonthNamesArray() {
  DEBUGIT(1)
  filename := MONTH_NAMES_DIR + "MonthNames.txt"
  file, err := os.Open(filename)
  if err != nil {
    LogBuf = "LoadMonthNamesArray - MonthNames file not found."
    LogIt(LogBuf)
    return
  }
  defer file.Close()
  pCalendar.MonthNames = nil
  scanner := bufio.NewScanner(file)
  for scanner.Scan() {
    line := scanner.Text()
    if line != "" {
      pCalendar.MonthNames = append(pCalendar.MonthNames, line)
    }
  }
}

// Saves the current in-game time to the calendar file
func SaveTime() {
  DEBUGIT(1)
  filename := CONTROL_DIR + "Calendar.txt"
  file, err := os.Create(filename)
  if err != nil {
    LogBuf = "Calendar::SaveTime - Open calendar file - Failed"
    LogIt(LogBuf)
    return
  }
  defer file.Close()
  buffer := fmt.Sprintf("%d %d %d %d %d", pCalendar.Year, pCalendar.Month, pCalendar.Day, pCalendar.Hour, pCalendar.DayOfWeek)
  _, err = file.WriteString(buffer)
  if err != nil {
    LogBuf = "Calendar::SaveTime - Error writing to calendar file: " + err.Error()
    LogIt(LogBuf)
    return
  }
}