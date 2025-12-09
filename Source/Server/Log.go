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

// Open log file
func OpenLogFile() {
	LogFileName := LOG_DIR + "SrvrLog.txt"
	if _, Err := os.Stat(LogFileName); Err == nil {
		LogTime := fmt.Sprintf("%d", time.Now().Unix())
		LogSaveFileName := LogFileName[:len(LogFileName)-4] + "." + LogTime + ".txt"
		os.Rename(LogFileName, LogSaveFileName)
	}
	var Err error
	LogFile, Err = os.Create(LogFileName)
	if Err != nil {
		fmt.Println("OpenLogFile() - Failed")
		fmt.Println("Hard Exit!")
		os.Exit(1)
	}
	fmt.Println("Log File is Open")
}

// Close log file
func CloseLogFile() {
	if LogFile != nil {
		LogFile.Close()
		LogFile = nil
	}
}

// Write log file
func LogIt(Message string) {
	if LogFile == nil {
		fmt.Println("Log file is not open")
		return
	}
	DisplayCurrentTime := time.Now().Format("2006-01-02 15:04:05 ")
	LogMessage := DisplayCurrentTime + Message + "\n"
	LogFile.WriteString(LogMessage)
	LogFile.Sync()
}