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

var HelpFile *os.File
var HelpScanner *bufio.Scanner
var HelpText string

func IsHelp() bool {
	DEBUGIT(1)
	Found := false
	if !OpenFile() {
		return false
	}
	HelpLookup := StrGetWord(CmdStr, 2)
	HelpLookup = StrMakeLower(HelpLookup)
	HelpText = "Not Done"
	for HelpText != "End of Help" {
		ReadLine()
		HelpText = StrTrimLeft(HelpText)
		TmpStr = StrLeft(HelpText, 5)
		if TmpStr == "Help:" {
			TmpStr = StrRight(HelpText, StrGetLength(HelpText)-5)
			TmpStr = StrMakeLower(TmpStr)
			if TmpStr == HelpLookup {
				Found = true
				ShowHelp()
				HelpText = "End of Help"
			}
		}
	}
	if HelpFile != nil {
		HelpFile.Close()
	}
	if Found {
		return true
	}
	return false
}

func OpenFile() bool {
	DEBUGIT(1)
	HelpFileName := HELP_DIR + "Help.txt"
	file, err := os.Open(HelpFileName)
	if err != nil {
		return false
	}
	HelpFile = file
	HelpScanner = bufio.NewScanner(HelpFile)
	return true
}

func ReadLine() {
	DEBUGIT(1)
	if HelpScanner != nil && HelpScanner.Scan() {
		HelpText = HelpScanner.Text()
	} else {
		HelpText = ""
	}
}

func ShowHelp() {
	TmpStr = StrLeft(HelpText, 13)
	for TmpStr != "Related help:" {
		ReadLine()
		pDnodeActor.PlayerOut += HelpText
		pDnodeActor.PlayerOut += "\r\n"
		TmpStr = StrLeft(HelpText, 13)
	}
	pDnodeActor.PlayerOut += "\r\n"
	CreatePrompt(pDnodeActor.pPlayer)
	pDnodeActor.PlayerOut += GetPlayerOutput(pDnodeActor.pPlayer)
}