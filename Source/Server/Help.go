//*******************************************************************
// HolyQuest: Powered by the Online Multi-User Game Server (OMugs)  *
// File:      Help.go                                               *
// Usage:     Provides help information for players                 *
// This file is part of the HolyQuestGo project.                    *
// It is licensed under the Unlicense.                              *
// See the LICENSE file in the project root for details.            *
//*******************************************************************

package server

import (
  "bufio"
  "os"
)

var HelpFile    *os.File
var HelpScanner *bufio.Scanner
var HelpText     string

// Is it a Help topic?
func IsHelp() bool {
  DEBUGIT(1)
  Found := false
  if !OpenHelpFile() {
    // If the file isn't there, then all Helps are not found, doh!
    return false
  }
  HelpLookup := StrGetWord(CmdStr, 2)
  HelpLookup  = StrMakeLower(HelpLookup)
  HelpText    = "Not Done"
  for HelpText != "End of Help" {
    // Loop until Help is found or end of file
    HelpReadLine()
    HelpText = StrTrimLeft(HelpText)
    TmpStr   = StrLeft(HelpText, 5)
    if TmpStr == "Help:" {
      // Ok, a Help entry has been found
      TmpStr = StrRight(HelpText, StrGetLength(HelpText)-5)
      TmpStr = StrMakeLower(TmpStr)
      if TmpStr == HelpLookup {
        // THE Help entry has been found, show it to player
        Found = true
        ShowHelp()
        HelpText = "End of Help"
      }
    }
  }
  CloseHelpFile()
  if Found {
    // Return true so command processor will exit properly
    return true
  } else {
    // Return false so command processor will tell player bad command
    return false
  }
}

// Open Help file
func OpenHelpFile() bool {
  DEBUGIT(1)
  HelpFileName := HELP_DIR
  HelpFileName += "Help.txt"
  var err error
  HelpFile, err = os.Open(HelpFileName)
  if err != nil {
    return false
  } else {
    HelpScanner = bufio.NewScanner(HelpFile)
    return true
  }
}

// Read a line from Help file
func HelpReadLine() {
  DEBUGIT(1)
  HelpScanner.Scan()
  HelpText = HelpScanner.Text()
}

// Show help to player
func ShowHelp() {
  DEBUGIT(1)
  TmpStr = StrLeft(HelpText, 13)
  for TmpStr != "Related help:" {
    HelpReadLine()
    pDnodeActor.PlayerOut += HelpText
    pDnodeActor.PlayerOut += "\r\n"
    TmpStr = StrLeft(HelpText, 13)
  }
  pDnodeActor.PlayerOut += "\r\n"
  CreatePrompt(pDnodeActor.pPlayer)
  pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
}

// Close Help file
func CloseHelpFile() {
  DEBUGIT(1)
  HelpFile.Close()
}
