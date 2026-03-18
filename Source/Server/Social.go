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
    // If the file isn't, then all socials are bad commands, doh!
    return false
  }
  MsgText = "Not Done"
  for MsgText != "End of Socials" {
    // Loop until social is found or end of file
    MsgText = ReadLine()
    MsgText = StrTrimLeft(MsgText)
    if StrLeft(MsgText, 9) == "Social : " {
      // Ok, a social has been found
      TmpStr = StrRight(MsgText, StrGetLength(MsgText)-9)
      if TmpStr == MudCmd {
        // THE social has been found
        Found = true
        MsgText = ReadLine()
        MinPos = StrRight(MsgText, StrGetLength(MsgText)-9)
        if PositionNotOk(pDnodeActor, MinPos) {
          // Player is not in the minimum position
          pDnodeActor.PlayerOut += "You are not in a position to that right now.\r\n"
          CreatePrompt(pDnodeActor.pPlayer)
          pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
          MsgText = ""
          break
        }
        // It is a social and player is in minimum position
        Socialize(MinPos, MsgText)
      }
    }
  }
  CloseSocialFile()
  if Found {
    // Return true so command processor will exit properly
    return true
  } else {
    // Return false so command processor will tell player bad command
    return false
  }
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
  var err error
  SocialFile, err = os.Open(SocialFileName)
  if err == nil {
    SocialScanner = bufio.NewScanner(SocialFile)
    return true
  } else {
    return false
  }
}

// Is player in a valid position for the social
func PositionNotOk(pDnode *Dnode, MinPos string) bool {
  var MinPosNbr int
  var PlayerPosNbr int

  MinPosNbr = PosNbr(MinPos)
  PlayerPosNbr = PosNbr(pDnodeActor.pPlayer.Position)
  if PlayerPosNbr < MinPosNbr {
    return true
  } else {
    return false
  }
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
  var MsgText string

  SocialScanner.Scan()
  MsgText = SocialScanner.Text()
  return MsgText
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
    // The social does not accept a target
    if TargetName != "" {
      // But, target was given
      pDnodeSrc.PlayerOut += MudCmd
      pDnodeSrc.PlayerOut += " does not use a target.\r\n"
      CreatePrompt(pDnodeSrc.pPlayer)
      pDnodeSrc.PlayerOut += GetOutput(pDnodeSrc.pPlayer)
      return
    }
  }
  // All checks complete, get on with the social
  if TargetName == "" {
    // Social without target
    MsgText = ReadLine()
    MsgText = PronounSubstitute(MsgText)
    SendToPlayer(MsgText)
    MsgText = ReadLine()
    MsgText = PronounSubstitute(MsgText)
    SendToRoom(pDnodeActor.pPlayer.RoomId, MsgText)
    return
  }
  if PlayerName == TargetName {
    // Social with target equal self
    for i = 1; i <= 3; i++ {
      MsgText = ReadLine()
    }
    MsgText = PronounSubstitute(MsgText)
    SendToPlayer(MsgText)
    return
  }
  // Do some checks to determine if target is valid
  pDnodeTgt = GetTargetDnode(TargetName)
  if pDnodeTgt == nil {
    // Target is not online and/or not in 'playing' state
    TargetNotHere = true
  } else {
    // Target is online and playing
    if pDnodeSrc.pPlayer.RoomId != pDnodeTgt.pPlayer.RoomId {
      // Target is not in the same room
      TargetNotHere = true
    }
  }
  if TargetNotHere {
    // Target is not playing or is not in same room as player
    for i = 1; i <= 4; i++ {
      MsgText = ReadLine()
    }
    MsgText = PronounSubstitute(MsgText)
    SendToPlayer(MsgText)
    return
  } else {
    // Target is playing and is in same room as player
    if PositionNotOk(pDnodeTgt, MinPos) {
      // Target is not in minimum position for social
      pDnodeSrc.PlayerOut += pDnodeTgt.PlayerName
      pDnodeSrc.PlayerOut += " is "
      pDnodeSrc.PlayerOut += pDnodeTgt.pPlayer.Position
      if pDnodeTgt.pPlayer.Position == "sit" {
        // Add the extra 't' so it comes out sitting vs siting
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
    // Message to player
    MsgText = PronounSubstitute(MsgText)
    SendToPlayer(MsgText)
    MsgText = ReadLine()
    // Message to target
    MsgText = PronounSubstitute(MsgText)
    SendToTarget(pDnodeTgt, MsgText)
    MsgText = ReadLine()
    // Message to the others
    MsgText = PronounSubstitute(MsgText)
    SendToRoom(pDnodeActor.pPlayer.RoomId, MsgText)
  }
}
