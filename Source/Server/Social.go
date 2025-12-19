//*******************************************************************
// HolyQuest: Powered by the Online Multi-User Game Server (OMugs)  *
// File:      Social.go                                             *
// Usage:     Manages social interactions between players           *
// This file is part of the HolyQuestGo project.                    *
// It is licensed under the Unlicense.                              *
// See the LICENSE file in the project root for details.            *
//*******************************************************************

package server

import (
  "bufio"
  "os"
)

var SocialFile    *os.File
var SocialScanner *bufio.Scanner

// Is MudCmd a social command?
func IsSocial() bool {
  var Found bool
  var MinPos string
  var MsgText string

  Found = false
  if !OpenSocialFile() {
    return false
  }
  MsgText = "Not Done"
  for MsgText != "End of Socials" {
    MsgText = ReadLine()
    MsgText = StrTrimLeft(MsgText)
    if StrLeft(MsgText, 9) == "Social : " {
      TmpStr = StrRight(MsgText, StrGetLength(MsgText)-9)
      if TmpStr == MudCmd {
        Found = true
        MsgText = ReadLine()
        MinPos = StrRight(MsgText, StrGetLength(MsgText)-9)
        if PositionNotOk(pDnodeActor, MinPos) {
          pDnodeActor.PlayerOut += "You are not in a position to that right now.\r\n"
          CreatePrompt(pDnodeActor.pPlayer)
          pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
          MsgText = ""
          break
        }
        Socialize(MinPos, MsgText)
      }
    }
  }
  CloseSocialFile()
  if Found {
    return true
  }
  return false
}

// Close social file
func CloseSocialFile() {
  if SocialFile != nil {
    SocialFile.Close()
  }
}

// Open social file
func OpenSocialFile() bool {
  var SocialFileName string

  SocialFileName = SOCIAL_DIR
  SocialFileName += "Social.txt"
  f, err := os.Open(SocialFileName)
  if err != nil {
    return false
  }
  SocialFile = f
  SocialScanner = bufio.NewScanner(SocialFile)
  return true
}

// Is player in a valid position for the social
func PositionNotOk(pDnode *Dnode, MinPos string) bool {
  var MinPosNbr int
  var PlayerPosNbr int

  MinPosNbr = PosNbr(MinPos)
  PlayerPosNbr = PosNbr(pDnodeActor.pPlayer.Position)
  if PlayerPosNbr < MinPosNbr {
    return true
  }
  return false
}

// Convert position to a number
func PosNbr(Position string) int {
  if Position == "sleep" {
    return 1
  }
  if Position == "sit" {
    return 2
  }
  if Position == "stand" {
    return 3
  }
  return -1
}

// Read a line from social file
func ReadLine() string {
  if SocialScanner != nil {
    if SocialScanner.Scan() {
      return SocialScanner.Text()
    }
  }
  return ""
}

// Send substituted message to player
func SendToPlayer(MsgText string) {
  pDnodeActor.PlayerOut += MsgText
  pDnodeActor.PlayerOut += "\r\n"
  CreatePrompt(pDnodeActor.pPlayer)
  pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
}

// Send substituted message to target
func SendToTarget(pDnodeTgt1 *Dnode, MsgText string) {
  pDnodeTgt = pDnodeTgt1
  pDnodeTgt.PlayerOut += "\r\n"
  pDnodeTgt.PlayerOut += MsgText
  pDnodeTgt.PlayerOut += "\r\n"
  CreatePrompt(pDnodeTgt.pPlayer)
  pDnodeTgt.PlayerOut += GetOutput(pDnodeTgt.pPlayer)
}

// Social command
func Socialize(MinPos string, MsgText string) {
  var i             int
  var LineCount     int
  var PlayerName    string
  var TargetName    string
  var TargetNotHere bool

  pDnodeSrc = pDnodeActor
  pDnodeTgt = nil
  TargetNotHere = false
  MsgText = ReadLine()
  LineCount = StrToInt(StrRight(MsgText, StrGetLength(MsgText)-9))
  PlayerName = pDnodeActor.PlayerName
  PlayerName = StrMakeLower(PlayerName)
  TargetName = StrGetWord(CmdStr, 2)
  TargetName = StrMakeLower(TargetName)
  if LineCount == 2 {
    if TargetName != "" {
      pDnodeSrc.PlayerOut += MudCmd
      pDnodeSrc.PlayerOut += " does not use a target.\r\n"
      CreatePrompt(pDnodeSrc.pPlayer)
      pDnodeSrc.PlayerOut += GetOutput(pDnodeSrc.pPlayer)
      return
    }
  }
  if TargetName == "" {
    MsgText = ReadLine()
    MsgText = PronounSubstitute(MsgText)
    SendToPlayer(MsgText)
    MsgText = ReadLine()
    MsgText = PronounSubstitute(MsgText)
    SendToRoom(pDnodeActor.pPlayer.RoomId, MsgText)
    return
  }
  if PlayerName == TargetName {
    for i = 1; i <= 3; i++ {
      MsgText = ReadLine()
    }
    MsgText = PronounSubstitute(MsgText)
    SendToPlayer(MsgText)
    return
  }
  pDnodeTgt = GetTargetDnode(TargetName)
  if pDnodeTgt == nil {
    TargetNotHere = true
  } else {
    if pDnodeSrc.pPlayer.RoomId != pDnodeTgt.pPlayer.RoomId {
      TargetNotHere = true
    }
  }
  if TargetNotHere {
    for i = 1; i <= 4; i++ {
      MsgText = ReadLine()
    }
    MsgText = PronounSubstitute(MsgText)
    SendToPlayer(MsgText)
    return
  } else {
    if PositionNotOk(pDnodeTgt, MinPos) {
      pDnodeSrc.PlayerOut += pDnodeTgt.PlayerName
      pDnodeSrc.PlayerOut += " is "
      pDnodeSrc.PlayerOut += pDnodeTgt.pPlayer.Position
      if pDnodeTgt.pPlayer.Position == "sit" {
        pDnodeSrc.PlayerOut += "t"
      }
      pDnodeSrc.PlayerOut += "ing and cannot participate.\r\n"
      CreatePrompt(pDnodeSrc.pPlayer)
      pDnodeSrc.PlayerOut += GetOutput(pDnodeSrc.pPlayer)
      return
    }
    for i = 1; i <= 5; i++ {
      MsgText = ReadLine()
    }
    MsgText = PronounSubstitute(MsgText)
    SendToPlayer(MsgText)
    MsgText = ReadLine()
    MsgText = PronounSubstitute(MsgText)
    SendToTarget(pDnodeTgt, MsgText)
    MsgText = ReadLine()
    MsgText = PronounSubstitute(MsgText)
    SendToRoom(pDnodeActor.pPlayer.RoomId, MsgText)
  }
}