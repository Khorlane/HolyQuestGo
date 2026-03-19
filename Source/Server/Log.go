//*******************************************************************
// HolyQuest: Powered by the Online Multi-User Game Server (OMugs)  *
// File:      Log.go                                                *
// Usage:     Logs messages to disk file                            *
// This file is part of the HolyQuestGo project.                    *
// It is licensed under the Unlicense.                              *
// See the LICENSE file in the project root for details.            *
//*******************************************************************

package server

import (
  "fmt"
  "os"
  "time"
)

var LogFile *os.File

// Close log file
func CloseLogFile() {
  LogFile.Close()
}

// Write log file
func LogIt(LogBuf string) {
  DisplayCurrentTime := time.Now().Format("2006-01-02 15:04:05 ")
  LogBuf = DisplayCurrentTime + LogBuf
  LogBuf += "\n"
  LogFile.WriteString(LogBuf)
  LogFile.Sync()
}

// Open log file
func OpenLogFile() {
  var LogFileName string
  var LogSaveFileName string
  var LogTime string

  PrintIt("Log::OpenLogFile()")
  LogFileName = LOG_DIR
  LogFileName += "SrvrLog.txt"
  if FileExist(LogFileName) {
    Buf = fmt.Sprintf("%d", GetTimeSeconds())
    LogTime = Buf

    LogSaveFileName = StrLeft(LogFileName, StrGetLength(LogFileName)-4)
    LogSaveFileName += "."
    LogSaveFileName += LogTime
    LogSaveFileName += ".txt"
    Rename(LogFileName, LogSaveFileName)
  }
  LogFile, ErrorCode = os.Create(LogFileName)
  if ErrorCode != nil {
    PrintIt("Log::OpenLogFile() - Failed")
    PrintIt("Hard Exit!")
    os.Exit(1)
  }
  PrintIt("Log File is Open")
}
